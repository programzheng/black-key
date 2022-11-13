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

	parseText := strings.Split(text, "|")

	if replayResult != nil {
		return replayResult, nil
	}

	//功能說明
	if len(parseText) == 1 {
		switch parseText[0] {
		case "c helper", "記帳說明", "記帳":
			return linebot.NewTextMessage("*記帳*\n將按照群組人數去做平均計算，使用記帳請使用以下格式輸入\n\"記帳|標題|總金額|備註\"\n例如:\n記帳|生日聚餐|1234|本人生日"), nil
		case "c list helper", "記帳列表說明":
			return linebot.NewTextMessage("*記帳列表*\n將回傳記帳紀錄的列表，格式為:\n日期時間 標題|金額| 平均金額 |付款人|備註"), nil
		case "c balance helper", "記帳結算說明", "結算說明":
			return linebot.NewTextMessage("*記帳結算說明*\n將刪除記帳資料，格式為:\n記帳結算|日期(可選)"), nil
		}
	}

	//功能
	switch parseText[0] {
	// Line相關資訊
	case "資訊":
		return getLineId(lineId)
	// c list||記帳列表
	case "c list", "記帳列表":
		return getLineBillings(lineId)
	// c||記帳|生日聚餐|1234|本人生日
	case "c", "記帳":
		return createBilling(lineId, text)
	// 記帳結算
	case "記帳結算", "結帳", "結算":
		return getBills(lineId, text)
	case "我的大頭貼":
		return getGroupMemberLineAvatar(lineId)
	case "猜拳", "石頭布剪刀", "剪刀石頭布", "rock-paper-scissors":
		return startRockPaperScissor(lineId, text)
	case "所有提醒", "所有通知", "All TODO":
		return getTodo(lineId)
	case "提醒", "通知", "TODO":
		return todo(lineId, text)
	case "多提醒", "多通知":
		return todos(lineId, text)
	}

	if helper.ConvertToBool(config.Cfg.GetString("LINE_MESSAGING_DEBUG")) {
		return linebot.NewTextMessage(text), nil
	}
	return nil, nil
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
