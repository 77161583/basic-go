package main

import (
	"basic-go/mybook/config"
	"basic-go/mybook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		//只会在初始化的过程中panic
		//panic相当于整个goroutine结束
		//一旦初始化出错，应用就不要再启动了
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
