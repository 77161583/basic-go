//go:build wireinject

//让wire 来注入这里的代码

package wire

import (
	"basic-go/wire/repository"
	"basic-go/wire/repository/dao"
	"github.com/google/wire"
)

func InitRepository() *repository.UserRepository {
	// 我只是在这里声明我要的各种东西，但是具体怎么构造，怎么编排顺序不管
	// 这个方法传入各个组件的初始化方法
	wire.Build(InitDB, repository.NewUserRepository, dao.NewUserDao)
	return new(repository.UserRepository)

}
