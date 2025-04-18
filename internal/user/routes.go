package user

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/raiymb/mappy/config"
	"github.com/raiymb/mappy/internal/middleware"
	"github.com/raiymb/mappy/internal/token"
	userhandler "github.com/raiymb/mappy/internal/user/handler"
	"github.com/raiymb/mappy/internal/user/repository"
	usersvc "github.com/raiymb/mappy/internal/user/service"
)

// Mount registers /auth/* and /users routes onto a sub‑router.
// Call from internal/app/routes.go   e.g.  user.RegisterRouter(apiGroup, …)
func RegisterRouter(r *gin.RouterGroup,
	repo repository.UserRepository,
	jwtCfg config.JWT,
	redisCli *redis.Client,
) {
	bl := token.NewBlacklist(redisCli)

	authSvc := usersvc.NewAuthService(repo, jwtCfg, bl)
	profSvc := usersvc.NewProfileService(repo)

	authH := userhandler.NewAuthHandler(authSvc)
	profH := userhandler.NewProfileHandler(profSvc)

	/* ---------- PUBLIC AUTH ROUTES ---------- */
	r.POST("/auth/register", authH.Register)
	r.POST("/auth/login", authH.Login)
	r.POST("/auth/refresh", authH.Refresh)
	r.POST("/auth/logout", authH.Logout)

	/* ---------- PROTECTED PROFILE ROUTES ---------- */
	users := r.Group("/users", middleware.Auth(jwtCfg.Secret, bl, "user"))
	{
		users.GET("/me", profH.Me)
	}
}
