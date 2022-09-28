package router

import (
	"github.com/programzheng/black-key/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setMiddleware(router *gin.Engine) {
	router.Use(middleware.CORSMiddleware())
}
