package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	testCase := []struct {
		name string
		//这里为什么不用 ctrl？
		//因为这里是sqlmock，不是 gomock
		mock   func(t *testing.T) *sql.DB
		ctx    context.Context
		user   User
		wanErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				res := sqlmock.NewResult(3, 1)
				//这边预期的是正则表达式
				//这个写法的意思是，只要 INSERT 到 users 的语句
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnResult(res)
				require.NoError(t, err)
				return mockDB
			},
			user: User{Email: sql.NullString{
				String: "123@qq.com",
				Valid:  true,
			}},
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				//这边预期的是正则表达式
				//这个写法的意思是，只要 INSERT 到 users 的语句
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnError(&mysql.MySQLError{
					Number: 1062,
				})
				require.NoError(t, err)
				return mockDB
			},
			user:   User{},
			wanErr: ErrUseDuplicate,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				//这边预期的是正则表达式
				//这个写法的意思是，只要 INSERT 到 users 的语句
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnError(errors.New("数据库错误"))
				require.NoError(t, err)
				return mockDB
			},
			user:   User{},
			wanErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(gormMysql.New(gormMysql.Config{
				Conn:                      tc.mock(t),
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			d := NewUserDao(db)
			u := tc.user
			err = d.Insert(tc.ctx, u)
			assert.Equal(t, tc.wanErr, err)
		})
	}
}
