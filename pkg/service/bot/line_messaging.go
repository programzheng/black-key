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

	if len(parseDate) == 0 {
		return generateErrorTextMessage(), nil
	}

	tt, err := getTimeByTimeString(parseDate[1])
	if err != nil {
		return generateErrorTextMessage(), err
	}

	//every day
	if parseDate[0] == "每天" ||
		parseDate[0] == "每日" ||
		parseDate[0] == "every" ||
		parseDate[0] == "every day" ||
		parseDate[0] == "every-day" {
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
		pdtl := *tt
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
	}

	//specify weekday
	wdtcs := strings.Split(parseDate[0], ",")
	wdens := []string{}
	for _, wdtc := range wdtcs {
		wden := helper.GetWeekDayByTraditionalChinese(wdtc)
		if wden == "" {
			break
		}
		wdens = append(wdens, wden)
	}
	if len(wdens) > 0 {
		templateJSONByte, err := linebot.NewTextMessage(replyText).MarshalJSON()
		if err != nil {
			return generateErrorTextMessage(), err
		}
		pdtl := *tt
		weekDays := strings.Join(wdens, ",")
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
		rpmg := fmt.Sprintf(
			"設置完成將於%s%s\n傳送訊息:%s",
			parseDate[0],
			parseDate[1],
			replyText,
		)
		return linebot.NewTextMessage(rpmg), nil

	}

	//specify date time
	if len(parseDate) == 1 {
		return linebot.NewTextMessage(
			fmt.Sprintf("需設置指定時間，例如: %s 2022-01-01 23:59:59", parseDate[0]),
		), nil
	}
	dts := fmt.Sprintf("%s %s", parseDate[0], parseDate[1])
	dtt, err := time.ParseInLocation("2006-01-02 15:04:05", dts, time.Now().Local().Location())
	if err != nil {
		return generateErrorTextMessage(), err
	}
	ccspm := checkCanSettingPushMessage(dtt)
	if !ccspm {
		return linebot.NewTextMessage(
			"請設置未來的時間",
		), nil
	}

	pdtl := dtt
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

func getTimeByTimeString(ts string) (*time.Time, error) {
	dt := fmt.Sprintf("%s %s", helper.GetNowDateTimeByFormat("2006-01-02"), ts)
	pdtl, err := time.ParseInLocation("2006-01-02 15:04:05", dt, time.Now().Local().Location())
	if err != nil {
		return nil, err
	}
	return &pdtl, nil
}

func checkCanSettingPushMessage(t time.Time) bool {
	return time.Now().Before(t)
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
