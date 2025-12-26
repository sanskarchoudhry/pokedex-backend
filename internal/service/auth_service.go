package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sanskarchoudhry/pokedex-backend/internal/models"
	"github.com/sanskarchoudhry/pokedex-backend/internal/repository"
	"github.com/sanskarchoudhry/pokedex-backend/internal/utils"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{
		userRepo: repo,
	}
}

func (s *authService) Register(ctx context.Context, email, password string) (*models.User, error) {

	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("email already in use")
	}

	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := &models.User{
		Email:    email,
		Password: hashedPwd,
	}

	if err := s.userRepo.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	// Sanitize response
	// We create a copy to return, removing sensitive info like the hash
	responseUser := *newUser
	responseUser.Password = ""

	return &responseUser, nil
}
