package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raiymb/mappy/config"
	"github.com/raiymb/mappy/internal/app/routes"    // root router
	maproutes "github.com/raiymb/mappy/internal/map"
	"github.com/raiymb/mappy/internal/middleware"
	"github.com/raiymb/mappy/internal/token"
	user "github.com/raiymb/mappy/internal/user"
	"github.com/raiymb/mappy/pkg/logger"
	"github.com/raiymb/mappy/storage"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	 "github.com/ilyakaznacheev/cleanenv"
)

func main() {
	/* ---------- 1. CONFIG ---------- */
	var cfg config.Config
	if err := cleanenv.ReadConfig("config/config.yaml", &cfg); err != nil {
		logger.Init(true)
		logger.L().Fatal("config load", zap.Error(err))
	}
	_ = cleanenv.ReadConfig("config/local.yaml", &cfg) // override for dev

	logger.Init(cfg.Server.Env == "dev")
	log := logger.L()

	/* ---------- 2. CONTEXT ---------- */
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	/* ---------- 3. STORAGE ---------- */
	mongoCli, err := storage.NewMongo(ctx, cfg.Mongo.URI)
	if err != nil {
		log.Fatal("mongo connect", zap.Error(err))
	}
	defer mongoCli.Disconnect(ctx)

	redisCli, err := storage.NewRedis(ctx, cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatal("redis connect", zap.Error(err))
	}
	defer redisCli.Close()

	bl := token.NewBlacklist(redisCli)

	/* ---------- 4. ROUTER + MIDDLEWARE ---------- */
	r := gin.New()
	if cfg.Server.Env == "dev" {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery(),
		middleware.RateLimiter(redisCli, cfg.RateLimit.Window, cfg.RateLimit.Max))

	// Health‑probe
	r.GET("/healthz", func(c *gin.Context) { c.String(200, "ok") })

	api := r.Group("/api")

	// User & Auth
	user.RegisterRouter(api, mongoCli.Database(cfg.Mongo.DB), cfg.JWT, redisCli)

	// Map Timeline
	maproutes.Register(api, mongoCli.Database(cfg.Mongo.DB), redisCli, cfg.JWT.Secret, bl)

	/* ---------- 5. HTTP SERVER ---------- */
	srv := &http.Server{
		Addr:         cfg.Server.ListenAddr(),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	go func() {
		log.Info("server started", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen", zap.Error(err))
		}
	}()

	/* ---------- 6. GRACEFUL SHUTDOWN ---------- */
	<-ctx.Done()
	stop()
	log.Info("shutdown signal received")

	shCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shCtx)
	log.Info("server exited")
}
