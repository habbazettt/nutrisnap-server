package services

import (
	"errors"

	"github.com/habbazettt/nutrisnap-server/internal/dto"
	"github.com/habbazettt/nutrisnap-server/internal/models"
	"github.com/habbazettt/nutrisnap-server/internal/repositories"
	"github.com/habbazettt/nutrisnap-server/pkg/constants"
	"github.com/habbazettt/nutrisnap-server/pkg/jwt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userRepo   repositories.UserRepository
	jwtManager *jwt.Manager
}

func NewAuthService(userRepo repositories.UserRepository, jwtManager *jwt.Manager) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
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
