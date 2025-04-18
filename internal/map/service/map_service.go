package service

import (
	"context"
	"slices"
	"strings"

	"github.com/raiymb/mappy/internal/map/dto"
	"github.com/raiymb/mappy/internal/map/model"
	"github.com/raiymb/mappy/internal/map/repository"
)

type Service interface {
	// Points returns markers for a year + optional type filter.
	Points(ctx context.Context, year int, types []model.PointType) ([]dto.MapPointResponse, error)
}

type mapService struct {
	repo repository.Repository
}

// New creates a timeline / layer‑filter service.
func New(repo repository.Repository) Service {
	return &mapService{repo: repo}
}

func (s *mapService) Points(ctx context.Context, year int, types []model.PointType) ([]dto.MapPointResponse, error) {
	raw, err := s.repo.PointsByYear(ctx, year)
	if err != nil {
		return nil, err
	}

	// optional layer toggle: keep only requested types
	if len(types) > 0 {
		raw = slices.DeleteFunc(raw, func(mp model.MapPoint) bool {
			return !contains(types, mp.Type)
		})
	}

	out := make([]dto.MapPointResponse, len(raw))
	for i, mp := range raw {
		out[i] = dto.ToDTO(mp)
	}
	return out, nil
}

/* ---------- helpers ---------- */

func contains(arr []model.PointType, t model.PointType) bool {
	for _, x := range arr {
		if x == t {
			return true
		}
	}
	return false
}

// ParseTypes converts CSV query string → []PointType, ignoring bad values.
func ParseTypes(csv string) []model.PointType {
	var out []model.PointType
	for _, s := range strings.Split(csv, ",") {
		switch strings.TrimSpace(strings.ToLower(s)) {
		case "event":
			out = append(out, model.TypeEvent)
		case "figure":
			out = append(out, model.TypeFigure)
		case "place":
			out = append(out, model.TypePlace)
		}
	}
	return out
}
