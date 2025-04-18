package dto

import "github.com/raiymb/mappy/internal/map/model"

// MapPointResponse is sent to the frontend.
type MapPointResponse struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Title     string  `json:"title"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	StartYear int     `json:"startYear"`
	EndYear   int     `json:"endYear"`
}

// ToDTO converts a model to transport‑ready struct.
func ToDTO(mp model.MapPoint) MapPointResponse {
	return MapPointResponse{
		ID:        mp.ID.Hex(),
		Type:      string(mp.Type),
		Title:     mp.Title,
		X:         mp.X,
		Y:         mp.Y,
		StartYear: mp.StartYear,
		EndYear:   mp.EndYear,
	}
}
