package web

import (
	"basic-go/mybook/internal/domain"
	"basic-go/mybook/internal/service"
	svcmocks "basic-go/mybook/internal/service/mocks"
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncrypt(t *testing.T) {
	password := "hello#world123"
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	//比较
	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
	assert.NoError(t, err)
}

func TestUserHandler_SignUp(t *testing.T) {
	// 这是一个测试用例切片，包含了多个测试场景
	testCases := []struct {
		/**
			name: 用于描述测试场景的名称。
			mock: 一个函数，用于创建模拟对象，模拟用户服务的行为。
			reqBody: 包含 JSON 数据的字符串，模拟 HTTP 请求体。
			wantCode: 预期的 HTTP 响应状态码。
			wantBody: 预期的 HTTP 响应正文。
		**/
		name     string
		mock     func(ctrl *gomock.Controller) service.UserServicePackage
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserServicePackage {
				usersvc := svcmocks.NewMockUserServicePackage(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "124@qq.com",
					Password: "Qq@adm331",
				}).Return(nil)
				//注册成功是 return nil
				return usersvc
			},
			reqBody: `{
   "email":"124@qq.com",
   "password":"Qq@adm331",
   "confirmPassword":"Qq@adm331"
}`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数不对， bind 失败",
			mock: func(ctrl *gomock.Controller) service.UserServicePackage {
				usersvc := svcmocks.NewMockUserServicePackage(ctrl)
				return usersvc
			},
			reqBody: `{
   "email":"124@qq.com",
   "password":"Qq@adm331",
   "confirmPassword":"Qq@adm331"
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserServicePackage {
				usersvc := svcmocks.NewMockUserServicePackage(ctrl)
				return usersvc
			},
			reqBody: `{
   "email":"124qq.com",
   "password":"Qq@adm331",
   "confirmPassword":"Qq@adm331"
}`,
			wantCode: http.StatusOK,
			wantBody: "Email格式不正确",
		},
		{
			name: "两次输入密码不匹配",
			mock: func(ctrl *gomock.Controller) service.UserServicePackage {
				usersvc := svcmocks.NewMockUserServicePackage(ctrl)
				return usersvc
			},
			reqBody: `{
   "email":"124@qq.com",
   "password":"Qq@adm3311",
   "confirmPassword":"Qq@adm331"
}`,
			wantCode: http.StatusOK,
			wantBody: "两次密码输入不一样！",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) service.UserServicePackage {
				usersvc := svcmocks.NewMockUserServicePackage(ctrl)
				return usersvc
			},
			reqBody: `{
   "email":"124@qq.com",
   "password":"Qq1234",
   "confirmPassword":"Qq1234"
}`,
			wantCode: http.StatusOK,
			wantBody: "密码必须大于8位，包含数字，特殊字符",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//使用 Gomock 创建一个控制器 ctrl，用于管理模拟对象。
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			//创建一个 GIN 引擎实例 server，用于处理 HTTP 请求
			server := gin.Default()
			// 用不上 codesvc
			//创建用户处理程序 h，并为其提供模拟用户服务和模拟验证码服务
			h := NewUserHandler(tc.mock(ctrl), nil)
			//创建 HTTP 请求 req，模拟用户注册请求，包括 URL 路径和 JSON 数据
			h.RegisterRoutes(server)
			//使用 httptest.NewRecorder() 创建一个 HTTP 响应记录器 resp，以捕获处理程序的响应。
			req, err := http.NewRequest(http.MethodPost,
				"http://localhost:8080/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))

			require.NoError(t, err)
			t.Log(req)
			//这里继续使用req
			resp := httptest.NewRecorder()
			t.Log(resp)
			//设置请求头的 Content-Type 为 JSON。数据格式是json
			req.Header.Set("Content-Type", "application/json")
			// 这就是 HTTP 请求进去 GIN 框架的入口
			// 当你这样调用的时候，GIN 就会处理这个请求
			// 响应写到 resp 里
			server.ServeHTTP(resp, req)
			//使用 assert 断言库检查实际的 HTTP 响应状态码和响应正文与预期是否一致。
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())

		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersvc := svcmocks.NewMockUserServicePackage(ctrl)

	//预期
	//usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
	//	Email: "112@qq.com",
	//}).Return(errors.New("mock error"))

	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).
		Return(errors.New("mock error"))

	err := usersvc.SignUp(context.Background(), domain.User{
		Email: "124@qq.com",
	})
	t.Log(err)
}
