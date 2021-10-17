package router

import (
	"github.com/gin-gonic/gin"

	"black-key/pkg/controller/admin"
	"black-key/pkg/controller/auth"
	"black-key/pkg/controller/file"
	"black-key/pkg/middleware"
)

func setAPIRoute(router *gin.Engine) {
	apiGroup := router.Group("/api/v1")
	adminGroup := apiGroup.Group("/admins")
	{
		adminGroup.POST("", admin.Register)
		adminGroup.POST("login", admin.Login)
		adminGroup.POST("auth", auth.VaildAdmin)
	}
	apiGroup.Use(middleware.ValidJSONWebToken())
	{
		adminsGroup := apiGroup.Group("/admins")
		{
			adminsGroup.GET("", admin.Get)
		}
		filesGroup := apiGroup.Group("/files")
		{
			filesGroup.POST("", file.Upload)
		}
	}

}
