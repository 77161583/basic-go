.PHONY: mock
mock:
	@mockgen -source=mybook/internal/service/user.go -package=svcmocks -destination=mybook/internal/service/mocks/user.mock.go
	@mockgen -source=mybook/internal/service/code.go -package=svcmocks -destination=mybook/internal/service/mocks/code.mock.go
	@go mod tidy