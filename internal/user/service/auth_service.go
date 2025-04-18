package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/raiymb/mappy/config"
	"github.com/raiymb/mappy/internal/token"
	"github.com/raiymb/mappy/internal/user/dto"
	"github.com/raiymb/mappy/internal/user/model"
	"github.com/raiymb/mappy/internal/user/repository"
	"github.com/raiymb/mappy/pkg/logger"
	"github.com/raiymb/mappy/pkg/util"
	"go.uber.org/zap"
)

var (
	ErrEmailTaken      = errors.New("email already registered")
	ErrInvalidCreds    = errors.New("invalid email or password")
	ErrBadRefreshToken = errors.New("refresh token invalid")
)

// AuthService orchestrates registration/login/refresh flows.
type AuthService struct {
	repo     repository.UserRepository
	jwtCfg   config.JWT
	bl       *token.Blacklist // can be nil
	validate *validator.Validate
}

// NewAuthService DI constructor.
func NewAuthService(repo repository.UserRepository, jwtCfg config.JWT, bl *token.Blacklist) *AuthService {
	return &AuthService{
		repo:     repo,
		jwtCfg:   jwtCfg,
		bl:       bl,
		validate: validator.New(),
	}
}

// Register a new account and return token pair.
func (s *AuthService) Register(ctx context.Context, req dto.RegisterDTO) (token.Pair, error) {
	if err := s.validate.Struct(req); err != nil {
		return token.Pair{}, err
	}
	if u, _ := s.repo.FindByEmail(ctx, req.Email); u != nil {
		return token.Pair{}, ErrEmailTaken
	}

	hash, _ := util.HashPassword(req.Password)
	user := &model.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hash,
		Role:         model.RoleUser,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return token.Pair{}, err
	}
	return token.NewPair(user.ID.Hex(), string(user.Role), s.jwtCfg)
}

// Login verifies credentials and returns token pair.
func (s *AuthService) Login(ctx context.Context, req dto.LoginDTO) (token.Pair, error) {
	if err := s.validate.Struct(req); err != nil {
		return token.Pair{}, err
	}
	user, _ := s.repo.FindByEmail(ctx, req.Email)
	if user == nil || !util.CheckPassword(user.PasswordHash, req.Password) {
		return token.Pair{}, ErrInvalidCreds
	}
	return token.NewPair(user.ID.Hex(), string(user.Role), s.jwtCfg)
}

// Refresh rotates refresh tokens (blacklists the old ‑ optional).
func (s *AuthService) Refresh(ctx context.Context, refresh string) (token.Pair, error) {
	claims, err := token.Parse(refresh, s.jwtCfg.Secret)
	if err != nil {
		return token.Pair{}, ErrBadRefreshToken
	}
	// blacklist old refresh
	if s.bl != nil {
		ttl := time.Until(claims.ExpiresAt.Time)
		_ = s.bl.Revoke(ctx, claims.ID, ttl)
	}
	return token.NewPair(claims.UID, claims.Role, s.jwtCfg)
}

// Logout simply blacklists the provided refresh token.
func (s *AuthService) Logout(ctx context.Context, refresh string) error {
	if s.bl == nil {
		return nil // no blacklist configured
	}
	claims, err := token.Parse(refresh, s.jwtCfg.Secret)
	if err != nil {
		return ErrBadRefreshToken
	}
	ttl := time.Until(claims.ExpiresAt.Time)
	return s.bl.Revoke(ctx, claims.ID, ttl)
}
