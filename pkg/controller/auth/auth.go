package auth

import (
	"errors"
	"strings"

	"github.com/programzheng/black-key/pkg/helper"
	"github.com/programzheng/black-key/pkg/service/auth"

	"github.com/gin-gonic/gin"
)

func VaildAdmin(ctx *gin.Context) {
	requestToken := ctx.GetHeader("Authorization")
	splitToken := strings.Split(requestToken, "Bearer")
	if len(splitToken) != 2 {
		//return not vaild
		helper.Unauthorized(ctx, errors.New("驗證失敗"))
		return
	}

	token := strings.TrimSpace(splitToken[1])

	//曾經有登入記錄
	adminLogin, err := (&auth.AdminLogin{
		Token: token,
	}).GetAdminLogin()

	if adminLogin.ID == 0 && err != nil {
		helper.Unauthorized(ctx, errors.New("請重新登入"))
		return
	}

	vaildResult := helper.ValidJSONWebToken(token)
	if !vaildResult {
		helper.Unauthorized(ctx, errors.New("請重新登入"))
		return
	}

	if adminLogin.Remember {

	}

	helper.Success(ctx, nil, nil)
	return
}
