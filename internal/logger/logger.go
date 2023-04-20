package logger

import (
	"fmt"
	"time"

	"github.com/programzheng/black-key/config"
	"github.com/programzheng/black-key/internal/helper"
)

func getLogFileNameByMode() string {
	logFileName := config.Cfg.GetString("APP_NAME")
	if config.Cfg.GetString("LOG_MODE") == "daily" {
		return fmt.Sprintf("%s_%s.log", logFileName, time.Now().Format(helper.Iso8601))
	}
	return logFileName
}
