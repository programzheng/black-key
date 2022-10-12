package bot

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/programzheng/black-key/internal/model/bot"
	"golang.org/x/exp/slices"
)

func RefreshTodoByAfterPushDateTime() error {
	now := time.Now().Local()

	lns, err := (&bot.LineNotification{}).GetAfterPushDateTime(now)
	if err != nil {
		return err
	}
	for _, ln := range lns {
		//check can set new push date time
		weekdays := strings.Split(ln.PushCycle, ",")
		if !slices.Contains(weekdays, now.Weekday().String()) {
			return errors.New("now is not push cycle")
		}
		//set new push date time
		pushDateTime := ln.PushDateTime
		oldPushTime := pushDateTime.Format("15:04:05")
		nowDate := now.Format("2006-01-02")
		newPushDateTimeString := fmt.Sprintf("%s %s", nowDate, oldPushTime)
		newPushDateTime, err := time.ParseInLocation("2006-01-02 15:04:05", newPushDateTimeString, now.Location())
		if err != nil {
			return err
		}
		ln.PushDateTime = newPushDateTime
		err = ln.Save()
		if err != nil {
			return err
		}
	}
	return nil
}
