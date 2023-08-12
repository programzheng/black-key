package bot

import (
	"encoding/json"

	"github.com/programzheng/black-key/config"
	"github.com/programzheng/black-key/internal/helper"

	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"
)

func GroupParseTextGenTemplate(lineId LineID, text string) (interface{}, error) {
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
		&GroupMemberLineAvatarStrategy{},
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

func GroupHandleReceiveImageMessage(
	lineId *LineID,
	messageContentResponse *linebot.MessageContentResponse,
) (interface{}, error) {
	//before handle
	replayResult, err := replayBeforeHandle(lineId, messageContentResponse)
	if err != nil {
		return nil, err
	}
	if replayResult != nil {
		return replayResult, nil
	}

	fs, staticFile := messageContentResponseToStaticFile(messageContentResponse)
	return linebot.NewImageMessage(
		fs.GetHostURL()+"/"+staticFile.Name,
		fs.GetHostURL()+"/"+staticFile.Name,
	), nil
}

func GroupParsePostBackGenTemplate(lineId LineID, postBack *linebot.Postback) (interface{}, error) {
	data := []byte(postBack.Data)
	lpba := LinePostBackAction{}
	err := json.Unmarshal(data, &lpba)
	if err != nil {
		log.Fatalf("line group GroupParsePostBackGenTemplate json unmarshal error: %v", err)
	}

	switch lpba.Action {
	case "delete line notification":
		return deleteTodoByPostBack(&lpba)
	case "結算":
		return bill(lineId, &lpba)
	case "猜拳":
		return rockPaperScissorTurn(&lpba)
	}
	if helper.ConvertToBool(config.Cfg.GetString("LINE_MESSAGING_DEBUG")) {
		return linebot.NewTextMessage(string(data)), nil
	}
	return nil, nil
}
