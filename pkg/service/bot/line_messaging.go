package bot

import (
	"fmt"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/pkg/helper"
	"github.com/programzheng/black-key/pkg/model/bot"
)

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
		pushDateTime := fmt.Sprintf("%s|%s", weekDays, parseDate[1])
		templateJSON := string(templateJSONByte)
		ln := &bot.LineNotification{
			Service:      "Messaging API",
			PushDateTime: pushDateTime,
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
		pdt, err := helper.ConvertStringToDateTimeString(date)
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
			PushDateTime: pdt,
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
