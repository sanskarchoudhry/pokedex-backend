package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sanskarchoudhry/pokedex-backend/internal/models"
	"github.com/sanskarchoudhry/pokedex-backend/internal/repository"
	"github.com/sanskarchoudhry/pokedex-backend/internal/utils"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, string, error) // Returns (accessToken, refreshToken, error)
}

type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (s *authService) Register(ctx context.Context, email, password string) (*models.User, error) {

	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("checking existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("email already in use")
	}

	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	newUser := &models.User{Email: email, Password: hashedPwd}
	if err := s.userRepo.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	responseUser := *newUser
	responseUser.Password = ""
	return &responseUser, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := utils.CheckPassword(password, user.Password); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", "", fmt.Errorf("generating access token: %w", err)
	}

	rawRefreshToken, tokenHash, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("generating refresh token: %w", err)
	}

	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.tokenRepo.CreateRefreshToken(ctx, refreshTokenModel); err != nil {
		return "", "", fmt.Errorf("saving refresh token: %w", err)
	}

	return accessToken, rawRefreshToken, nil
}
