package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raiymb/mappy/internal/user/service"
)

type ProfileHandler struct {
	svc *service.ProfileService
}

func NewProfileHandler(svc *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

func (h *ProfileHandler) Me(c *gin.Context) {
	uid, _ := c.Get("uid") // set by auth middleware
	resp, err := h.svc.Get(c, uid.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if resp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
