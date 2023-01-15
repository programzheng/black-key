package router

import (
	"github.com/programzheng/black-key/internal/controller/file"

	"github.com/gin-gonic/gin"
)

func setRoute(router *gin.Engine) {

	// router.LoadHTMLGlob("dist/view/*")

	router.Static("static", "./storage/upload")
	router.GET("files/:hash_id", file.Get)
}
