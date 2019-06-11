// Copyright 2019 Axetroy. All rights reserved. MIT license.
package resource

import (
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/module/uploader"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func Image(ctx *gin.Context) {
	filename := ctx.Param("filename")
	originImagePath := path.Join(uploader.Config.Path, uploader.Config.Image.Path, filename)
	if fs.PathExists(originImagePath) == false {
		// if the path not found
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	http.ServeFile(ctx.Writer, ctx.Request, originImagePath)
}