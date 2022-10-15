package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	underscore "github.com/ahl5esoft/golang-underscore"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/internal/model"
	"github.com/programzheng/black-key/internal/model/bot"
	"github.com/programzheng/black-key/internal/service/billing"
	log "github.com/sirupsen/logrus"
)

func deleteTodoByPostBack(lpba *LinePostBackAction) (interface{}, error) {
	id := uint(lpba.Data["ID"].(float64))
	ln, err := bot.LineNotificationFirstByID(id)
	if err != nil {
		return nil, err
	}
	err = ln.Delete()
	if err != nil {
		return linebot.NewTextMessage(
			"刪除失敗",
		), nil
	}
	var tp linebot.TextMessage
	data := []byte(ln.Template)
	err = json.Unmarshal(data, &tp)
	if err != nil {
		log.Printf("pkg/service/bot/line_messaging deleteTodoByPostBack json.Unmarshal error: %v", err)
		return nil, err
	}
	text := fmt.Sprintf("刪除ID為%d \"%s\" 的提醒成功", ln.ID, tp.Text)
	return linebot.NewTextMessage(text), nil
}

func bill(lineId LineID, lpba *LinePostBackAction) (interface{}, error) {
	lineUserID := lpba.Data["LineUserID"].(string)
	if lineUserID != lineId.UserID {
		return linebot.NewTextMessage("操作者不同，請自行輸入\"結算\""), nil
	}
	date := lpba.Data["Date"].(string)
	lineIdMap := getLineIdMap(lineId)
	var lbs []bot.LineBilling
	err := model.DB.Where(lineIdMap).Where("updated_at < ?", date).Preload("Billing").Find(&lbs).Error
	if err != nil {
		log.Fatalf("line group GroupParsePostBackGenTemplate 結算 Get LineBilling failed: %v", err)
	}
	memberName := "Unknow"
	lineMember, _ := botClient.GetGroupMemberProfile(lineId.GroupID, lineId.UserID).Do()
	memberName = lineMember.DisplayName
	if len(lbs) == 0 {
		return linebot.NewTextMessage(fmt.Sprintf("%v:%v以前沒有記帳紀錄哦", memberName, date)), nil
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

	return linebot.NewTextMessage(fmt.Sprintf("%v:成功刪除 *%v* 以前的記帳資料", memberName, date)), nil
}

func rockPaperScissorTurn(lpba *LinePostBackAction) (interface{}, error) {
	lineGroupID := lpba.Data["LineGroupID"].(string)
	lineUserID := lpba.Data["LineUserID"].(string)
	key := "rock-paper-scissors-" + lineGroupID
	exist := rdb.Exists(ctx, key).Val()
	if exist == 0 {
		return linebot.NewTextMessage("請輸入\"猜拳\"開始賽局"), nil
	}
	action := lpba.Data["Action"].(string)
	if ok, _ := rdb.SIsMember(ctx, key, lineUserID+"-out").Result(); ok {
		memberName := "Unknow"
		lineMember, _ := botClient.GetGroupMemberProfile(lineGroupID, lineUserID).Do()
		memberName = lineMember.DisplayName
		return linebot.NewTextMessage(memberName + "已出局"), nil
	}
	if ok, _ := rdb.SIsMember(ctx, key, lineUserID+"-rock").Result(); ok {
		memberName := "Unknow"
		lineMember, _ := botClient.GetGroupMemberProfile(lineGroupID, lineUserID).Do()
		memberName = lineMember.DisplayName
		return linebot.NewTextMessage(memberName + "已出過"), nil
	}
	if ok, _ := rdb.SIsMember(ctx, key, lineUserID+"-paper").Result(); ok {
		memberName := "Unknow"
		lineMember, _ := botClient.GetGroupMemberProfile(lineGroupID, lineUserID).Do()
		memberName = lineMember.DisplayName
		return linebot.NewTextMessage(memberName + "已出過"), nil
	}
	if ok, _ := rdb.SIsMember(ctx, key, lineUserID+"-scissors").Result(); ok {
		memberName := "Unknow"
		lineMember, _ := botClient.GetGroupMemberProfile(lineGroupID, lineUserID).Do()
		memberName = lineMember.DisplayName
		return linebot.NewTextMessage(memberName + "已出過"), nil
	}
	es, err := rdb.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("get a rock-paper-scissors set error:%v", err))
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
						return nil, errors.New(fmt.Sprintf("rock-paper-scissors out rem error:%v", err))
					}
					err = rdb.SAdd(ctx, key, oldUserId+"-out").Err()
					if err != nil {
						return nil, errors.New(fmt.Sprintf("rock-paper-scissors out add error:%v", err))
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
						return nil, errors.New(fmt.Sprintf("rock-paper-scissors rem error:%v", err))
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
				return nil, errors.New(fmt.Sprintf("rock-paper-scissors is end error:%v", err))
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
		return messages, nil
	}
	err = rdb.SAdd(ctx, key, lineUserID+"-"+action).Err()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("create a rock-paper-scissors error:%v", err))
	}
	return nil, nil
}
