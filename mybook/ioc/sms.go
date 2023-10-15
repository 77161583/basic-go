package ioc

import (
	"basic-go/mybook/internal/service/sms"
	"basic-go/mybook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	//这里可以换内存，或者换其他
	return memory.NewService()
}
