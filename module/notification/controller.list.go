// Copyright 2019 Axetroy. All rights reserved. MIT license.
package notification

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/notification/notification_model"
	"github.com/axetroy/go-server/module/notification/notification_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

// Query params
type Query struct {
	schema.Query
}

// GetList get notification list
func GetListByUser(context schema.Context, input Query) (res schema.List) {
	var (
		err  error
		data = make([]notification_schema.Notification, 0)
		meta = &schema.Meta{}
		tx   *gorm.DB
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
			res.Message = err.Error()
			res.Data = nil
			res.Meta = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
			res.Meta = meta
		}
	}()

	query := input.Query

	query.Normalize()

	tx = database.Db.Begin()

	var total int64

	list := make([]notification_model.Notification, 0)

	if err = tx.Table(new(notification_model.Notification).TableName()).Limit(query.Limit).Offset(query.Limit * query.Page).Find(&list).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := notification_schema.Notification{}
		if er := mapstructure.Decode(v, &d.NotificationPure); er != nil {
			err = er
			return
		}
		d.CreatedAt = v.CreatedAt.Format(time.RFC3339Nano)
		d.UpdatedAt = v.UpdatedAt.Format(time.RFC3339Nano)

		// 查询用户是否已读通知
		mark := notification_model.NotificationMark{
			Id:  v.Id,
			Uid: context.Uid,
		}

		if err = tx.Last(&mark).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				d.Read = false
				d.ReadAt = ""
				err = nil
			} else {
				break
			}
		} else {
			d.Read = true
			d.ReadAt = mark.CreatedAt.Format(time.RFC3339Nano)
		}

		data = append(data, d)
	}

	meta.Total = total
	meta.Num = len(data)
	meta.Page = query.Page
	meta.Limit = query.Limit

	return
}

// GetList get notification list
func GetListAdmin(context schema.Context, input Query) (res schema.List) {
	var (
		err  error
		data = make([]notification_schema.NotificationAdmin, 0)
		meta = &schema.Meta{}
		tx   *gorm.DB
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
			res.Message = err.Error()
			res.Data = nil
			res.Meta = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
			res.Meta = meta
		}
	}()

	query := input.Query

	query.Normalize()

	tx = database.Db.Begin()

	var total int64

	list := make([]notification_model.Notification, 0)

	if err = tx.Table(new(notification_model.Notification).TableName()).Limit(query.Limit).Offset(query.Limit * query.Page).Find(&list).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := notification_schema.NotificationAdmin{}
		if er := mapstructure.Decode(v, &d.NotificationPureAdmin); er != nil {
			err = er
			return
		}
		d.CreatedAt = v.CreatedAt.Format(time.RFC3339Nano)
		d.UpdatedAt = v.UpdatedAt.Format(time.RFC3339Nano)
		data = append(data, d)
	}

	meta.Total = total
	meta.Num = len(data)
	meta.Page = query.Page
	meta.Limit = query.Limit

	return
}

// GetListRouter get list router
func GetListUserRouter(ctx *gin.Context) {
	var (
		err   error
		res   = schema.List{}
		input Query
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindQuery(&input); err != nil {
		err = exception.ErrInvalidParams
		return
	}

	res = GetListByUser(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}

// GetListRouter get list router
func GetListAdminRouter(ctx *gin.Context) {
	var (
		err   error
		res   = schema.List{}
		input Query
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindQuery(&input); err != nil {
		err = exception.ErrInvalidParams
		return
	}

	res = GetListAdmin(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}