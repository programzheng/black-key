package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetRouter(router *gin.Engine) {
	// 中間件
	setMiddleware(router)
	// 設置Web Route
	setRoute(router)
	fmt.Println(gin.Mode())
	// 設置API Router
	setAPIRoute(router)
	// 設置Bot Router
	setBotRouter(router)
}
