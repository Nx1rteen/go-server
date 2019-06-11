// Copyright 2019 Axetroy. All rights reserved. MIT license.
package message_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/message"
	"github.com/axetroy/go-server/module/message/message_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetMessage(t *testing.T) {
	var (
		messageId string
	)

	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	// 2. 先创建一篇消息作为测试
	{
		var (
			title   = "test"
			content = "test"
		)

		r := message.Create(schema.Context{
			Uid: adminInfo.Id,
		}, message.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := message_schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		messageId = n.Id

		defer message.DeleteMessageById(n.Id)
	}

	// 3. 获取文章公告
	{
		r := message.Get(schema.Context{
			Uid: userInfo.Id,
		}, messageId)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		messageInfo := r.Data.(message_schema.Message)

		assert.Equal(t, "test", messageInfo.Title)
		assert.Equal(t, "test", messageInfo.Content)
	}

}

func TestGetAdmin(t *testing.T) {
	var (
		messageId string
	)

	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	// 2. 先创建一篇消息作为测试
	{
		var (
			title   = "test"
			content = "test"
		)

		r := message.Create(schema.Context{
			Uid: adminInfo.Id,
		}, message.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := message_schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		messageId = n.Id

		defer message.DeleteMessageById(n.Id)
	}

	// 3. 获取文章公告
	{
		r := message.GetByAdmin(schema.Context{
			Uid: adminInfo.Id,
		}, messageId)

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		messageInfo := r.Data.(message_schema.MessageAdmin)

		assert.Equal(t, "test", messageInfo.Title)
		assert.Equal(t, "test", messageInfo.Content)
	}

}

func TestGetRouter(t *testing.T) {
	var (
		messageId string
	)

	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	// 2. 先创建一篇消息作为测试
	{
		var (
			title   = "test"
			content = "test"
		)

		r := message.Create(schema.Context{
			Uid: adminInfo.Id,
		}, message.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := message_schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		messageId = n.Id

		defer message.DeleteMessageById(n.Id)
	}

	// 用户接口获取
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + userInfo.Token,
		}

		r := tester.HttpUser.Get("/v1/message/m/"+messageId, nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		n := message_schema.Message{}

		assert.Nil(t, tester.Decode(res.Data, &n))

		assert.Equal(t, "test", n.Title)
		assert.Equal(t, "test", n.Content)
		assert.IsType(t, "string", n.CreatedAt)
		assert.IsType(t, "string", n.UpdatedAt)
	}

}

func TestGetAdminRouter(t *testing.T) {
	var (
		messageId string
	)

	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	// 2. 先创建一篇消息作为测试
	{
		var (
			title   = "test"
			content = "test"
		)

		r := message.Create(schema.Context{
			Uid: adminInfo.Id,
		}, message.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := message_schema.Message{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		messageId = n.Id

		defer message.DeleteMessageById(n.Id)
	}

	// 管理员接口获取
	{
		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		r := tester.HttpAdmin.Get("/v1/message/m/"+messageId, nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		n := message_schema.MessageAdmin{}

		assert.Nil(t, tester.Decode(res.Data, &n))

		assert.Equal(t, "test", n.Title)
		assert.Equal(t, "test", n.Content)
		assert.IsType(t, "string", n.CreatedAt)
		assert.IsType(t, "string", n.UpdatedAt)
	}
}