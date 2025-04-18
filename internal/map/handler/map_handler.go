package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raiymb/mappy/internal/map/service"
)

// MapHandler exposes REST endpoints for map data.
type MapHandler struct {
	svc service.Service
}

func New(svc service.Service) *MapHandler { return &MapHandler{svc: svc} }

/* ---------- PUBLIC ---------- */

// GET /map-points?year=1700&type=event,figure
func (h *MapHandler) ListPoints(c *gin.Context) {
	yearStr := c.DefaultQuery("year", "")
	if yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year query param required"})
		return
	}
	year, _ := strconv.Atoi(yearStr)
	pointTypes := service.ParseTypes(c.Query("type")) // optional

	items, err := h.svc.Points(c, year, pointTypes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

/* --------- ADMIN CRUD (examples) ---------- */

// You already have repo CRUD; here are wrappers.
//
// POST /admin/map-points – create
// PUT  /admin/map-points/:id – update
// DELETE /admin/map-points/:id – delete
//
// Implement these if/when you wire the admin routes.
