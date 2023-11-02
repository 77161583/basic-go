package repository

import (
	"basic-go/mybook/internal/domain"
	"basic-go/mybook/internal/repository/cache"
	cachemocks "basic-go/mybook/internal/repository/cache/mocks"
	"basic-go/mybook/internal/repository/dao"
	daomocks "basic-go/mybook/internal/repository/dao/mocks"
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	now := time.Now()
	//去除毫秒
	now = time.UnixMilli(now.UnixMilli())
	testCase := []struct {
		name string

		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		ctx      context.Context
		id       int64
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				//缓存未命中，差了缓存，但是没结果
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)

				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.User{
						Id: 123,
						Email: sql.NullString{
							String: "123@qq.com",
							Valid:  true,
						},
						Password: "123123",
						Phone: sql.NullString{
							String: "13511111111",
							Valid:  true,
						},
						NickName:   "",
						CreateTime: now.UnixMilli(),
						UpdateTime: now.UnixMilli(),
					}, nil)

				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:         123,
					Email:      "123@qq.com",
					Password:   "123123",
					Phone:      "13511111111",
					CreateTime: now,
					UpdateTime: now,
				}).Return(nil)
				return d, c
			},

			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:         123,
				Email:      "123@qq.com",
				Password:   "123123",
				Phone:      "13511111111",
				CreateTime: now,
				UpdateTime: now,
			},
			wantErr: nil,
		},
		{
			name: "缓存未命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				//缓存未命中，差了缓存，但是没结果
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{
						Id:         123,
						Email:      "123@qq.com",
						Password:   "123123",
						Phone:      "13511111111",
						CreateTime: now,
						UpdateTime: now,
					}, nil)

				d := daomocks.NewMockUserDAO(ctrl)
				return d, c
			},

			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:         123,
				Email:      "123@qq.com",
				Password:   "123123",
				Phone:      "13511111111",
				CreateTime: now,
				UpdateTime: now,
			},
			wantErr: nil,
		},
		{
			name: "查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				//缓存未命中，差了缓存，但是没结果
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)

				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.User{}, errors.New("mock db 错误"))

				return d, c
			},

			ctx:      context.Background(),
			id:       123,
			wantUser: domain.User{},
			wantErr:  errors.New("mock db 错误"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tc.mock(ctrl)
			repo := NewUserRepository(ud, uc)
			u, err := repo.FindById(tc.ctx, tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
			//测goRouting 可以用这个，但是不保证每次都会触发
			time.Sleep(time.Second)

		})
	}
}
