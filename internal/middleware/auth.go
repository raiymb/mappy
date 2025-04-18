package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raiymb/mappy/internal/token"
	"github.com/raiymb/mappy/pkg/logger"
)

// Auth returns a Gin middleware that:
//   1. extracts `Authorization: Bearer <jwt>`,
//   2. validates signature / expiry,
//   3. checks Redis blacklist (optional),
//   4. enforces the minimal allowed role.
// It then stores the claims in ctx for downstream handlers.
func Auth(secret string, bl *token.Blacklist, minRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		raw := strings.TrimPrefix(h, "Bearer ")
		claims, err := token.Parse(raw, secret)
		if err != nil {
			logger.L().Warn("token parse failed", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// blacklist check (skip if nil)
		if bl != nil {
			if ok, _ := bl.IsBlacklisted(c, claims.ID); ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
				return
			}
		}

		// role hierarchy: guest < user < moderator < admin
		const (
			guest = iota
			user
			mod
			admin
		)
		rank := map[string]int{
			"guest":     guest,
			"user":      user,
			"moderator": mod,
			"admin":     admin,
		}
		if rank[claims.Role] < rank[minRole] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}

		// attach claims
		c.Set("uid", claims.UID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
