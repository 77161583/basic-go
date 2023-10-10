package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"mybook/internal/domain"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache struct {
	//传 单机 Redis 可以
	//传 cluster 的 Redis 也可以
	client     redis.Cmdable
	expiration time.Duration
}

// A 用到了 B，B 一定是接口 =》保证面向接口
// A 用到了 B，B 一定是 A 的字段 =》规避包变量，包方法，都非常缺乏扩展性
// A 用到了 B，A 绝对不初始化 B，而是外面注入 =》保持依赖注入（DI，Dependency Injection）和以来反转
func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

// 只要error 为 nil ，就认为缓存里面有数据
// 如果没有数据，返回一个特定的 error
func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.Key(id)
	//如果数据不存在，err = redis.Nil
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}
func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.Key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *UserCache) Key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

/***********优雅做法 start *************/
//type CacheV1 interface {
//	//正常是中间件去做
//	Get(ctx context.Context,key string)(any error)
//}
//
//type UserCache struct {
//	cache CacheV1
//}
//
//func (u *UserCache) GetUser(ctx context.Context, id int64)(domain.User,error){
//
//}
/***********优雅做法 end *************/
