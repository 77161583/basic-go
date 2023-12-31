// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"basic-go/mybook/internal/repository"
	"basic-go/mybook/internal/repository/cache"
	"basic-go/mybook/internal/repository/dao"
	"basic-go/mybook/internal/service"
	"basic-go/mybook/internal/web"
	"basic-go/mybook/ioc"
	"github.com/gin-gonic/gin"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := ioc.InitRedis()
	v := ioc.InitMiddleware(cmdable)
	db := InitDB()
	userDAO := dao.NewUserDao(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	userServicePackage := service.NewUserService(userRepository)
	codeCache := cache.NewLocalCodeCache()
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMSService()
	codeServicePackage := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userServicePackage, codeServicePackage)
	engine := ioc.InitGin(v, userHandler)
	return engine
}
