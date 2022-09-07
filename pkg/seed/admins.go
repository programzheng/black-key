package seed

import (
	"github.com/programzheng/black-key/pkg/helper"
	"github.com/programzheng/black-key/pkg/model/admin"

	"github.com/jinzhu/gorm"
)

func CreateAdmin(db *gorm.DB, account string, password string) error {
	password = helper.CreateHash(password)
	return db.Create(&admin.Admin{Account: account, Password: password}).Error
}
