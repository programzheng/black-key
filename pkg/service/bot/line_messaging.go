package bot

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/pkg/helper"
	"github.com/programzheng/black-key/pkg/model/bot"
	log "github.com/sirupsen/logrus"
)

func getTodo(lineId LineID) (interface{}, error) {
	lns, err := (&bot.LineNotification{}).Get(map[string]interface{}{
		"user_id":  lineId.UserID,
		"group_id": lineId.GroupID,
		"room_id":  lineId.RoomID,
	}, nil)
	if err != nil {
		return nil, err
	}
	if len(lns) == 0 {
		return linebot.NewTextMessage("沒有資料"), nil
	}
	carouselColumns := []*linebot.CarouselColumn{}
	for _, ln := range lns {
		var tp linebot.TextMessage
		data := []byte(ln.Template)
		err := json.Unmarshal(data, &tp)
		if err != nil {
			log.Printf("pkg/service/bot/line_messaging getTodo json.Unmarshal error: %v", err)
			return nil, err
		}
		pushDateTime := ln.PushDateTime.String()
		deletePostBackAction := LinePostBackAction{
			Action: "delete line notification",
			Data: map[string]interface{}{
				"ID": ln.ID,
			},
		}
		deletePostBackActionJson, err := json.Marshal(deletePostBackAction)
		if err != nil {
			log.Printf("pkg/service/bot/line_messaging getTodo deletePostBackActionJson json.Marshal error: %v", err)
			return nil, err
		}
		carouselColumn := linebot.NewCarouselColumn(
			"",
			tp.Text,
			fmt.Sprintf("發送時間: %s", pushDateTime),
			linebot.NewPostbackAction(
				"刪除",
				string(deletePostBackActionJson),
				"",
				"",
			),
		)
		carouselColumns = append(carouselColumns, carouselColumn)
	}
	carouselTemplate := linebot.NewCarouselTemplate(carouselColumns...)

	return linebot.NewTemplateMessage("所有提醒", carouselTemplate), nil
}

func convertPushDateTime(pdt string) string {
	s := strings.Split(pdt, "|")
	if len(s) == 1 {
		return pdt
	}
	period := s[0]
	dateTime := s[1]
	switch period {
	case "Sunday,Monday,Tuesday,Wednesday,Thursday,Friday,Saturday":
		return fmt.Sprintf("每天 %s", dateTime)
	}
	return ""
}

func todo(lineId LineID, text string) (interface{}, error) {
	parseText := strings.Split(text, "|")
	date := parseText[1]
	replyText := parseText[2]
	parseDate := strings.Split(date, " ")
	switch parseDate[0] {
	// TODO|every 19:55|測試29號13:30送出
	case "每天", "每日", "every", "every day", "every-day":
		if len(parseDate) == 1 {
			return linebot.NewTextMessage(
				fmt.Sprintf("需設置指定時間，例如: %s 23:59:59", parseDate[0]),
			), nil
		}

		templateJSONByte, err := linebot.NewTextMessage(replyText).MarshalJSON()
		if err != nil {
			return generateErrorTextMessage(), err
		}
		weekDays := strings.Join(helper.GetWeekDays(), ",")
		dt := fmt.Sprintf("%s %s", helper.GetNowDateTimeByFormat("2006-01-02"), parseDate[1])
		pdtl, err := time.ParseInLocation("2006-01-02 15:04:05", dt, time.Now().Local().Location())
		if err != nil {
			return linebot.NewTextMessage(
				"時間格式錯誤，請重新輸入，例如:每天 23:59:59",
			), nil
		}
		templateJSON := string(templateJSONByte)
		ln := &bot.LineNotification{
			Service:      "Messaging API",
			PushCycle:    weekDays,
			PushDateTime: pdtl,
			Limit:        -1,
			UserID:       lineId.UserID,
			GroupID:      lineId.GroupID,
			RoomID:       lineId.RoomID,
			Type:         string(linebot.MessageTypeText),
			Template:     templateJSON,
		}
		_, err = ln.Add()
		if err != nil {
			return generateErrorTextMessage(), err
		}
		return linebot.NewTextMessage("設置完成將於每天" + parseDate[1] + "\n傳送訊息:" + replyText), nil
	// TODO|2020/02/29 13:00|測試29號13:30送出
	default:
		if len(parseDate) == 1 {
			return linebot.NewTextMessage(
				fmt.Sprintf("需設置指定時間，例如: %s 2022-01-01 23:59:59", parseDate[0]),
			), nil
		}

		pdt, err := helper.ConvertStringToDateTimeString(date)
		pdtl, err := time.ParseInLocation("2006-01-02 15:04:05", pdt, time.Now().Local().Location())
		if err != nil {
			return generateErrorTextMessage(), err
		}
		templateJSONByte, err := linebot.NewTextMessage(replyText).MarshalJSON()
		if err != nil {
			return generateErrorTextMessage(), err
		}
		templateJSON := string(templateJSONByte)
		ln := &bot.LineNotification{
			Service:      "Messaging API",
			PushCycle:    "specify",
			PushDateTime: pdtl,
			Limit:        1,
			UserID:       lineId.UserID,
			GroupID:      lineId.GroupID,
			RoomID:       lineId.RoomID,
			Type:         string(linebot.MessageTypeText),
			Template:     templateJSON,
		}
		_, err = ln.Add()
		if err != nil {
			return generateErrorTextMessage(), err
		}

		return linebot.NewTextMessage("設置完成將於" + date + "\n傳送訊息:" + replyText), nil
	}
}

func deleteTodoByPostBack(lpba *LinePostBackAction) interface{} {
	id := uint(lpba.Data["ID"].(float64))
	ln, err := bot.LineNotificationFirstByID(id)
	if err != nil {
		return nil
	}
	err = ln.Delete()
	if err != nil {
		return linebot.NewTextMessage(
			"刪除失敗",
		)
	}

	return linebot.NewTextMessage("刪除成功")
}
