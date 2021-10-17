package router

import (
	"black-key/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func setMiddleware(router *gin.Engine) {
	router.Use(middleware.CORSMiddleware())
}
