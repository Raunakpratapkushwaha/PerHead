package service

import (
	"context"
	"errors"

	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/config"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/model"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/repository"
	"github.com/Raunakpratapkushwaha/Batwara/backend/pkg/password"
	"github.com/Raunakpratapkushwaha/Batwara/backend/pkg/token"
)

type AuthService struct {
	repo       repository.UserRepository
	tokenMaker *token.TokenMaker
	config     *config.Config
}

func NewAuthService(repo repository.UserRepository, maker *token.TokenMaker, cfg *config.Config) *AuthService {
	return &AuthService{
		repo:       repo,
		tokenMaker: maker,
		config:     cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.TokenResponse, error) {
	params := password.Argon2Params{
		Memory:      s.config.Argon2Memory,
		Iterations:  s.config.Argon2Iterations,
		Parallelism: s.config.Argon2Parallelism,
		SaltLength:  s.config.Argon2SaltLength,
		KeyLength:   s.config.Argon2KeyLength,
	}

	hashedPassword, err := password.HashPassword(req.Password, &params)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return s.generateUserTokens(user)
}

func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.TokenResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	valid, err := password.CheckPasswordHash(req.Password, user.PasswordHash)
	if err != nil || !valid {
		return nil, errors.New("invalid email or password")
	}

	return s.generateUserTokens(user)
}

func (s *AuthService) RefreshToken(ctx context.Context, tokenStr string) (*model.TokenResponse, error) {
	claims, err := s.tokenMaker.VerifyToken(tokenStr, true)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.repo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	return s.generateUserTokens(user)
}

func (s *AuthService) generateUserTokens(user *model.User) (*model.TokenResponse, error) {
	accessToken, err := s.tokenMaker.CreateToken(user.ID, s.config.AccessTokenDuration, false)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenMaker.CreateToken(user.ID, s.config.RefreshTokenDuration, true)
	if err != nil {
		return nil, err
	}

	return &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}