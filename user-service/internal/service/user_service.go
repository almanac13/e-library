package service

import (
	"errors"
	"strings"
	"user-service/internal/model"
	"user-service/internal/repository"

	"github.com/google/uuid"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(name, email, password, role string) (*model.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))
	role = strings.TrimSpace(strings.ToLower(role))

	if name == "" || email == "" || password == "" {
		return nil, errors.New("name, email and password are required")
	}

	if role == "" {
		role = "student"
	}

	if role != "student" && role != "admin" {
		return nil, errors.New("role must be student or admin")
	}

	user := model.User{
		ID:       uuid.New().String(),
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
	}

	err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return &user, nil
}

func (s *UserService) Login(email, password string) (*model.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.Password != password {
		return nil, errors.New("invalid email or password")
	}

	user.Password = ""
	return user, nil
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("user id is required")
	}

	return s.repo.GetUserByID(id)
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.repo.GetAllUsers()
}
