package bot

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/internal/helper"
	"github.com/programzheng/black-key/internal/model/bot"
	log "github.com/sirupsen/logrus"
)

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
		var tp linebot.TextMessage
		data := []byte(ln.Template)
		err := json.Unmarshal(data, &tp)
		if err != nil {
			log.Printf("pkg/service/bot/line_messaging getTodo json.Unmarshal error: %v", err)
			return nil, err
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

		templateJSONByte, err := linebot.NewTextMessage(replyText).MarshalJSON()
		if err != nil {
			return generateErrorTextMessage(), err
		}
		weekDays := strings.Join(helper.GetWeekDays(), ",")
		pdtl := *tt
		templateJSON := string(templateJSONByte)
		ln := &bot.LineNotification{
			Service:      "Messaging API",
			PushCycle:    weekDays,
			PushDateTime: pdtl,
			Limit:        -1,
			UserID:       lineId.UserID,
			GroupID:      lineId.GroupID,
			RoomID:       lineId.RoomID,
			Type:         string(linebot.MessageTypeText),
			Template:     templateJSON,
		}
		_, err = ln.Add()
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
		templateJSONByte, err := linebot.NewTextMessage(replyText).MarshalJSON()
		if err != nil {
			return generateErrorTextMessage(), err
		}
		pdtl := *tt
		weekDays := strings.Join(wdens, ",")
		templateJSON := string(templateJSONByte)
		ln := &bot.LineNotification{
			Service:      "Messaging API",
			PushCycle:    weekDays,
			PushDateTime: pdtl,
			Limit:        -1,
			UserID:       lineId.UserID,
			GroupID:      lineId.GroupID,
			RoomID:       lineId.RoomID,
			Type:         string(linebot.MessageTypeText),
			Template:     templateJSON,
		}
		_, err = ln.Add()
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

	pdtl := dtt
	templateJSONByte, err := linebot.NewTextMessage(replyText).MarshalJSON()
	if err != nil {
		return generateErrorTextMessage(), err
	}
	templateJSON := string(templateJSONByte)
	ln := &bot.LineNotification{
		Service:      "Messaging API",
		PushCycle:    "specify",
		PushDateTime: pdtl,
		Limit:        1,
		UserID:       lineId.UserID,
		GroupID:      lineId.GroupID,
		RoomID:       lineId.RoomID,
		Type:         string(linebot.MessageTypeText),
		Template:     templateJSON,
	}
	_, err = ln.Add()
	if err != nil {
		return generateErrorTextMessage(), err
	}

	return linebot.NewTextMessage("設置完成將於" + date + "\n傳送訊息:" + replyText), nil

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

func deleteTodoByPostBack(lpba *LinePostBackAction) interface{} {
	id := uint(lpba.Data["ID"].(float64))
	ln, err := bot.LineNotificationFirstByID(id)
	if err != nil {
		return nil
	}
	err = ln.Delete()
	if err != nil {
		return linebot.NewTextMessage(
			"刪除失敗",
		)
	}

	return linebot.NewTextMessage("刪除成功")
}

func startRockPaperScissor(lineId LineID, text string) (interface{}, error) {
	groupMemberCount := getGroupMemberCount(lineId.GroupID)
	// if groupMemberCount <= 1 {
	// 	return linebot.NewTextMessage("此功能需要群組大於(包含)2人"), nil
	// }
	key := "rock-paper-scissors-" + lineId.GroupID
	minutes := "5"
	m, _ := time.ParseDuration(minutes + "m")
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

func rockPaperScissorTurn(lpba *LinePostBackAction) interface{} {
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
				winCount := conditionRockPaperScissors(oldAction, es, numberOfPeople)
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
	return nil
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
