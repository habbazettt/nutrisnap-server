package services

import (
	"github.com/habbazettt/nutrisnap-server/internal/dto"
	"github.com/habbazettt/nutrisnap-server/internal/models"
	"github.com/habbazettt/nutrisnap-server/internal/repositories"
)

type AdminService interface {
	GetAllUsers(page, limit int) ([]models.User, int64, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUserRole(userID string, role models.UserRole) (*models.User, error)
	DeleteUser(userID string) error
	GetStats() (*dto.AdminStatsResponse, error)
}

type adminService struct {
	userRepo repositories.UserRepository
}

func NewAdminService(userRepo repositories.UserRepository) AdminService {
	return &adminService{userRepo: userRepo}
}

func (s *adminService) GetAllUsers(page, limit int) ([]models.User, int64, error) {
	offset := (page - 1) * limit
	return s.userRepo.FindAll(offset, limit)
}

func (s *adminService) GetUserByID(id string) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *adminService) UpdateUserRole(userID string, role models.UserRole) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	user.Role = role
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *adminService) DeleteUser(userID string) error {
	return s.userRepo.Delete(userID)
}

func (s *adminService) GetStats() (*dto.AdminStatsResponse, error) {
	totalUsers, err := s.userRepo.Count()
	if err != nil {
		return nil, err
	}

	return &dto.AdminStatsResponse{
		TotalUsers: totalUsers,
	}, nil
}
