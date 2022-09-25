package line

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/pkg/library/line/bot/template"
	model "github.com/programzheng/black-key/pkg/model/bot"
	"github.com/programzheng/black-key/pkg/service/bot"
	"golang.org/x/exp/slices"
)

type Todo struct {
	BotClient *linebot.Client
	ToID      string
	Template  *linebot.TextMessage
}

func RunSchedule() {
	//get line notifications from database
	ln := &model.LineNotification{}
	lns, err := ln.Get(nil, nil)
	if err != nil {
		log.Printf("pkg/job/line/todo RunSchedule Get error: %v", err)
	}
	for _, ln := range lns {
		canPush := checkCanPushLineNotification(ln)
		if !canPush {
			continue
		}
		switch ln.Type {
		case string(linebot.MessageTypeText):
			var tp linebot.TextMessage
			data := []byte(ln.Template)
			err := json.Unmarshal(data, &tp)
			if err != nil {
				log.Printf("pkg/job/line/todo RunSchedule json.Unmarshal error: %v", err)
				return
			}
			pushID := getPushID(ln)
			if pushID != "" {
				err := bot.LinePushMessage(pushID, template.Text(tp.Text))
				if err != nil {
					log.Printf("pkg/job/line/todo RunSchedule LinePushMessage error: %v", err)
				}
				err = afterPushLineNotification(ln)
				if err != nil {
					log.Printf("pkg/job/line/todo RunSchedule PermanentlyDelete error: %v", err)
				}
			}
		}
	}
}

func convertDateTimeToOnlyDateTimeString(dateTime string) string {
	s := strings.Split(dateTime, "|")
	if len(s) > 1 {
		return s[1]
	}
	return dateTime
}

func convertTimeToPushDateTime(dateTime string) string {
	s := strings.Split(dateTime, "|")
	if len(s) > 1 {
		//is by weekday push
		weekDays := strings.Split(s[0], ",")
		nt := time.Now()
		nowWeekDay := nt.Weekday().String()
		for _, weekDay := range weekDays {
			if nowWeekDay == weekDay {
				return fmt.Sprintf("%s %s", nt.Format("2006-01-02"), s[1])
			}
		}
	}
	return dateTime
}

func checkCanPushLineNotification(ln *model.LineNotification) bool {
	pushDateTime := ln.PushDateTime

	if ln.PushCycle != "specify" {
		nowWeekDay := time.Now().Weekday()
		pcs := strings.Split(ln.PushCycle, ",")
		if slices.Contains(pcs, nowWeekDay.String()) {
			nowDate := time.Now().Format("2006-01-02")
			st := pushDateTime.Format("15:04:05")
			pdts := fmt.Sprintf("%s %s", nowDate, st)
			pushDateTime, err := time.ParseInLocation("2006-01-02 15:04:05", pdts, time.Now().Local().Location())
			if err != nil {
				return false
			}
			minTolerantDateTime := time.Now().Add(-30 * time.Second)
			maxTolerantDateTime := time.Now().Add(30 * time.Second)
			return minTolerantDateTime.Before(pushDateTime) && maxTolerantDateTime.After(pushDateTime)
		}
	}
	minTolerantDateTime := time.Now().Add(-30 * time.Second)
	maxTolerantDateTime := time.Now().Add(30 * time.Second)
	return minTolerantDateTime.Before(pushDateTime) && maxTolerantDateTime.After(pushDateTime)
}

func afterPushLineNotification(ln *model.LineNotification) error {
	var err error

	//Limit < 0 is unlimited
	if ln.Limit > 0 {
		ln.Limit -= 1
		err := ln.Save()
		if err != nil {
			return err
		}
	}

	// if ln.PushCycle != "specify" && ln.Limit == -1 {
	// 	//push weekday is next weekday
	// 	nextDateTime := ln.PushDateTime.AddDate(0, 0, 1)
	// 	wds := nextDateTime.Weekday().String()
	// 	weekDayCycle := strings.Split(ln.PushCycle, ",")
	// 	if slices.Contains(weekDayCycle, wds) {
	// 		ln.PushDateTime = nextDateTime
	// 		err := ln.Save()
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	if ln.Limit == 0 {
		err = ln.Delete()
		if err != nil {
			return err
		}
	}

	return nil
}

func getPushID(ln *model.LineNotification) string {
	if ln.RoomID != "" {
		return ln.RoomID
	}
	if ln.GroupID != "" {
		return ln.GroupID
	}
	if ln.UserID != "" {
		return ln.UserID
	}
	return ""
}
