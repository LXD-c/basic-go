package ioc

import (
	"github.com/LXD-c/basic-go/webook/internal/web"
	ijwt "github.com/LXD-c/basic-go/webook/internal/web/jwt"
	"github.com/LXD-c/basic-go/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(middlewares []gin.HandlerFunc, u *web.UserHandler, w *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	//middlewares...将切片解包为多个独立的参数
	server.Use(middlewares...)
	u.RegisterRoutes(server)
	w.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable, hdl ijwt.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),

		middleware.NewLoginJWTMiddlewareBuilder(hdl).
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/users/login").Build(),

		//ratelimit.NewBuilder(redisClient,)
	}
}

func corsHdl() gin.HandlerFunc {
	//使用use方法注册middleware，这个中间件是用于解决 CORS 的 middleware
	return cors.New(cors.Config{
		//是否允许带上用户认证信息（比如 cookie）
		AllowCredentials: true,
		AllowHeaders:     []string{"content-type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
		//哪些来源是允许的
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "live.webook.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
