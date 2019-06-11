// Copyright 2019 Axetroy. All rights reserved. MIT license.
package email

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/module/user/user_error"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/axetroy/go-server/service/email"
	"github.com/axetroy/go-server/service/redis"
	"github.com/axetroy/go-server/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type SendResetPasswordEmailParams struct {
	To string `json:"to"` // 发送给谁
}

func GenerateResetCode(uid string) string {
	// 生成重置码
	var codeId = "reset-" + util.GenerateId() + uid
	return util.MD5(codeId)
}

func SendResetPasswordEmail(input SendResetPasswordEmailParams) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.ErrUnknown
			}
		}

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			res.Data = nil
			res.Message = err.Error()
			res.Data = false
		} else {
			res.Data = true
			res.Status = schema.StatusSuccess
		}
	}()

	userInfo := user_model.User{
		Email: &input.To,
	}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user_error.ErrUserNotExist
		}
		return
	}

	// 生成重置码
	var code = GenerateResetCode(userInfo.Id)

	// set activationCode to redis
	if err = redis.ResetCodeClient.Set(code, userInfo.Id, time.Minute*30).Err(); err != nil {
		return
	}

	e := email.NewMailer()

	// send email
	if err = e.SendForgotPasswordEmail(input.To, code); err != nil {
		// 邮件没发出去的话，删除redis的key
		_ = redis.ResetCodeClient.Del(code).Err()
		return
	}

	return

}

func SendResetPasswordEmailRouter(ctx *gin.Context) {
	var (
		input SendResetPasswordEmailParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = exception.ErrInvalidParams
		return
	}

	res = SendResetPasswordEmail(input)
}