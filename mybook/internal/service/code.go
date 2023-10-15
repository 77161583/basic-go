package service

import (
	"basic-go/mybook/internal/repository"
	"basic-go/mybook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

const codeTplId = "213123"

var (
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
)

type CodeServicePackage interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CodeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeServicePackage {
	return &CodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

// Send 发送验证码 需要什么参数
func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	//biz 区别业务场景
	//生成一个验证码
	code := svc.generateCode()
	//塞入到redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	//发出去
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	return err
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *CodeService) generateCode() string {
	//六位数 num 在 0，999999之间
	num := rand.Intn(1000000)
	//不够六位的，加上前导0    例： 001234
	return fmt.Sprintf("%06d", num)
}

//func (svc *CodeService) VerifyV1(ctx context.Context, biz string) error {
//
//}
