package instances

import "basic-go/mybook/internal/repository"

var userRepo repository.UserRepository

func InitUserRepo(repo repository.UserRepository) {
	userRepo = repo
}
