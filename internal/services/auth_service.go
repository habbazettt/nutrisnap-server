package services

import (
	"context"
	"errors"
	"time"

	"github.com/habbazettt/nutrisnap-server/internal/dto"
	"github.com/habbazettt/nutrisnap-server/internal/models"
	"github.com/habbazettt/nutrisnap-server/internal/repositories"
	"github.com/habbazettt/nutrisnap-server/pkg/constants"
	"github.com/habbazettt/nutrisnap-server/pkg/jwt"
	"github.com/habbazettt/nutrisnap-server/pkg/oauth"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrOAuthFailed        = errors.New("OAuth authentication failed")
)

type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetGoogleAuthURL(state string) string
	GoogleCallback(ctx context.Context, code string) (*dto.LoginResponse, error)
}

type authService struct {
	userRepo    repositories.UserRepository
	jwtManager  *jwt.Manager
	googleOAuth *oauth.GoogleOAuth
}

func NewAuthService(userRepo repositories.UserRepository, jwtManager *jwt.Manager, googleOAuth *oauth.GoogleOAuth) AuthService {
	return &authService{
		userRepo:    userRepo,
		jwtManager:  jwtManager,
		googleOAuth: googleOAuth,
	}
}

func (s *authService) Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// Create user model
	user := &models.User{
		Email: req.Email,
		Name:  req.Name,
		Role:  models.RoleUser,
	}

	// Hash password
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Build response
	return &dto.RegisterResponse{
		User:    s.toUserResponse(user),
		Message: constants.GetStatusMessage(constants.StatusCreated),
	}, nil
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokenResponse(user)
}

// GetGoogleAuthURL returns the URL to redirect user for Google OAuth
func (s *authService) GetGoogleAuthURL(state string) string {
	return s.googleOAuth.GetAuthURL(state)
}

// GoogleCallback handles the Google OAuth callback
func (s *authService) GoogleCallback(ctx context.Context, code string) (*dto.LoginResponse, error) {
	// Exchange code for token
	token, err := s.googleOAuth.Exchange(ctx, code)
	if err != nil {
		return nil, ErrOAuthFailed
	}

	// Get user info from Google
	googleUser, err := s.googleOAuth.GetUserInfo(ctx, token)
	if err != nil {
		return nil, ErrOAuthFailed
	}

	// Find or create user
	user, err := s.userRepo.FindByEmail(googleUser.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			// Create new user
			now := time.Now()
			user = &models.User{
				Email:           googleUser.Email,
				Name:            googleUser.Name,
				Role:            models.RoleUser,
				AvatarURL:       &googleUser.Picture,
				EmailVerifiedAt: &now, // Google emails are verified
				GoogleID:        &googleUser.ID,
			}

			if err := s.userRepo.Create(user); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		// Update Google ID if not set
		if user.GoogleID == nil {
			user.GoogleID = &googleUser.ID
			if googleUser.Picture != "" && user.AvatarURL == nil {
				user.AvatarURL = &googleUser.Picture
			}
			if err := s.userRepo.Update(user); err != nil {
				return nil, err
			}
		}
	}

	return s.generateTokenResponse(user)
}

func (s *authService) generateTokenResponse(user *models.User) (*dto.LoginResponse, error) {
	// Generate access token
	accessToken, expiresAt, err := s.jwtManager.GenerateAccessToken(
		user.ID.String(),
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, _, err := s.jwtManager.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		User:         s.toUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *authService) toUserResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:              user.ID.String(),
		Email:           user.Email,
		Name:            user.Name,
		Role:            string(user.Role),
		EmailVerifiedAt: user.EmailVerifiedAt,
		CreatedAt:       user.CreatedAt,
	}
}
