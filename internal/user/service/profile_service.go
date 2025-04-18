package service

import (
	"context"

	"github.com/raiymb/mappy/internal/internal/user/dto"
	"github.com/raiymb/mappy/internal/internal/user/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileService struct {
	repo repository.UserRepository
}

func NewProfileService(repo repository.UserRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) Get(ctx context.Context, uid string) (*dto.ProfileResponse, error) {
	user, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	return &dto.ProfileResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Name:      user.Name,
		Role:      string(user.Role),
		AvatarURL: user.AvatarURL,
	}, nil
}
