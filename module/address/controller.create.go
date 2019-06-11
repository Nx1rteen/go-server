// Copyright 2019 Axetroy. All rights reserved. MIT license.
package address

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/address/address_model"
	"github.com/axetroy/go-server/module/address/address_schema"
	"github.com/axetroy/go-server/module/user/user_error"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type CreateAddressParams struct {
	Name         string `json:"name" valid:"required~请填写收货人"`         // 收货人
	Phone        string `json:"phone" valid:"required~请输入收货人电话号码"`    // 收货人手机号
	ProvinceCode string `json:"province_code" valid:"required~请选择省份"` // 省份代码
	CityCode     string `json:"city_code" valid:"required~请选择城市"`     // 城市代码
	AreaCode     string `json:"area_code" valid:"required~请选择区域"`     // 区域代码
	Address      string `json:"address" valid:"required~请输入详细地址"`     // 详细的地址
	IsDefault    *bool  `json:"is_default"`                           // 是否是默认地址
}

func Create(context schema.Context, input CreateAddressParams) (res schema.Response) {
	var (
		err          error
		data         address_schema.Address
		tx           *gorm.DB
		isDefault    = false
		isValidInput bool
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
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = exception.ErrInvalidParams
		return
	}

	// 校验 省份代码
	if _, ok := ProvinceCode[input.ProvinceCode]; !ok {
		err = ErrAddressInvalidProvinceCode
		return
	}

	// 校验 城市代码
	if _, ok := CityCode[input.CityCode]; !ok {
		err = ErrAddressInvalidCityCode
		return
	}

	// 校验 区域代码
	if _, ok := CountryCode[input.AreaCode]; !ok {
		err = ErrAddressInvalidAreaCode
		return
	}

	tx = database.Db.Begin()

	userInfo := user_model.User{
		Id: context.Uid,
	}

	if err = tx.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user_error.ErrUserNotExist
		}
		return
	}

	if input.IsDefault != nil {
		isDefault = *input.IsDefault

		defaultAddress := address_model.Address{
			Uid:       context.Uid,
			IsDefault: true,
		}
		if err = tx.First(&defaultAddress).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = nil
			} else {
				return
			}
		} else {
			// 如果存在了默认地址，则取消它的默认属性
			if err = tx.Model(&defaultAddress).UpdateColumn(address_model.Address{
				IsDefault: false,
			}).Error; err != nil {
				return
			}
		}

	} else {
		firstAddress := address_model.Address{
			Uid: context.Uid,
		}
		if err = tx.Where(&firstAddress).First(&firstAddress).Error; err != nil {
			// 如果还没有设置过地址，那么这次设置就是默认地址
			if err == gorm.ErrRecordNotFound {
				err = nil
				isDefault = true
			} else {
				return
			}
		}
	}

	AddressInfo := address_model.Address{
		Uid:          context.Uid,
		Name:         input.Name,
		Phone:        input.Phone,
		ProvinceCode: input.ProvinceCode,
		CityCode:     input.CityCode,
		AreaCode:     input.AreaCode,
		Address:      input.Address,
		IsDefault:    isDefault,
	}

	if err = tx.Create(&AddressInfo).Error; err != nil {
		return
	}

	if er := mapstructure.Decode(AddressInfo, &data.AddressPure); er != nil {
		err = er
		return
	}

	data.CreatedAt = AddressInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = AddressInfo.UpdatedAt.Format(time.RFC3339Nano)
	return
}

func CreateRouter(ctx *gin.Context) {
	var (
		input CreateAddressParams
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

	res = Create(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}