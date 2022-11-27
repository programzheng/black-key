package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/internal/cache"
)

func replayBeforeHandle(lineId *LineID, input interface{}) (interface{}, error) {
	if checkTodosExist(lineId) {
		return appendTodos(lineId, input)
	}
	return nil, nil
}

func checkTodosExist(lineId *LineID) bool {
	rdb := cache.GetRedisClient()
	ctx := context.Background()
	exist := rdb.Exists(ctx, lineId.getTodosCacheKey()).Val()

	return exist != 0
}

func appendTodos(lineId *LineID, input interface{}) (interface{}, error) {
	rdb := cache.GetRedisClient()
	ctx := context.Background()
	templatesJSON := rdb.HGet(ctx, lineId.getTodosCacheKey(), "templates").Val()
	templates := []interface{}{}
	err := json.Unmarshal([]byte(templatesJSON), &templates)
	if err != nil {
		return generateErrorTextMessage(), err
	}
	replayText := ""
	switch value := input.(type) {
	case string:
		if value == "結束" {
			dts := rdb.HGet(ctx, lineId.getTodosCacheKey(), "date_time").Val()
			dtt, err := time.ParseInLocation("2006-01-02 15:04:05", dts, time.Now().Local().Location())
			if err != nil {
				return generateErrorTextMessage(), err
			}
			_, err = createLineNotificationByTemplatesJSON(*lineId, "specify", dtt, -1, templatesJSON)
			if err != nil {
				return generateErrorTextMessage(), err
			}
			err = rdb.HDel(ctx, lineId.getTodosCacheKey(), "date_time", "templates").Err()
			if err != nil {
				return generateErrorTextMessage(), err
			}
			err = rdb.Del(ctx, lineId.getTodosCacheKey()).Err()
			if err != nil {
				return generateErrorTextMessage(), err
			}
			return linebot.NewTextMessage("結束設置多通知"), nil
		}
		templates = append(templates, linebot.NewTextMessage(value))
		replayText = fmt.Sprintf("完成設定文字通知:%s，請輸入\"結束\"進行儲存", value)
	case *linebot.MessageContentResponse:
		imageMessage, err := getImageMessageAppendToTodos(lineId, input.(*linebot.MessageContentResponse))
		if err != nil {
			return generateErrorTextMessage(), err
		}
		templates = append(templates, imageMessage)
		replayText = "完成設定圖片通知，請輸入\"結束\"進行儲存"
	}
	b, err := json.Marshal(templates)
	if err != nil {
		return generateErrorTextMessage(), err
	}
	err = rdb.HSet(ctx, lineId.getTodosCacheKey(), "templates", string(b)).Err()
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
