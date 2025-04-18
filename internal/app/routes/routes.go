package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raiymb/mappy/config"
	"github.com/raiymb/mappy/internal/map"
	maproutes "github.com/raiymb/mappy/internal/map"
	"github.com/raiymb/mappy/internal/middleware"
	"github.com/raiymb/mappy/internal/token"
	user "github.com/raiymb/mappy/internal/user"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

// Register wires *all* domain routers under /api.
// Call this once from cmd/server/main.go so that you don’t repeat wiring logic.
func Register(
	engine *gin.Engine,
	cfg config.Config,
	mongoDB *mongo.Database,
	redisCli *redis.Client,
	bl *token.Blacklist,
) {
	// top‑level middleware already added in main.go (logger, recovery, rate‑limit)
	api := engine.Group("/api")

	// ----- user/auth -----
	user.RegisterRouter(api, mongoDB, cfg.JWT, redisCli)

	// ----- timeline map -----
	maproutes.Register(api,
		mongoDB,
		redisCli,
		cfg.JWT.Secret,
		bl)

	// (Future) more bounded‑contexts:
	// article.RegisterRouter(api, …)
	// quiz.RegisterRouter(api, …)

	// A simple liveliness endpoint lives outside /api for K8s probes
	engine.GET("/healthz", func(c *gin.Context) { c.String(200, "ok") })
}
