package service

import (
	"basic-go/mybook/internal/domain"
	"basic-go/mybook/internal/repository"
	repomocks "basic-go/mybook/internal/repository/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestUserService_Login(t *testing.T) {
	//公共时间
	now := time.Now()
	testCase := []struct {
		name string
		mock func(repo *gomock.Controller) repository.UserRepository

		//输入
		email    string
		password string

		//输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:      "123@qq.com",
						Password:   "$2a$10$0X2rBtMoM3moeaboBTZS3.C1gnnTWHZCqAT1fJH0yEtosnWprNH7S",
						Phone:      "13511111111",
						CreateTime: now,
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "Qq@adm331",
			wantUser: domain.User{
				Email:      "123@qq.com",
				Password:   "$2a$10$0X2rBtMoM3moeaboBTZS3.C1gnnTWHZCqAT1fJH0yEtosnWprNH7S",
				Phone:      "13511111111",
				CreateTime: now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFund)
				return repo
			},
			email:    "123@qq.com",
			password: "Qq@adm331",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("mock db 错误"))
				return repo
			},
			email:    "123@qq.com",
			password: "Qq@adm331",
			wantUser: domain.User{},
			wantErr:  errors.New("mock db 错误"),
		},
		{
			name: "密码不一致",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, ErrInvalidUserOrPassword)
				return repo
			},
			email:    "123@qq.com",
			password: "Qq@adm33111",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			//具体测试代码
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl))
			u, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

func TestEncrypted(t *testing.T) {
	res, error := bcrypt.GenerateFromPassword([]byte("Qq@adm331"), bcrypt.DefaultCost)
	if error == nil {
		t.Log(string(res))
	}
}
