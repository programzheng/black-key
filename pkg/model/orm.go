package model

import (
	"log"
	"os"
	"time"

	"github.com/programzheng/black-key/config"
	_ "github.com/programzheng/black-key/config"
	"github.com/programzheng/black-key/pkg/helper"

	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func init() {
	var err error
	//?parseTime=true for the database table column type is TIMESTAMP
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?loc=Local&parseTime=true",
		config.Cfg.GetString("DB_USERNAME"),
		config.Cfg.GetString("DB_PASSWORD"),
		config.Cfg.GetString("DB_HOST"),
		config.Cfg.GetString("DB_PORT"),
		config.Cfg.GetString("DB_DATABASE"))
	fmt.Printf("connect: %v database\n", dsn)
	gormConfig := &gorm.Config{}
	if helper.ConvertToBool(config.Cfg.GetString("DB_DEBUG")) {
		gormConfig.Logger = getLogger()
	}
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)

	if err != nil {
		log.Println("DataBase error:", err)
	}
}

func getLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
}

func Migrate(models ...interface{}) {
	DB.AutoMigrate(models...)
}
