package bot

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"
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
	case "所有提醒", "所有通知", "All TODO":
		return getTodo(lineId)
	case "提醒", "通知", "TODO":
		return todo(lineId, text)
	}
	return linebot.NewTextMessage(text), nil
}

func UserParsePostBackGenTemplate(lineId LineID, postBack *linebot.Postback) interface{} {
	data := []byte(postBack.Data)
	lpba := LinePostBackAction{}
	err := json.Unmarshal(data, &lpba)
	if err != nil {
		log.Fatalf("line group GroupParsePostBackGenTemplate json unmarshal error: %v", err)
	}

	switch lpba.Action {
	case "delete line notification":
		return deleteTodoByPostBack(&lpba)
	}
	return nil
}
