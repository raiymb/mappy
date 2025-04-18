package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raiymb/mappy/internal/token"
	"github.com/raiymb/mappy/internal/user/dto"
	"github.com/raiymb/mappy/internal/user/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pair, err := h.svc.Register(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pair)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pair, err := h.svc.Login(c, req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pair)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var body struct {
		Refresh string `json:"refresh"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Refresh == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
		return
	}
	newPair, err := h.svc.Refresh(c, body.Refresh)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newPair)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var body struct {
		Refresh string `json:"refresh"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Refresh == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
		return
	}
	if err := h.svc.Logout(c, body.Refresh); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
