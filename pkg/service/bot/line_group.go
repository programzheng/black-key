package bot

import (
	"black-key/config"
	"black-key/pkg/cache"
	"black-key/pkg/helper"
	"black-key/pkg/library/line/bot/template"
	"black-key/pkg/model"
	"black-key/pkg/model/bot"
	"black-key/pkg/service/billing"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	underscore "github.com/ahl5esoft/golang-underscore"
	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()
var rdb = cache.GetRedisClient()

func GroupParseTextGenTemplate(lineId LineID, text string) interface{} {
	parseText := strings.Split(text, "|")

	//功能說明
	if len(parseText) == 1 {
		switch parseText[0] {
		case "c helper", "記帳說明", "記帳":
			return linebot.NewTextMessage("*記帳*\n將按照群組人數去做平均計算，使用記帳請使用以下格式輸入\n\"記帳|標題|總金額|備註\"\n例如:\n記帳|生日聚餐|1234|本人生日")
		case "c list helper", "記帳列表說明":
			return linebot.NewTextMessage("*記帳列表*\n將回傳記帳紀錄的列表，格式為:\n日期時間 標題|金額| 平均金額 |付款人|備註")
		case "c balance helper", "記帳結算說明", "結算說明":
			return linebot.NewTextMessage("*記帳結算說明*\n將刪除記帳資料，格式為:\n記帳結算|日期(可選)")
		}
	}

	lineIdMap := getLineIDMap(lineId)
	//功能
	switch parseText[0] {
	// Line相關資訊
	case "資訊":
		return linebot.NewTextMessage(fmt.Sprintf("RoomID:%v\nGroupID:%v\nUserID:%v", lineId.RoomID, lineId.GroupID, lineId.UserID))
	// c list||記帳列表
	case "c list", "記帳列表":
		messages := []linebot.SendingMessage{}

		var lbs []bot.LineBilling
		err := model.DB.Where(lineIdMap).Preload("Billing").Find(&lbs).Error
		if err != nil {
			log.Fatalf("Get failed: %v", err)
		}
		//沒有記帳資料
		if len(lbs) == 0 {
			return linebot.NewTextMessage("目前沒有記帳紀錄哦！")
		}
		dstByUserID := getDistinctByUserID(lbs)
		listText := getLineBillingList(lineId, lbs, dstByUserID)
		messages = append(messages, linebot.NewTextMessage(listText))
		totalText := getLineBillingTotalAmount(lineId, lbs, dstByUserID)
		messages = append(messages, linebot.NewTextMessage(totalText))

		return messages
	// c||記帳|生日聚餐|1234|本人生日
	case "c", "記帳":
		title := parseText[1]
		amount := helper.ConvertToInt(parseText[2])
		note := ""
		//如果有輸入備註
		if len(parseText) == 4 {
			note = parseText[3]
		}
		billingAction(lineId, amount, title, note)
		amountFloat64 := helper.ConvertToFloat64(amount)
		amountAvg, amountAvgBase := calculateAmount(lineId.GroupID, amountFloat64)
		return linebot.NewTextMessage(title + ":記帳完成," + parseText[2] + "/" + helper.ConvertToString(int(amountAvgBase)) + " = " + "*" + helper.ConvertToString(amountAvg) + "*")
	// 記帳結算
	case "記帳結算", "結帳", "結算":
		messages := []linebot.SendingMessage{}

		date := time.Now().Format(helper.Yyyymmddhhmmss)
		//如果有輸入限制日期
		if len(parseText) == 2 {
			date = parseText[1]
		}
		var lbs []bot.LineBilling
		err := model.DB.Where(lineIdMap).Where("updated_at < ?", date).Preload("Billing").Find(&lbs).Error
		if err != nil {
			log.Fatalf("Get failed: %v", err)
		}
		//沒有記帳資料
		if len(lbs) == 0 {
			return linebot.NewTextMessage(fmt.Sprintf("%v以前沒有記帳紀錄哦！", date))
		}
		dstByUserID := getDistinctByUserID(lbs)
		listText := getLineBillingList(lineId, lbs, dstByUserID)
		messages = append(messages, linebot.NewTextMessage(listText))

		//template
		postBack := LinePostBackAction{
			Action: "結算",
			Data: map[string]interface{}{
				"LineRoomID":  lineId.RoomID,
				"LineGroupID": lineId.GroupID,
				"LineUserID":  lineId.UserID,
				"Date":        date,
			},
		}
		postBackJson, err := json.Marshal(postBack)
		if err != nil {
			log.Fatalf("Marshal failed: %v", err)
		}
		leftBtn := linebot.NewPostbackAction("是", string(postBackJson), "", "")
		rightBtn := linebot.NewMessageAction("否", "記帳列表")

		confirmTemplate := linebot.NewConfirmTemplate("確定要刪除以上紀錄?", leftBtn, rightBtn)
		messages = append(messages, linebot.NewTemplateMessage("確定要刪除以上紀錄?", confirmTemplate))

		return messages
	case "我的大頭貼":
		lineMember, err := botClient.GetGroupMemberProfile(lineId.GroupID, lineId.UserID).Do()
		if err != nil {
			return nil
		}
		return linebot.NewImageMessage(lineMember.PictureURL, lineMember.PictureURL)
	case "猜拳", "石頭布剪刀", "剪刀石頭布", "rock-paper-scissors":
		groupMemberCount := getGroupMemberCount(lineId.GroupID)
		// if groupMemberCount <= 1 {
		// 	return linebot.NewTextMessage("此功能需要群組大於(包含)2人")
		// }
		key := "rock-paper-scissors-" + lineId.GroupID
		minutes := "5"
		m, _ := time.ParseDuration(minutes + "m")
		exist := rdb.Exists(ctx, key).Val()
		if exist > 0 {
			return rockPaperScissorsTemplate(lineId, "已有猜拳正在進行中", minutes)
		}
		err := rdb.SAdd(ctx, key, groupMemberCount).Err()
		if err != nil {
			log.Fatalf("create a rock-paper-scissors error:%v", err)
		}
		err = rdb.Expire(ctx, key, m).Err()
		if err != nil {
			log.Fatalf("set expire rock-paper-scissors time error:%v", err)
		}
		return rockPaperScissorsTemplate(lineId, "剪刀石頭布", minutes)
	case "TODO":
		date := parseText[1]
		replyText := parseText[2]
		parseDate := strings.Split(date, " ")
		switch parseDate[0] {
		case "every":
			// TODO|every 19:55|測試29號13:30送出
			todoAction(lineId.UserID, "every", parseDate[1], template.TODO(replyText))
			return linebot.NewTextMessage("設置完成將於每天" + parseDate[1] + "\n傳送訊息:" + replyText)
		default:
			// TODO|2020/02/29 13:00|測試29號13:30送出
			todoAction(lineId.UserID, "once", date, template.TODO(replyText))
			return linebot.NewTextMessage("設置完成將於" + date + "\n傳送訊息:" + replyText)
		}

	}
	if helper.ConvertToBool(config.Cfg.GetString("LINE_MESSAGING_DEBUG")) {
		return linebot.NewTextMessage("目前沒有此功能")
	}
	return nil
}

func GroupParsePostBackGenTemplate(lineId LineID, postBack *linebot.Postback) interface{} {
	data := []byte(postBack.Data)
	lpba := LinePostBackAction{}
	err := json.Unmarshal(data, &lpba)
	if err != nil {
		log.Fatalf("line group GroupParsePostBackGenTemplate json unmarshal error: %v", err)
	}

	lineIdMap := getLineIDMap(lineId)
	switch lpba.Action {
	case "結算":
		lineUserID := lpba.Data["LineUserID"].(string)
		if lineUserID != lineId.UserID {
			return linebot.NewTextMessage("操作者不同，請自行輸入\"結算\"")
		}
		date := lpba.Data["Date"].(string)
		var lbs []bot.LineBilling
		err := model.DB.Where(lineIdMap).Where("updated_at < ?", date).Preload("Billing").Find(&lbs).Error
		if err != nil {
			log.Fatalf("line group GroupParsePostBackGenTemplate 結算 Get LineBilling failed: %v", err)
		}
		memberName := "Unknow"
		lineMember, _ := botClient.GetGroupMemberProfile(lineId.GroupID, lineId.UserID).Do()
		memberName = lineMember.DisplayName
		if len(lbs) == 0 {
			return linebot.NewTextMessage(fmt.Sprintf("%v:%v以前沒有記帳紀錄哦", memberName, date))
		}
		//delete Billing
		var bID []uint
		underscore.Chain(lbs).SelectBy("BillingID").Value(&bID)
		var bs []billing.Billing
		err = model.DB.Where(bID).Delete(&bs).Error
		if err != nil {
			log.Fatalf("line group GroupParsePostBackGenTemplate 結算 Delete Billing failed: %v", err)
		}

		//delete LineBilling
		err = model.DB.Model(lbs).Delete(&lbs).Error
		if err != nil {
			log.Fatalf("line group GroupParsePostBackGenTemplate 結算 Delete LineBilling failed: %v", err)
		}

		return linebot.NewTextMessage(fmt.Sprintf("%v:成功刪除 *%v* 以前的記帳資料", memberName, date))
	case "猜拳":
		lineGroupID := lpba.Data["LineGroupID"].(string)
		lineUserID := lpba.Data["LineUserID"].(string)
		key := "rock-paper-scissors-" + lineGroupID
		exist := rdb.Exists(ctx, key).Val()
		if exist == 0 {
			return linebot.NewTextMessage("請輸入\"猜拳\"開始賽局")
		}
		action := lpba.Data["Action"].(string)
		if ok, _ := rdb.SIsMember(ctx, key, lineUserID+"-out").Result(); ok {
			memberName := "Unknow"
			lineMember, _ := botClient.GetGroupMemberProfile(lineGroupID, lineUserID).Do()
			memberName = lineMember.DisplayName
			return linebot.NewTextMessage(memberName + "已出局")
		}
		if ok, _ := rdb.SIsMember(ctx, key, lineUserID+"-rock").Result(); ok {
			memberName := "Unknow"
			lineMember, _ := botClient.GetGroupMemberProfile(lineGroupID, lineUserID).Do()
			memberName = lineMember.DisplayName
			return linebot.NewTextMessage(memberName + "已出過")
		}
		if ok, _ := rdb.SIsMember(ctx, key, lineUserID+"-paper").Result(); ok {
			memberName := "Unknow"
			lineMember, _ := botClient.GetGroupMemberProfile(lineGroupID, lineUserID).Do()
			memberName = lineMember.DisplayName
			return linebot.NewTextMessage(memberName + "已出過")
		}
		if ok, _ := rdb.SIsMember(ctx, key, lineUserID+"-scissors").Result(); ok {
			memberName := "Unknow"
			lineMember, _ := botClient.GetGroupMemberProfile(lineGroupID, lineUserID).Do()
			memberName = lineMember.DisplayName
			return linebot.NewTextMessage(memberName + "已出過")
		}
		es, err := rdb.SMembers(ctx, key).Result()
		if err != nil {
			log.Fatalf("get a rock-paper-scissors set error:%v", err)
		}
		numberOfPeople := 4
		//判斷結果
		if len(es) == numberOfPeople {
			messages := []linebot.SendingMessage{}
			es = append(es, lineUserID+"-"+action)
			end := false
			tieCount := 0
			var everyBuilder strings.Builder
			var outBuilder strings.Builder
			var resultBuilder strings.Builder
			for _, s := range es {
				result := strings.Split(s, "-")
				if len(result) > 1 {
					currentMemberName := "Unknow"
					oldUserId := result[0]
					currentLineMember, err := botClient.GetGroupMemberProfile(lineGroupID, oldUserId).Do()
					if err == nil {
						currentMemberName = currentLineMember.DisplayName
					}
					oldAction := result[1]
					winCount := jugdeRockPaperScissors(oldAction, es, numberOfPeople)
					everyBuilder.WriteString(currentMemberName + "出" + convertRockPaperScissors(oldAction) + "\n")
					//出局
					if winCount == 0 {
						err = rdb.SRem(ctx, key, s).Err()
						if err != nil {
							log.Fatalf("rock-paper-scissors out rem error:%v", err)
						}
						err = rdb.SAdd(ctx, key, oldUserId+"-out").Err()
						if err != nil {
							log.Fatalf("rock-paper-scissors out add error:%v", err)
						}
						outBuilder.WriteString(currentMemberName + "出局\n")
						//有獲勝者
					} else if winCount == (numberOfPeople - 1) {
						end = true
						resultBuilder.WriteString("*" + currentMemberName + "獲勝*\n")
					} else {
						tieCount++
						err = rdb.SRem(ctx, key, s).Err()
						if err != nil {
							log.Fatalf("rock-paper-scissors rem error:%v", err)
						}
					}
					//流局
					if tieCount == numberOfPeople {
						end = true
						resultBuilder.WriteString("流局\n")
					}
				}
			}
			if end {
				err = rdb.Del(ctx, key).Err()
				if err != nil {
					log.Fatalf("rock-paper-scissors is end error:%v", err)
				}
			}
			if everyBuilder.Len() > 0 {
				messages = append(messages, linebot.NewTextMessage(strings.TrimSuffix(everyBuilder.String(), "\n")))
			}
			if outBuilder.Len() > 0 {
				messages = append(messages, linebot.NewTextMessage(strings.TrimSuffix(outBuilder.String(), "\n")))
			}
			if resultBuilder.Len() > 0 {
				messages = append(messages, linebot.NewTextMessage(strings.TrimSuffix(resultBuilder.String(), "\n")))
			}
			return messages
		}
		err = rdb.SAdd(ctx, key, lineUserID+"-"+action).Err()
		if err != nil {
			log.Fatalf("create a rock-paper-scissors error:%v", err)
		}
	}
	return nil
}

func getLineIDMap(lineId LineID) map[string]interface{} {
	lineIdMap := make(map[string]interface{})
	lineIdMap["room_id"] = lineId.RoomID
	lineIdMap["group_id"] = lineId.GroupID
	lineIdMap["user_id"] = lineId.UserID

	return lineIdMap
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

func getLineBillingList(lineId LineID, lbs []bot.LineBilling, dstByUserID map[string]string) string {
	var sbList strings.Builder
	sbList.Grow(len(lbs))
	for key, lb := range lbs {
		var memberName string
		amountAvg, amountAvgBase := calculateAmount(lineId.GroupID, helper.ConvertToFloat64(lb.Billing.Amount))
		//check line member display name is exist
		if _, ok := dstByUserID[lb.UserID]; ok {
			memberName = dstByUserID[lb.UserID]
		}
		text := fmt.Sprintf("%v\n%v|%v/%v= *%v* |%v", lb.Billing.CreatedAt.Format(helper.Yyyymmddhhmmss), lb.Billing.Title, helper.ConvertToString(lb.Billing.Amount), helper.ConvertToString(amountAvgBase), helper.ConvertToString(amountAvg), memberName)
		if lb.Billing.Note != "" {
			text = text + "|" + lb.Billing.Note
		}
		if len(lbs)-1 != key {
			text = text + "\n"
		}
		sbList.WriteString(text)
	}
	return string(sbList.String())
}

func getLineBillingTotalAmount(lineId LineID, lbs []bot.LineBilling, dstByUserID map[string]string) string {
	lbUserIDAmount := make(map[string]float64, 0)
	var sbTotal strings.Builder
	sbTotal.Grow(len(dstByUserID))
	for _, lb := range lbs {
		amountAvg, _ := calculateAmount(lineId.GroupID, helper.ConvertToFloat64(lb.Billing.Amount))
		if _, ok := dstByUserID[lb.UserID]; ok {
			lbUserIDAmount[lb.UserID] = lbUserIDAmount[lb.UserID] + amountAvg
		}
	}
	text := "總付款金額：\n"
	sbTotal.WriteString(text)
	for userID, name := range dstByUserID {
		text = fmt.Sprintf("%v: *%v*\n", name, helper.ConvertToString(lbUserIDAmount[userID]))
		sbTotal.WriteString(text)
	}
	return string(sbTotal.String())
}

func getGroupMemberCount(groupID string) int {
	groupMemberCount, err := botClient.GetGroupMemberCount(groupID).Do()
	if err != nil {
		log.Fatal("line messaging api get group member count error:", err)
	}
	return groupMemberCount.Count
}

// This feature is available only for verified or premium accounts
func getGroupMemberIds(groupID string, continuationToken string) []string {
	groupMemberIds, err := botClient.GetGroupMemberIDs(groupID, continuationToken).Do()
	if err != nil {
		log.Fatal("line messaging api get group member ids error:", err)
	}
	return groupMemberIds.MemberIDs
}

func calculateAmount(groupID string, amount float64) (float64, int) {
	//預設平均計算基數
	amountAvgBase := 3.0
	groupMemberCount := getGroupMemberCount(groupID)
	amountAvgBase = helper.ConvertToFloat64(groupMemberCount)
	amountAvg := amount / amountAvgBase
	return amountAvg, groupMemberCount
}

func rockPaperScissorsTemplate(lineId LineID, templateTitle string, minutes string) *linebot.TemplateMessage {
	if minutes == "" {
		minutes = "5"
	}
	rockPostBack := LinePostBackAction{
		Action: "猜拳",
		Data: map[string]interface{}{
			"LineRoomID":  lineId.RoomID,
			"LineGroupID": lineId.GroupID,
			"LineUserID":  lineId.UserID,
			"Action":      "rock",
		},
	}
	rockPostBackJson, err := json.Marshal(rockPostBack)
	if err != nil {
		log.Fatalf("rock post back json failed: %v", err)
	}
	rockBtn := linebot.NewPostbackAction("石頭", string(rockPostBackJson), "", "")
	paperPostBack := LinePostBackAction{
		Action: "猜拳",
		Data: map[string]interface{}{
			"LineRoomID":  lineId.RoomID,
			"LineGroupID": lineId.GroupID,
			"LineUserID":  lineId.UserID,
			"Action":      "paper",
		},
	}
	paperPostBackJson, err := json.Marshal(paperPostBack)
	if err != nil {
		log.Fatalf("paper post back json failed: %v", err)
	}
	paperBtn := linebot.NewPostbackAction("布", string(paperPostBackJson), "", "")
	scissorsPostBack := LinePostBackAction{
		Action: "猜拳",
		Data: map[string]interface{}{
			"LineRoomID":  lineId.RoomID,
			"LineGroupID": lineId.GroupID,
			"LineUserID":  lineId.UserID,
			"Action":      "scissors",
		},
	}
	scissorsPostBackJson, err := json.Marshal(scissorsPostBack)
	if err != nil {
		log.Fatalf("scissors post back json failed: %v", err)
	}
	scissorsBtn := linebot.NewPostbackAction("剪刀", string(scissorsPostBackJson), "", "")
	buttonTemplate := linebot.NewButtonsTemplate("https://images.unsplash.com/photo-1614032686099-e648d6dea9b3", templateTitle, minutes+"分鐘內結束", rockBtn, paperBtn, scissorsBtn)
	return linebot.NewTemplateMessage("開始剪刀石頭布", buttonTemplate)
}

func convertRockPaperScissors(target string) string {
	switch target {
	case "rock":
		return "石頭"
	case "paper":
		return "布"
	case "scissors":
		return "剪刀"
	}
	return "Unknow"
}

func jugdeRockPaperScissors(target string, all []string, numberOfPeople int) int {
	winCount := 0
	for _, s := range all {
		result := strings.Split(s, "-")
		if len(result) > 1 {
			action := result[1]
			switch target {
			case "rock":
				if action == "scissors" {
					winCount++
				}
			case "paper":
				if action == "rock" {
					winCount++
				}
			case "scissors":
				if action == "paper" {
					winCount++
				}
			}
		}
	}
	return winCount
}
