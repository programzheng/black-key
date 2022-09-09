package bot

import (
	"fmt"
	"strings"

	"github.com/programzheng/black-key/pkg/helper"
	"github.com/programzheng/black-key/pkg/model/bot"

	"github.com/line/line-bot-sdk-go/linebot"
)

func UserParseTextGenTemplate(lineId LineID, text string) (interface{}, error) {
	parseText := strings.Split(text, "|")

	if len(parseText) == 1 {

	}
	switch parseText[0] {
	// Line相關資訊
	case "資訊":
		return linebot.NewTextMessage(fmt.Sprintf("RoomID:%v\nGroupID:%v\nUserID:%v", lineId.RoomID, lineId.GroupID, lineId.UserID)), nil
	case "我的大頭貼":
		lineMember, err := botClient.GetGroupMemberProfile(lineId.GroupID, lineId.UserID).Do()
		if err != nil {
			return generateErrorTextMessage(), err
		}
		return linebot.NewImageMessage(lineMember.PictureURL, lineMember.PictureURL), nil
	case "提醒", "通知", "TODO":
		date := parseText[1]
		replyText := parseText[2]
		parseDate := strings.Split(date, " ")
		switch parseDate[0] {
		// TODO|every 19:55|測試29號13:30送出
		case "每天", "每日", "every", "every day":
			templateJSONByte, err := linebot.NewTextMessage(replyText).MarshalJSON()
			if err != nil {
				return nil, err
			}
			templateJSON := string(templateJSONByte)
			ln := &bot.LineNotification{
				Service:      "Messaging API",
				PushDateTime: "every day",
				Limit:        -1,
				UserID:       lineId.UserID,
				GroupID:      lineId.GroupID,
				RoomID:       lineId.RoomID,
				Type:         string(linebot.MessageTypeText),
				Template:     templateJSON,
			}
			_, err = ln.Add()
			if err != nil {
				return nil, err
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
	return linebot.NewTextMessage(text), nil
}

func generateErrorTextMessage() linebot.Message {
	return linebot.NewTextMessage("系統錯誤，請重新再試或是通知管理員")
}
