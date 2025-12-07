package services

import (
	"errors"

	"github.com/habbazettt/nutrisnap-server/internal/dto"
	"github.com/habbazettt/nutrisnap-server/internal/models"
	"github.com/habbazettt/nutrisnap-server/internal/repositories"
	"github.com/habbazettt/nutrisnap-server/pkg/constants"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
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
