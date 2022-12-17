package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/i18n"
	"github.com/programzheng/black-key/internal/cache"
	"github.com/programzheng/black-key/internal/helper"
	"github.com/programzheng/black-key/internal/model"
	"github.com/programzheng/black-key/internal/model/bot"
	log "github.com/sirupsen/logrus"
)

func getLineIdMap(lineId LineID) map[string]interface{} {
	lineIdMap := make(map[string]interface{})
	lineIdMap["room_id"] = lineId.RoomID
	lineIdMap["group_id"] = lineId.GroupID
	lineIdMap["user_id"] = lineId.UserID

	return lineIdMap
}

func getLineId(lineId LineID) (interface{}, error) {
	return linebot.NewTextMessage(fmt.Sprintf("RoomID:%v\nGroupID:%v\nUserID:%v", lineId.RoomID, lineId.GroupID, lineId.UserID)), nil
}

func getMemberLineAvatar(lineId LineID) (interface{}, error) {
	lineMember, err := botClient.GetProfile(lineId.UserID).Do()
	if err != nil {
		return generateErrorTextMessage(), err
	}
	return linebot.NewImageMessage(lineMember.PictureURL, lineMember.PictureURL), nil
}

func getGroupMemberLineAvatar(lineId LineID) (interface{}, error) {
	lineMember, err := botClient.GetGroupMemberProfile(lineId.GroupID, lineId.UserID).Do()
	if err != nil {
		return generateErrorTextMessage(), err
	}
	return linebot.NewImageMessage(lineMember.PictureURL, lineMember.PictureURL), nil
}

func setTodoHelper() interface{} {
	t := &(i18n.Translation{})
	s := t.Translate("LINE_Messaging_Todo_Notification_Helper")

	return linebot.NewTextMessage(s)
}

func setTodosHelper() interface{} {
	t := &(i18n.Translation{})
	s := t.Translate("LINE_Messaging_Todos_Notification_Helper")

	return linebot.NewTextMessage(s)
}

func getTodo(lineId LineID) (interface{}, error) {
	lns, err := (&bot.LineNotification{}).Get(map[string]interface{}{
		"user_id":  lineId.UserID,
		"group_id": lineId.GroupID,
		"room_id":  lineId.RoomID,
	}, nil)
	if err != nil {
		return nil, err
	}
	if len(lns) == 0 {
		return linebot.NewTextMessage("沒有資料"), nil
	}
	carouselColumns := []*linebot.CarouselColumn{}

	for _, ln := range lns {
		tps := []interface{}{}
		err := json.Unmarshal([]byte(ln.Template), &tps)
		if err != nil {
			log.Printf("internal/service/bot/line_text_messaging getTodo tps json.Unmarshal error: %v", err)
		}
		textTemplate, err := json.Marshal(tps[0])
		if err != nil {
			log.Printf("internal/service/bot/line_text_messaging getTodo first tps json.Marshal error: %v", err)
		}
		var tp linebot.TextMessage
		err = json.Unmarshal(textTemplate, &tp)
		if err != nil {
			log.Printf("pkg/service/bot/line_messaging getTodo tp json.Unmarshal error: %v", err)
		}
		deletePostBackAction := LinePostBackAction{
			Action: "delete line notification",
			Data: map[string]interface{}{
				"ID": ln.ID,
			},
		}
		deletePostBackActionJson, err := json.Marshal(deletePostBackAction)
		if err != nil {
			log.Printf("pkg/service/bot/line_messaging getTodo deletePostBackActionJson json.Marshal error: %v", err)
			return nil, err
		}
		pushDateTime := ln.PushDateTime.Local().Format(helper.Yyyymmddhhmmss)
		pushCycleString := func(ln *bot.LineNotification) string {
			r := ""
			switch ln.PushCycle {
			case "specify":
				r = "指定時間"
			default:
				var buf strings.Builder
				pcs := strings.Split(ln.PushCycle, ",")
				buf.WriteString("星期")
				for i, pc := range pcs {
					buf.WriteString(helper.GetWeekDayShortTraditionalChineseByEnglish(pc))
					if i != len(pcs)-1 {
						buf.WriteString("、")
					}
				}
				r = buf.String()
			}

			return r
		}(ln)
		title := fmt.Sprintf(
			"%d, %s",
			ln.ID,
			tp.Text,
		)
		text := fmt.Sprintf(
			"發送週期:%s \n下次發送時間:%s",
			pushCycleString,
			pushDateTime,
		)
		carouselColumn := linebot.NewCarouselColumn(
			"",
			title,
			text,
			linebot.NewPostbackAction(
				"刪除",
				string(deletePostBackActionJson),
				"",
				"",
			),
		)
		carouselColumns = append(carouselColumns, carouselColumn)
	}
	messages := []linebot.SendingMessage{}
	chunkSize := 10
	for i := 0; i < len(carouselColumns); i += chunkSize {
		end := i + chunkSize
		if end > len(carouselColumns) {
			end = len(carouselColumns)
		}
		carouselTemplate := linebot.NewCarouselTemplate(carouselColumns[i:end]...)
		templateMessage := linebot.NewTemplateMessage(
			fmt.Sprintf("所有提醒-%d", len(messages)+1),
			carouselTemplate,
		)
		messages = append(messages, templateMessage)
	}

	return messages, nil
}

func convertPushDateTime(pdt string) string {
	s := strings.Split(pdt, "|")
	if len(s) == 1 {
		return pdt
	}
	period := s[0]
	dateTime := s[1]
	switch period {
	case "Sunday,Monday,Tuesday,Wednesday,Thursday,Friday,Saturday":
		return fmt.Sprintf("每天 %s", dateTime)
	}
	return ""
}

func todo(lineId LineID, text string) (interface{}, error) {
	parseText := strings.Split(text, "|")
	if len(parseText) == 1 {
		return setTodoHelper(), nil
	}
	date := parseText[1]
	replyText := parseText[2]
	parseDate := strings.Split(date, " ")

	if len(parseDate) == 0 {
		return generateErrorTextMessage(), nil
	}

	tt, err := getTimeByTimeString(parseDate[1])
	if err != nil {
		return generateErrorTextMessage(), err
	}

	//every day
	if parseDate[0] == "每天" ||
		parseDate[0] == "每日" ||
		parseDate[0] == "every" ||
		parseDate[0] == "every day" ||
		parseDate[0] == "every-day" {
		if len(parseDate) == 1 {
			return linebot.NewTextMessage(
				fmt.Sprintf("需設置指定時間，例如: %s 23:59:59", parseDate[0]),
			), nil
		}

		weekDays := strings.Join(helper.GetWeekDays(), ",")
		_, err := createLineNotificationByText(lineId, weekDays, *tt, -1, replyText)
		if err != nil {
			return generateErrorTextMessage(), err
		}
		return linebot.NewTextMessage("設置完成將於每天" + parseDate[1] + "\n傳送訊息:" + replyText), nil
	}

	//specify weekday
	wdtcs := strings.Split(parseDate[0], ",")
	wdens := []string{}
	for _, wdtc := range wdtcs {
		wden := helper.GetWeekDayEnglishByTraditionalChinese(wdtc)
		if wden == "" {
			break
		}
		wdens = append(wdens, wden)
	}
	if len(wdens) > 0 {
		weekDays := strings.Join(wdens, ",")
		_, err = createLineNotificationByText(lineId, weekDays, *tt, -1, replyText)
		if err != nil {
			return generateErrorTextMessage(), err
		}
		rpmg := fmt.Sprintf(
			"設置完成將於%s%s\n傳送訊息:%s",
			parseDate[0],
			parseDate[1],
			replyText,
		)
		return linebot.NewTextMessage(rpmg), nil

	}

	//specify date time
	if len(parseDate) == 1 {
		return linebot.NewTextMessage(
			fmt.Sprintf("需設置指定時間，例如: %s 2022-01-01 23:59:59", parseDate[0]),
		), nil
	}
	dts := fmt.Sprintf("%s %s", parseDate[0], parseDate[1])
	dtt, err := time.ParseInLocation("2006-01-02 15:04:05", dts, time.Now().Local().Location())
	if err != nil {
		return generateErrorTextMessage(), err
	}
	ccspm := checkCanSettingPushMessage(dtt)
	if !ccspm {
		return linebot.NewTextMessage(
			"請設置未來的時間",
		), nil
	}

	_, err = createLineNotificationByText(lineId, "specify", dtt, -1, replyText)
	if err != nil {
		return generateErrorTextMessage(), err
	}

	return linebot.NewTextMessage("設置完成將於" + date + "\n傳送訊息:" + replyText), nil

}

func todos(lineId LineID, text string) (interface{}, error) {
	parseText := strings.Split(text, "|")
	if len(parseText) == 1 {
		return setTodosHelper(), nil
	}
	dateTime := parseText[1]
	replyText := parseText[2]

	rdb := cache.GetRedisClient()
	ctx := context.Background()
	todosCacheKey := lineId.getTodosCacheKey()
	templates := []interface{}{}
	templates = append(templates, generateLineMessagingTemplate(replyText))
	templatesJSONByte, err := json.Marshal(templates)
	if err != nil {
		return generateErrorTextMessage(), err
	}
	templatesJSON := string(templatesJSONByte)
	err = rdb.HSet(ctx, todosCacheKey, "date_time", dateTime, "templates", templatesJSON).Err()
	if err != nil {
		return generateErrorTextMessage(), err
	}

	return linebot.NewTextMessage("設置將於" + dateTime + "\n傳送標題為:" + replyText + "\n請繼續輸入其他內容(例如:圖片)"), nil
}

func getTimeByTimeString(ts string) (*time.Time, error) {
	dt := fmt.Sprintf("%s %s", helper.GetNowDateTimeByFormat("2006-01-02"), ts)
	pdtl, err := time.ParseInLocation("2006-01-02 15:04:05", dt, time.Now().Local().Location())
	if err != nil {
		return nil, err
	}
	return &pdtl, nil
}

func checkCanSettingPushMessage(t time.Time) bool {
	return time.Now().Before(t)
}

func startRockPaperScissor(lineId LineID, text string) (interface{}, error) {
	groupMemberCount := getGroupMemberCount(lineId.GroupID)
	// if groupMemberCount <= 1 {
	// 	return linebot.NewTextMessage("此功能需要群組大於(包含)2人"), nil
	// }
	key := "rock-paper-scissors-" + lineId.GroupID
	minutes := "5"
	m, _ := time.ParseDuration(minutes + "m")
	rdb := cache.GetRedisClient()
	ctx := context.Background()
	exist := rdb.Exists(ctx, key).Val()
	if exist > 0 {
		return rockPaperScissorsTemplate(lineId, "已有猜拳正在進行中", minutes), nil
	}
	err := rdb.SAdd(ctx, key, groupMemberCount).Err()
	if err != nil {
		log.Fatalf("create a rock-paper-scissors error:%v", err)
	}
	err = rdb.Expire(ctx, key, m).Err()
	if err != nil {
		log.Fatalf("set expire rock-paper-scissors time error:%v", err)
	}
	return rockPaperScissorsTemplate(lineId, "剪刀石頭布", minutes), nil
}

func createBilling(lineId LineID, text string) (interface{}, error) {
	parseText := strings.Split(text, "|")

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
	return linebot.NewTextMessage(title + ":記帳完成," + parseText[2] + "/" + helper.ConvertToString(int(amountAvgBase)) + " = " + "*" + helper.ConvertToString(amountAvg) + "*"), nil
}

func getLineBillings(lineId LineID) (interface{}, error) {
	messages := []linebot.SendingMessage{}

	lineIdMap := getLineIdMap(lineId)
	var lbs []bot.LineBilling
	err := model.DB.Where(lineIdMap).Preload("Billing").Find(&lbs).Error
	if err != nil {
		return generateErrorTextMessage(), err
	}
	//沒有記帳資料
	if len(lbs) == 0 {
		return linebot.NewTextMessage("目前沒有記帳紀錄哦！"), nil
	}
	dstByUserID := getDistinctByUserID(lbs)
	listText := getLineBillingList(lineId, lbs, dstByUserID)
	messages = append(messages, linebot.NewTextMessage(listText))
	totalText := getLineBillingTotalAmount(lineId, lbs, dstByUserID)
	messages = append(messages, linebot.NewTextMessage(totalText))

	return messages, nil
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

func getBills(lineId LineID, text string) (interface{}, error) {
	parseText := strings.Split(text, "|")
	lineIdMap := getLineIdMap(lineId)

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
		return linebot.NewTextMessage(fmt.Sprintf("%v以前沒有記帳紀錄哦！", date)), nil
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

	return messages, nil
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

func conditionRockPaperScissors(target string, all []string, numberOfPeople int) int {
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

// This feature is available only for verified or premium accounts
func getGroupMemberIds(groupID string, continuationToken string) []string {
	groupMemberIds, err := botClient.GetGroupMemberIDs(groupID, continuationToken).Do()
	if err != nil {
		log.Fatal("line messaging api get group member ids error:", err)
	}
	return groupMemberIds.MemberIDs
}
