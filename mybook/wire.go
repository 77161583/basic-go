//go:build wireinject

package main

import (
	"basic-go/mybook/internal/repository"
	"basic-go/mybook/internal/repository/cache"
	"basic-go/mybook/internal/repository/dao"
	"basic-go/mybook/internal/service"
	"basic-go/mybook/internal/web"
	"basic-go/mybook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		//最基础的第三方依赖
		InitDB, ioc.InitRedis,
		//初始化 dao
		dao.NewUserDao,
		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewUserService,
		service.NewCodeService,
		//基于内存实现
		ioc.InitSMSService,
		web.NewUserHandler,
		//
		ioc.InitGin,
		ioc.InitMiddleware,
	)
	return new(gin.Engine)
}
