package bot

import (
	"encoding/json"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/internal/model/bot"
)

func createLineNotification(lineId LineID, pushCycle string, pushDateTime time.Time, limit int, replyText string) (*bot.LineNotification, error) {
	templates := []interface{}{}
	templates = append(templates, generateLineMessagingTemplate(replyText))
	templatesJSONByte, err := json.Marshal(templates)
	if err != nil {
		return nil, err
	}
	templatesJSON := string(templatesJSONByte)

	return createLineNotificationByTemplatesJSON(lineId, pushCycle, pushDateTime, limit, templatesJSON)
}

func createLineNotificationByTemplatesJSON(lineId LineID, pushCycle string, pushDateTime time.Time, limit int, templatesJSON string) (*bot.LineNotification, error) {
	ln := &bot.LineNotification{
		Service:      "Messaging API",
		PushCycle:    pushCycle,
		PushDateTime: pushDateTime,
		Limit:        1,
		UserID:       lineId.UserID,
		GroupID:      lineId.GroupID,
		RoomID:       lineId.RoomID,
		Type:         string(linebot.MessageTypeText),
		Template:     templatesJSON,
	}
	result, err := ln.Add()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func generateLineMessagingTemplate(input interface{}) interface{} {
	switch value := input.(type) {
	case string:
		return linebot.NewTextMessage(value)
	}

	return nil
}
