package bot

import (
	"encoding/json"

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

	strategies := []TextParsingStrategy{
		&HelpStrategy{},
		&InfoStrategy{},
		&BillingStrategy{},
		&MemberLineAvatarStrategy{},
		&RockPaperScissorStrategy{},
		&TodoStrategy{},
		&ProxyStrategy{},
		&DefaultStrategy{},
	}

	for _, strategy := range strategies {
		result, err := strategy.Execute(lineId, text)
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
		log.Fatalf("line user UserParsePostBackGenTemplate json unmarshal error: %v", err)
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
