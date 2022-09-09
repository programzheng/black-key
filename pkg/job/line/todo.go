package line

import (
	"encoding/json"
	"log"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/pkg/library/line/bot/template"
	model "github.com/programzheng/black-key/pkg/model/bot"
	"github.com/programzheng/black-key/pkg/service/bot"
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
		lnDateTime, err := time.ParseInLocation("2006-01-02 15:04:05", ln.PushDateTime, time.Now().Local().Location())
		if err != nil {
			log.Printf("pkg/job/line/todo RunSchedule time.Parse error: %v", err)
		}
		nowDateTime := time.Now()
		if nowDateTime.Before(lnDateTime) {
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
			if ln.UserID != "" {
				err := bot.LinePushMessage(ln.UserID, template.Text(tp.Text))
				if err != nil {
					log.Printf("pkg/job/line/todo RunSchedule LinePushMessage error: %v", err)
				}
				err = ln.Delete()
				if err != nil {
					log.Printf("pkg/job/line/todo RunSchedule PermanentlyDelete error: %v", err)
				}
			}
		}
	}
}
