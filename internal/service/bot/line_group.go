package bot

import (
	"encoding/json"
	"strings"

	"github.com/programzheng/black-key/config"
	"github.com/programzheng/black-key/internal/helper"
	"github.com/programzheng/black-key/internal/model/bot"

	underscore "github.com/ahl5esoft/golang-underscore"
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

	parseText := strings.Split(text, "|")

	strategies := []TextParsingStrategy{
		&InfoStrategy{},
		&BillingStrategy{},
		&GroupMemberLineAvatarStrategy{},
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

func getDistinctByUserID(lbs []bot.LineBilling) map[string]string {
	//user id line member display name
	dstByUserID := make(map[string]string, 0)
	underscore.Chain(lbs).DistinctBy("UserID").SelectMany(func(lb bot.LineBilling, _ int) map[string]string {
		dst := make(map[string]string)
		lineMember, err := botClient.GetGroupMemberProfile(lb.GroupID, lb.UserID).Do()
		if err != nil {
			dst[lb.UserID] = "Unknow"
			return dst
		}
		dst[lb.UserID] = lineMember.DisplayName
		return dst
	}).Value(&dstByUserID)

	return dstByUserID
}

func getGroupMemberCount(groupID string) int {
	groupMemberCount, err := botClient.GetGroupMemberCount(groupID).Do()
	if err != nil {
		log.Fatal("line messaging api get group member count error:", err)
	}
	return groupMemberCount.Count
}
func calculateAmount(groupID string, amount float64) (float64, int) {
	//預設平均計算基數
	amountAvgBase := 3.0
	groupMemberCount := getGroupMemberCount(groupID)
	amountAvgBase = helper.ConvertToFloat64(groupMemberCount)
	amountAvg := amount / amountAvgBase
	return amountAvg, groupMemberCount
}
