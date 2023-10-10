package repository

import (
	"context"
	"database/sql"
	"mybook/internal/domain"
	"mybook/internal/repository/cache"
	"mybook/internal/repository/dao"
	"time"
)

var ErrUseDuplicate = dao.ErrUseDuplicate
var ErrUserNotFund = dao.ErrUserNotFund

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *UserRepository) Created(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.DomainToEntity(u))
}

func (r *UserRepository) Edit(ctx context.Context, u domain.User) error {
	return r.dao.Edit(ctx, dao.User{
		Id:           u.Id,
		NickName:     u.NickName,
		Birthday:     u.Birthday,
		Introduction: u.Introduction,
	})
}

func (r *UserRepository) FindById(ctx context.Context, userId int64) (domain.User, error) {
	u, err := r.dao.FindById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	//t1 := u.CreateTime / 1000
	//t2 := u.UpdateTime / 1000
	return r.entityToDomain(u), nil
}

func (r *UserRepository) CacheFindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}
	//没有这个数据
	if err == cache.ErrKeyNotExist {
		//去数据里面加载
	}

	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = r.entityToDomain(ue)

	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			//这里打日志做监控
			//return
		}
	}()

	return u, nil

	//这里怎么办？ err = io.Err
	//要不要去数据库加载？
	//看起来我不应该加载？
	//看起来我好像也要加载

	//选加载 -- 做好兜底，万一 Redis 真的崩了，要保护住你的数据库
	// 数据库限流

	// 选不加载 -- 用户体验差一些

	// 缓存里面有数据
	// 缓存里面没有数据
	// 缓存出错了

}

func (r *UserRepository) DomainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password:   u.Password,
		CreateTime: u.CreateTime.UnixMilli(),
	}
}

func (r *UserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:           u.Id,
		Email:        u.Email.String,
		Phone:        u.Phone.String,
		Password:     u.Password,
		NickName:     u.NickName,
		Birthday:     u.Birthday,
		Introduction: u.Introduction,
		CreateTime:   time.UnixMilli(u.CreateTime),
		UpdateTime:   time.UnixMilli(u.UpdateTime),
	}
}
