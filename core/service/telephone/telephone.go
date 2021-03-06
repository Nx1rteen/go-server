// Copyright 2019 Axetroy. All rights reserved. MIT license.
package telephone

import (
	"fmt"
	"github.com/axetroy/go-server/core/config"
	"log"
)

type provider string

var (
	client          *Telephone             // 发送短信的客户端
	providerAliyun  provider   = "aliyun"  // 阿里云
	providerTencent provider   = "tencent" // 腾讯云
)

// 邮箱提供这应提供的对象
type Telephone interface {
	getAuthTemplateID() string                                                 // 身份验证的模版 ID
	getResetPasswordTemplateID() string                                        // 重置密码的模版 ID
	getRegisterTemplateID() string                                             // 注册帐号的模版 ID
	send(phone string, templateID string, templateMap map[string]string) error // 发送验证码
	SendRegisterCode(phone string, code string) error                          // 发送注册验证码
	SendAuthCode(phone string, code string) error                              // 发送身份验证码
	SendResetPasswordCode(phone string, code string) error                     // 发送重置密码验证码
}

func init() {
	switch provider(config.Telephone.Provider) {
	case providerAliyun:
		initClient(NewAliyun())
		break
	case providerTencent:
		initClient(NewTencent())
		break
	default:
		log.Fatal(fmt.Sprintf(`Invalid telephone provider "%s"`, config.Telephone.Provider))
	}
}

func initClient(s Telephone) {
	client = &s
}

func GetClient() Telephone {
	return *client
}
