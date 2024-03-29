package bot

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/internal/cache"
	"github.com/programzheng/black-key/internal/helper"
)

func replayBeforeHandle(lineId *LineID, input interface{}) (interface{}, error) {
	if checkTodosExist(lineId) {
		return appendTodos(lineId, input)
	}
	return nil, nil
}

func checkTodosExist(lineId *LineID) bool {
	cd, err := cache.GetCacheDriver("")
	if err != nil {
		return false
	}
	exist, err := cd.Exists(lineId.getTodosCacheKey())
	return exist != 0 && err == nil
}

func appendTodos(lineId *LineID, input interface{}) (interface{}, error) {
	cd, err := cache.GetCacheDriver("")
	if err != nil {
		return generateErrorTextMessage(), err
	}
	templatesJSON, _ := cd.HGet(lineId.getTodosCacheKey(), "templates")
	templates := []interface{}{}
	err = json.Unmarshal([]byte(templatesJSON), &templates)
	if err != nil {
		return generateErrorTextMessage(), err
	}
	replayText := ""
	switch value := input.(type) {
	case string:
		if value == "結束" {
			date, _ := cd.HGet(lineId.getTodosCacheKey(), "date_time")
			if helper.IsDateTime(date) {
				dtt, err := time.ParseInLocation("2006-01-02 15:04:05", date, time.Now().Local().Location())
				if err != nil {
					return generateErrorTextMessage(), err
				}
				_, err = createLineNotificationByTemplatesJSON(
					*lineId,
					"specify",
					dtt,
					-1,
					"multi",
					templatesJSON,
				)
				if err != nil {
					return generateErrorTextMessage(), err
				}
			} else {
				dtt, err := helper.GetDateTimeByTraditionalChinese(date)
				if err != nil {
					return generateErrorTextMessage(), err
				}
				if dtt.IsZero() {
					return generateErrorTextMessage(), fmt.Errorf(
						"please setting the date time",
					)
				}
				shortTc := strings.Split(date, " ")[0]
				if helper.ShortDateIsEveryDay(shortTc) {
					weekDays := strings.Join(helper.GetWeekDays(), ",")
					_, err = createLineNotificationByTemplatesJSON(
						*lineId,
						weekDays,
						dtt,
						-1,
						"multi",
						templatesJSON,
					)
					if err != nil {
						return generateErrorTextMessage(), err
					}
				} else {
					_, err = createLineNotificationByTemplatesJSON(
						*lineId,
						"specify",
						dtt,
						1,
						"multi",
						templatesJSON,
					)
					if err != nil {
						return generateErrorTextMessage(), err
					}
				}
			}

			_, err = cd.HDel(lineId.getTodosCacheKey(), "date_time", "templates")
			if err != nil {
				return generateErrorTextMessage(), err
			}
			_, err = cd.Del(lineId.getTodosCacheKey())
			if err != nil {
				return generateErrorTextMessage(), err
			}
			return linebot.NewTextMessage("結束設置多通知"), nil
		}
		templates = append(templates, linebot.NewTextMessage(value))
		replayText = fmt.Sprintf("完成設定文字通知:%s\n請輸入\"結束\"進行儲存", value)
	case *linebot.MessageContentResponse:
		imageMessage, err := getImageMessageAppendToTodos(lineId, input.(*linebot.MessageContentResponse))
		if err != nil {
			return generateErrorTextMessage(), err
		}
		templates = append(templates, imageMessage)
		replayText = "完成設定圖片通知\n請輸入\"結束\"進行儲存"
	}
	b, err := json.Marshal(templates)
	if err != nil {
		return generateErrorTextMessage(), err
	}
	_, err = cd.HSet(lineId.getTodosCacheKey(), "templates", string(b))
	if err != nil {
		return generateErrorTextMessage(), err
	}
	return linebot.NewTextMessage(replayText), nil
}

func getImageMessageAppendToTodos(
	lineId *LineID,
	messageContentResponse *linebot.MessageContentResponse,
) (interface{}, error) {
	fs, staticFile := messageContentResponseToStaticFile(messageContentResponse)

	return linebot.NewImageMessage(
		fs.GetHostURL()+"/"+staticFile.Name,
		fs.GetHostURL()+"/"+staticFile.Name,
	), nil
}
