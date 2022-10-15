package bot

import (
	"encoding/json"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/config"
	"github.com/programzheng/black-key/internal/helper"
	log "github.com/sirupsen/logrus"
)

func UserParseTextGenTemplate(lineId LineID, text string) (interface{}, error) {
	parseText := strings.Split(text, "|")

	switch parseText[0] {
	// Line相關ID資訊
	case "資訊":
		return getLineId(lineId)
	case "我的大頭貼":
		return getMemberLineAvatar(lineId)
	case "所有提醒", "所有通知", "All TODO":
		return getTodo(lineId)
	case "提醒", "通知", "TODO":
		return todo(lineId, text)
	}
	if helper.ConvertToBool(config.Cfg.GetString("LINE_MESSAGING_DEBUG")) {
		return linebot.NewTextMessage(text), nil
	}
	return nil, nil
}

func UserParsePostBackGenTemplate(lineId LineID, postBack *linebot.Postback) (interface{}, error) {
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
	if helper.ConvertToBool(config.Cfg.GetString("LINE_MESSAGING_DEBUG")) {
		return linebot.NewTextMessage(string(data)), nil
	}
	return nil, nil
}
