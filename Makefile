.PHONY: mock
mock:
	@mockgen -source=mybook/internal/service/user.go -package=svcmocks -destination=mybook/internal/service/mocks/user.mock.go
	@mockgen -source=mybook/internal/service/code.go -package=svcmocks -destination=mybook/internal/service/mocks/code.mock.go
	@mockgen -source=mybook/internal/repository/user.go -package=repomocks -destination=mybook/internal/repository/mocks/user.mock.go
	@mockgen -source=mybook/internal/repository/code.go -package=repomocks -destination=mybook/internal/repository/mocks/code.mock.go
	@mockgen -source=mybook/internal/repository/dao/user.go -package=daomocks -destination=mybook/internal/repository/dao/mocks/user.mock.go
	@mockgen -source=mybook/internal/repository/cache/user.go -package=cachemocks -destination=mybook/internal/repository/cache/mocks/user.mock.go
	@mockgen -package=redismocks -destination=mybook/internal/repository/cache/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy