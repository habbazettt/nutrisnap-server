package services

import (
	"errors"

	"github.com/habbazettt/nutrisnap-server/internal/dto"
	"github.com/habbazettt/nutrisnap-server/internal/models"
	"github.com/habbazettt/nutrisnap-server/internal/repositories"
)

var (
	ErrPasswordMismatch = errors.New("current password is incorrect")
)

type UserService interface {
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	UpdateProfile(userID string, req *dto.UpdateProfileRequest) (*models.User, error)
	ChangePassword(userID string, req *dto.ChangePasswordRequest) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetByID(id string) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) GetByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

func (s *userService) UpdateProfile(userID string, req *dto.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	// Save changes
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) ChangePassword(userID string, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if !user.CheckPassword(req.CurrentPassword) {
		return ErrPasswordMismatch
	}

	// Set new password
	if err := user.SetPassword(req.NewPassword); err != nil {
		return err
	}

	// Save changes
	return s.userRepo.Update(user)
}
