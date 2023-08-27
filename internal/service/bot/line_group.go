package bot

import (
	"fmt"
	"time"

	"github.com/programzheng/black-key/config"
	"github.com/programzheng/black-key/internal/helper"

	"github.com/line/line-bot-sdk-go/linebot"
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
	lpba := createLinePostBackActionByDataAndParams([]byte(postBack.Data))
	if postBack.Params != nil {
		lpba.Params = LinePostBackActionParams{
			Date:     postBack.Params.Date,
			Time:     postBack.Params.Time,
			Datetime: postBack.Params.Datetime,
		}
	}

	switch lpba.Action {
	case "delete line notification":
		return deleteTodoByPostBack(lpba)
	case "結算":
		return bill(lineId, lpba)
	case "猜拳":
		return rockPaperScissorTurn(lpba)
	case "todo":
		parsedTime, err := time.Parse("2006-01-02T15:04", lpba.Params.Datetime)
		if err != nil {
			return nil, fmt.Errorf("UserParsePostBackGenTemplate Action:todo error: %v", err)
		}
		text := fmt.Sprintf("提醒|%s|%s", parsedTime.Format("2006-01-02 15:04:05"), lpba.Data["Text"].(string))
		return todo(lineId, text)
	}
	if helper.ConvertToBool(config.Cfg.GetString("LINE_MESSAGING_DEBUG")) {
		return linebot.NewTextMessage(postBack.Data), nil
	}
	return nil, nil
}
