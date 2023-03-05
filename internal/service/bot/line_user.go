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
	//before handle
	replayResult, err := replayBeforeHandle(&lineId, text)
	if err != nil {
		return nil, err
	}
	if replayResult != nil {
		return replayResult, nil
	}

	parseText := strings.Split(text, "|")

	strategies := []TextParsingStrategy{
		&InfoStrategy{},
		&BillingStrategy{},
		&MemberLineAvatarStrategy{},
		&RockPaperScissorStrategy{},
		&TodoStrategy{},
		&DefaultStrategy{},
	}
	actionText := parseText[0]

	for _, strategy := range strategies {
		result, err := strategy.Execute(lineId, actionText)
		if err != nil {
			return nil, err
		}
		if result != nil {
			return result, nil
		}
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
