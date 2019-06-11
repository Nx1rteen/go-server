// Copyright 2019 Axetroy. All rights reserved. MIT license.
package banner_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/module/banner"
	"github.com/axetroy/go-server/module/banner/banner_model"
	"github.com/axetroy/go-server/module/banner/banner_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetList(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	ctx := schema.Context{
		Uid: adminInfo.Id,
	}

	{
		var (
			image    = "test"
			href     = "test"
			platform = banner_model.BannerPlatformApp
		)

		r := banner.Create(schema.Context{
			Uid: adminInfo.Id,
		}, banner.CreateParams{
			Image:    image,
			Href:     href,
			Platform: platform,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := banner_schema.Banner{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		defer banner.DeleteBannerById(n.Id)
	}

	// 获取列表
	{
		r := banner.GetList(ctx, banner.Query{})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		banners := make([]banner_schema.Banner, 0)

		assert.Nil(t, tester.Decode(r.Data, &banners))

		assert.Equal(t, schema.DefaultLimit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
		assert.IsType(t, 1, r.Meta.Num)
		assert.IsType(t, int64(1), r.Meta.Total)

		assert.True(t, len(banners) >= 1)

		for _, b := range banners {
			assert.IsType(t, "string", b.Image)
			assert.IsType(t, "string", b.Href)
			assert.IsType(t, banner_model.BannerPlatformApp, b.Platform)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}

func TestGetListRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	{
		var (
			image    = "test"
			href     = "test"
			platform = banner_model.BannerPlatformApp
		)

		r := banner.Create(schema.Context{
			Uid: adminInfo.Id,
		}, banner.CreateParams{
			Image:    image,
			Href:     href,
			Platform: platform,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := banner_schema.Banner{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		defer banner.DeleteBannerById(n.Id)
	}

	{
		r := tester.HttpAdmin.Get("/v1/banner", nil, &header)

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		banners := make([]banner_schema.Banner, 0)

		assert.Nil(t, tester.Decode(res.Data, &banners))

		for _, b := range banners {
			assert.IsType(t, "string", b.Image)
			assert.IsType(t, "string", b.Href)
			assert.IsType(t, banner_model.BannerPlatformApp, b.Platform)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}