package ioc

import (
	"github.com/LXD-c/basic-go/webook/internal/service/sms"
	"github.com/LXD-c/basic-go/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	// 具体实现随便你换
	return memory.NewService()
}
