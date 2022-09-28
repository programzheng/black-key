package router

import (
	"github.com/programzheng/black-key/internal/controller/file"
	"github.com/programzheng/black-key/internal/controller/job"

	"github.com/gin-gonic/gin"
)

func setRoute(router *gin.Engine) {
	router.GET("/jobrunner/json", job.JobJson)

	// router.LoadHTMLGlob("dist/view/*")

	// router.GET("/jobrunner/html", job.JobHtml)

	router.GET("files/:hash_id", file.Get)

}
