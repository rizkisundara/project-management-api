package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rizkisundara/project-management-api/models"
	"github.com/rizkisundara/project-management-api/models/dto"
	"github.com/rizkisundara/project-management-api/repositories"
	"github.com/rizkisundara/project-management-api/utils"
)

type UserService interface {
	Register(req *dto.UserRegisterRequest) (*models.User, error)
	Login(email, password string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByPublicID(publicID string) (*models.User, error)
	GetAllUsersWithPagination(q utils.PaginationQuery) ([]models.User, int64, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(req *dto.UserRegisterRequest) (*models.User, error) {
	existingUser, _ := s.repo.FindByEmail(req.Email)
	if existingUser != nil && existingUser.InternalID != 0 {
		return nil, errors.New("email already registered")
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		PublicID: uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
		Role:     "user",
	}

	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userService) Login(email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("Invalid Credential")
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("Invalid Credential")
	}

	return user, nil
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) GetByPublicID(publicID string) (*models.User, error) {
	return s.repo.FindByPublicID(publicID)
}

func (s *userService) GetAllUsersWithPagination(q utils.PaginationQuery) ([]models.User, int64, error) {
	return s.repo.FindAllUsersWithPagination(q)
}
