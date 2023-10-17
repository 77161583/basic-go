package repository

import (
	"basic-go/mybook/internal/repository/cache"
	"context"
)

var (
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
	ErrCodeInvalid            = cache.ErrCodeInvalid
	ErrCodeTimeOut            = cache.ErrCodeTimeOut
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CacheCodeRepository struct {
	//cache cache.CodeRedisCache
	cache cache.CodeCache
}

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		cache: c,
	}
}

func (repo *CacheCodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.LocalSet(ctx, biz, phone, code)
}

func (repo *CacheCodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return repo.cache.LocalVerify(ctx, biz, phone, inputCode)
}
