package maproutes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raiymb/mappy/internal/map/handler"
	"github.com/raiymb/mappy/internal/map/repository"
	"github.com/raiymb/mappy/internal/map/service"
	"github.com/raiymb/mappy/internal/middleware"
	"github.com/raiymb/mappy/internal/token"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

// Register attaches /map‑points routes (plus admin)
//  r  – parent group, e.g. api := router.Group("/api")
//  db – *mongo.Database
func Register(r *gin.RouterGroup,
	db *mongo.Database,
	redisCli *redis.Client,
	jwtSecret string,
	bl *token.Blacklist, // can be nil if you don’t use blacklist
) {
	// 1. repo with cache
	baseRepo := repository.NewMongo(db)
	cached    := repository.NewCached(baseRepo, redisCli, 6*time.Hour)

	// 2. wire service & handler
	svc := service.New(cached)
	h   := handler.New(svc)

	/* ---------- PUBLIC ---------- */
	r.GET("/map-points", h.ListPoints)

	/* ---------- ADMIN (optional) ---------- */
	admin := r.Group("/admin", middleware.Auth(jwtSecret, bl, "admin"))
	{
		// admin.POST("/map-points", h.CreatePoint)
		// admin.PUT("/map-points/:id", h.UpdatePoint)
		// admin.DELETE("/map-points/:id", h.DeletePoint)
	}
}
