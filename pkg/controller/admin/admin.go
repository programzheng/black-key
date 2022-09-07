package admin

import (
	"errors"

	"github.com/programzheng/black-key/config"
	"github.com/programzheng/black-key/pkg/helper"
	"github.com/programzheng/black-key/pkg/service/admin"
	"github.com/programzheng/black-key/pkg/service/auth"

	"github.com/gin-gonic/gin"
)

var adminService admin.Admin

func Register(ctx *gin.Context) {
	if err := ctx.Bind(&adminService); err != nil {
		helper.BadRequest(ctx, err)
		return
	}

	//hash password
	adminService.Password = helper.CreateHash(adminService.Password)
	result, err := adminService.Add()
	if err != nil {
		helper.Fail(ctx, err)
		return
	}

	helper.Success(ctx, result, nil)
	return
}

func Login(ctx *gin.Context) {
	login := auth.Login{}
	if err := ctx.Bind(&login); err != nil {
		helper.BadRequest(ctx, err)
		return
	}
	admin, err := (&admin.Admin{
		Account: login.Account,
	}).GetForLogin()
	if err != nil {
		helper.Fail(ctx, errors.New("帳號錯誤"))
		return
	}
	err = helper.CheckHash(admin.Password, login.Password)
	if err != nil {
		helper.Fail(ctx, errors.New("密碼錯誤"))
		return
	}
	secret := []byte(config.Cfg.GetString("JWT_SECRET"))
	token := helper.CreateJWT(secret)
	adminLogin := auth.AdminLogin{
		AdminID: admin.ID,
		Token:   token.Token,
		IP:      ctx.ClientIP(),
	}
	if err := adminLogin.AddAdminLogin(); err != nil {
		helper.Fail(ctx, err)
		return
	}

	helper.Success(ctx, token, nil)
	return
}

func Get(ctx *gin.Context) {
	adminService := admin.Admin{}
	if err := ctx.Bind(&adminService); err != nil {
		helper.BadRequest(ctx, err)
		return
	}
	admins, err := adminService.Get()
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	data := make(map[string]interface{})
	data["list"] = admins
	// data["Total"] = total
	helper.Success(ctx, data, nil)
	return
}
