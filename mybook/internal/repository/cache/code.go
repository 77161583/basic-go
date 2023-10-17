package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	_ "github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknowForCode          = errors.New("我也不知道发生了什么，反正是跟 code 有关")
	ErrCodeInvalid            = errors.New("验证码失效，请重新获取验证码")
	ErrCodeTimeOut            = errors.New("c！")
)

// 编译器会在编译的时候，把 set_code 的代码放进来这个 luaSetCode 变量里

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeRedisCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeRedisCache {
	return &RedisCodeCache{
		client: client,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		//no problem
		return nil
	case -1:
		//发送频繁
		return ErrCodeSendTooMany
	//case -2:

	default:
		//系统错误
		return errors.New("系统错误")
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		//如果频繁出现这个错误，你需要告警，有人在搞你
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		return false, nil
	}
	return false, ErrUnknowForCode
}

func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phoneCode:%s:%s", biz, phone)
}

// CodeCache 第四次作业 开始 ——————————————————————————————————————————————————————————————————
type CodeCache interface {
	LocalSet(ctx context.Context, biz, phone, code string) error
	LocalVerify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type LocalCodeCache struct {
	//创建一个新的 sync.map
	localCache sync.Map
}

func NewLocalCodeCache() CodeCache {
	return &LocalCodeCache{
		localCache: sync.Map{},
	}
}

// 实现接口的对应方法
var mu sync.Mutex

func (c *LocalCodeCache) LocalSet(ctx context.Context, biz, phone, code string) error {
	//因为使用本地缓存sync.map 不需要在执行 lua 脚本
	key := c.LocalKey(biz, phone)
	cntKey := key + ":cnt"
	// 在本地缓存中存储验证码
	c.localCache.Store(key, code)
	// 在本地缓存中存储验证次数，一个验证码最多验证三次
	c.localCache.Store(cntKey, 3)
	// 设置验证码的过期时间为 1 分钟（
	expirationTime := 1 * time.Minute

	//好像会有并发问题，暂时没有想到好的解决办法
	mu.Lock()
	defer mu.Unlock()

	// 使用 time.AfterFunc 来在过期后删除验证码
	time.AfterFunc(expirationTime, func() {
		//删除时候，监测下key 是否存在
		if _, ok := c.localCache.Load(key); ok {
			c.localCache.Delete(key)
		}
		if _, ok := c.localCache.Load(cntKey); ok {
			c.localCache.Delete(cntKey)
		}

		//设置验证码是否过期
		c.localCache.Store(key+":timestamp", 1)
	})

	return nil
}

func (c *LocalCodeCache) LocalVerify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	key := c.LocalKey(biz, phone)
	cntKey := key + ":cnt"

	// 从本地缓存中获取验证码和验证次数
	code, codeExists := c.localCache.Load(key)
	cnt, cntExists := c.localCache.Load(cntKey)
	_, timeExists := c.localCache.Load(key + ":timestamp")
	// 查看时间戳是否存在
	if timeExists {
		// 如果存在说明，验证码超时了
		return false, ErrCodeTimeOut
	}

	// 检查验证码、验证次数是否存在
	if !codeExists || !cntExists {
		// 没有找到相应的数据，返回错误
		return false, errors.New("系统错误")
	}

	// 因为使用 sync.Map 存储值，它将所有值存储为 interface{} 类型，需要进行类型断言以获取存储的值的实际类型
	// 将验证码转换为字符串
	codeStr, ok := code.(string)
	if !ok {
		return false, errors.New("系统错误")
	}

	// 将验证次数转换为整数
	cntInt, ok := cnt.(int)
	if !ok {
		return false, errors.New("系统错误")
	}

	if cntInt <= 0 {
		// 用户一直输错或已经用过
		return false, ErrCodeInvalid
	}

	if inputCode == codeStr {
		// 输入正确，标记验证码为已用
		c.localCache.Store(cntKey, "-1")
		return true, nil
	}

	// 用户输入错误，减少验证次数
	cntInt--
	c.localCache.Store(cntKey, cntInt)
	return false, nil
}

func (c *LocalCodeCache) LocalKey(biz, phone string) string {
	return fmt.Sprintf("phoneCode:%s:%s", biz, phone)
}
