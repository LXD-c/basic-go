//go:build wireinject

package main

import (
	"github.com/LXD-c/basic-go/webook/internal/repository"
	"github.com/LXD-c/basic-go/webook/internal/repository/article"
	"github.com/LXD-c/basic-go/webook/internal/repository/cache"
	"github.com/LXD-c/basic-go/webook/internal/repository/dao"
	article2 "github.com/LXD-c/basic-go/webook/internal/repository/dao/article"
	"github.com/LXD-c/basic-go/webook/internal/service"
	"github.com/LXD-c/basic-go/webook/internal/web"
	"github.com/LXD-c/basic-go/webook/internal/web/jwt"
	"github.com/LXD-c/basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		//最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,

		//初始化 DAO
		dao.NewUserDAO,
		article2.NewGORMArticleDAO,

		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,
		article.NewArticleRepository,
		//article.NewArticleReaderRepository,
		//article.NewArticleAuthorRepository,

		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,

		// 直接基于内存实现
		ioc.InitSMSService,
		ioc.InitWechatService,

		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		ioc.NewWechatHandlerConfig,
		jwt.NewRedisJWTHandler,
		web.NewArticleHandler,
		// 你中间件呢？
		// 你注册路由呢？
		// 你这个地方没有用到前面的任何东西
		// gin.Default,
		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
