package services

import (
	"e-wallet-go/internal/models"
	"e-wallet-go/internal/repository"
	"errors"
)

type UserService interface {
	GetAllUsers(requesterRole string) ([]models.User, error)
	GetUserByID(requesterID, requesterRole, targetID string) (*models.User, error)
	UpdateUser(requesterID, requesterRole, targetID string, updateData models.User) (*models.User, error)
	DeleteUser(requesterID, requesterRole, targetID string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetAllUsers(requesterRole string) ([]models.User, error) {
	if requesterRole != "admin" {
		return nil, errors.New("unauthorized: only admin can view all users")
	}
	return s.userRepo.FindAll()
}

func (s *userService) GetUserByID(requesterID, requesterRole, targetID string) (*models.User, error) {
	if requesterRole != "admin" && requesterID != targetID {
		return nil, errors.New("forbidden: you can only view your own profile")
	}
	return s.userRepo.FindByID(targetID)
}

func (s *userService) UpdateUser(requesterID, requesterRole, targetID string, updateData models.User) (*models.User, error) {
	if requesterRole != "admin" && requesterID != targetID {
		return nil, errors.New("forbidden: you can only update your own profile")
	}

	user, err := s.userRepo.FindByID(targetID)
	if err != nil {
		return nil, err
	}

	if updateData.Username != "" {
		user.Username = updateData.Username
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(requesterID, requesterRole, targetID string) error {
	if requesterRole != "admin" && requesterID != targetID {
		return errors.New("forbidden: you can only delete your own account")
	}
	return s.userRepo.Delete(targetID)
}