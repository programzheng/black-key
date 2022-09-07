package model

import (
	"github.com/programzheng/black-key/config"
	_ "github.com/programzheng/black-key/config"

	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

func init() {
	var err error
	//?parseTime=true for the database table column type is TIMESTAMP
	setting := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?loc=Local&parseTime=true",
		config.Cfg.GetString("DB_USERNAME"),
		config.Cfg.GetString("DB_PASSWORD"),
		config.Cfg.GetString("DB_HOST"),
		config.Cfg.GetString("DB_PORT"),
		config.Cfg.GetString("DB_DATABASE"))
	fmt.Printf("connect: %v database\n", setting)
	DB, err = gorm.Open(config.Cfg.GetString("DB_CONNECTION"), setting)

	if err != nil {
		log.Println("DataBase error:", err)
	}
}

func Migrate(models ...interface{}) {
	DB.AutoMigrate(models...)
}
