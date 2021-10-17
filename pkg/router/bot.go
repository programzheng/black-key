package router

import (
	"github.com/gin-gonic/gin"

	"black-key/pkg/controller/bot"
)

func setBotRouter(router *gin.Engine) {
	botGroup := router.Group("/bot")
	lineGroup := botGroup.Group("/line")
	{
		lineGroup.POST("", bot.LineWebHook)
		lineGroup.POST("push", bot.LinePush)
	}
}
