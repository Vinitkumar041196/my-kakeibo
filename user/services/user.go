package services

import (
	"kakeibo-api/domain"
)

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) domain.UserService {
	return &userService{userRepo: userRepo}
}

func (us *userService) CreateUser(u *domain.User) error {
	return us.userRepo.CreateUser(u)
}

func (us *userService) UpdateUser(u *domain.User) error {
	return us.userRepo.UpdateUser(u)
}

func (us *userService) GetUser(email string) (*domain.User, error) {
	return us.userRepo.GetUser(email)
}

func (us *userService) DeleteUser(email string) error {
	return us.userRepo.DeleteUser(email)
}
