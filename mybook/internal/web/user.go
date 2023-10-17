package web

import (
	"basic-go/mybook/internal/domain"
	"basic-go/mybook/internal/service"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const biz = "login"

// 确保 UserHandler 上实现了 handler 接口
var _ handler = &UserHandler{}

// 这个更优雅
var _ handler = (*UserHandler)(nil)

// UserHandler 定义和跟用户有关的路由
type UserHandler struct {
	svc         service.UserServicePackage
	codeSvc     service.CodeServicePackage
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	birthdayExp *regexp.Regexp
	phoneExp    *regexp.Regexp
}

func NewUserHandler(svc service.UserServicePackage, codeSvc service.CodeServicePackage) *UserHandler {
	//校验参数
	const (
		emailRegExpPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
		passwordRegExpPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
		birthdayRegExpPattern = `^\d{4}-\d{2}-\d{2}$`
		phoneExpPattern       = `^1[3456789]\d{9}$`
	)
	emailRegExp := regexp.MustCompile(emailRegExpPattern, 0)
	passwordRegExp := regexp.MustCompile(passwordRegExpPattern, 0)
	birthdayRegExp := regexp.MustCompile(birthdayRegExpPattern, 0)
	phoneExp := regexp.MustCompile(phoneExpPattern, 0)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailRegExp,
		passwordExp: passwordRegExp,
		birthdayExp: birthdayRegExp,
		phoneExp:    phoneExp,
		codeSvc:     codeSvc,
	}
}

func (u *UserHandler) RegisterRoutes(serve *gin.Engine) {
	ug := serve.Group("/users")
	//ug.POST("login", u.Login)
	ug.POST("login", u.LoginJWT)
	ug.POST("signup", u.SignUp)
	ug.POST("edit", u.Edit)
	ug.POST("profile", u.ProfileJWT)
	//put “login/sms/code”发送验证码
	//put “login/sms/code” 校验验证码
	//put “login/sms/code”发送验证码
	ug.POST("login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("login_sms", u.LoginSMS)
}

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		switch err {
		case service.ErrCodeInvalid:
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "验证码失效，请重新获取验证码",
			})
		case service.ErrCodeTimeOut:
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "验证码已过期!",
			})
		default:
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "系统错误！",
			})
		}
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证有误！",
		})
		return
	}
	//有可能是新用户
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误 ！",
		})
		return
	}
	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误！",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 4,
		Msg:  "验证码校验通过",
	})
}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//这里可以校验下手机号，可以用正则
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "手机号不能为空！",
		})
		return
	}
	isPhone, _ := u.phoneExp.MatchString(req.Phone)
	if !isPhone {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "手机号码格式不正确！",
		})
		return
	}
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "发送频繁，请稍后重试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
	}

}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	//Bind 方法回根据 Content-type 来解析你的数据到req 里面
	//解析错了会直接返回400 的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "Email格式不正确")
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次密码输入不一样！")
		return
	}
	isPassword, err := u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误1")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字，特殊字符")
		return
	}

	//调用service 层了
	err = u.svc.SignUp(ctx, domain.User{Email: req.Email, Password: req.Password})
	if err == service.ErrUseDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常2")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名/邮箱或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	//步骤2  使用JWT 设置登录状态
	//生成一个 JWT token

	err = u.setJWTToken(ctx, user.Id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	fmt.Println(user)
	ctx.String(http.StatusOK, "登陆成功")
	return

}

func (u *UserHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       uid,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("iF9BZyZtFYktKQtS9bsJAByiT1aVyt06"))
	if err != nil {
		return err
	}

	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名/邮箱或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	//登陆成功之后获取sessions
	//登陆之后保存登陆信息 步骤2
	sess := sessions.Default(ctx)
	//设置seesion的值
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		MaxAge: 20,
	})
	sess.Save()
	ctx.String(http.StatusOK, "登陆成功2")
	return

}

func (u *UserHandler) LogOut(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	//设置seesion的值
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	ctx.String(http.StatusOK, "登出成功")
	return
}

// Edit 作业 根据用户id修改数据
func (u *UserHandler) Edit(ctx *gin.Context) {
	//定义接收数据
	type UserDataReq struct {
		UserId       int64  `json:"userId"`
		NickName     string `json:"nickName"`
		Birthday     string `json:"birthday"`
		Introduction string `json:"introduction"`
	}
	//实例化一个req
	var req UserDataReq
	//用bing 获取数据
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求数据绑定失败"})
	}
	if req.UserId == 0 {
		ctx.String(http.StatusOK, "用户Id丢失，无法编辑")
		return
	}
	//查看该条用户是否存在
	UserData, _ := u.svc.FindById(ctx, req.UserId)
	if UserData.Id == 0 {
		ctx.String(http.StatusOK, "没有改用户，无法编辑")
		return
	}
	//参数验证
	if req.NickName != "" {
		if len(req.NickName) < 6 || len(req.NickName) > 30 {
			ctx.String(http.StatusOK, "昵称大小需要保持2~10个汉字")
			return
		}
	}
	if req.Birthday != "" {
		isBirthday, _ := u.birthdayExp.MatchString(req.Birthday)
		if !isBirthday {
			ctx.String(http.StatusOK, "生日日期格式不正确，应为YYYY-MM-DD格式")
			return
		}
		if req.Introduction == "" {
			ctx.String(http.StatusOK, "个人简介不能为空")
			return
		}
	}
	//调用service
	err := u.svc.Edit(ctx, domain.User{
		Id:           req.UserId,
		NickName:     req.NickName,
		Birthday:     req.Birthday,
		Introduction: req.Introduction,
	})
	if err != nil {
		ctx.String(http.StatusOK, "修改失败")
		return
	}
	//ctx.JSON(http.StatusOK, UserData)
	ctx.String(http.StatusOK, "修改成功")
}

// Profile 作业回显
func (u *UserHandler) Profile(ctx *gin.Context) {
	//定义接收数据
	type UserAllDataReq struct {
		UserId int64 `json:"userId"`
	}
	//实例化一个req
	var reqSelect UserAllDataReq
	//用bing 获取数据
	if err := ctx.Bind(&reqSelect); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求数据绑定失败"})
	}
	if reqSelect.UserId == 0 {
		ctx.String(http.StatusOK, "用户Id丢失，无法编辑")
		return
	}
	//查看该条用户是否存在
	UserData, _ := u.svc.FindById(ctx, reqSelect.UserId)
	if UserData.Id == 0 {
		ctx.String(http.StatusOK, "未找到该用户")
		return
	}
	ctx.JSON(http.StatusOK, UserData)
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, ok := ctx.Get("claims")
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	//断言
	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	println(claims.Uid)
}

type UserClaims struct {
	jwt.RegisteredClaims
	//声明自己要放进token里的数据
	Uid       int64
	UserAgent string
}
