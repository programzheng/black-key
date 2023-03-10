package bot

import (
	"encoding/json"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/internal/model/bot"
)

func createLineNotificationByText(
	lineId LineID,
	pushCycle string,
	pushDateTime time.Time,
	limit int,
	replyText string,
) (*bot.LineNotification, error) {
	templates := []interface{}{}
	templates = append(templates, linebot.NewTextMessage(replyText))
	templatesJSONByte, err := json.Marshal(templates)
	if err != nil {
		return nil, err
	}
	templatesJSON := string(templatesJSONByte)

	return createLineNotificationByTemplatesJSON(
		lineId,
		pushCycle,
		pushDateTime,
		limit,
		string(linebot.MessageTypeText),
		templatesJSON,
	)
}

func createLineNotificationByTemplatesJSON(
	lineId LineID,
	pushCycle string,
	pushDateTime time.Time,
	limit int,
	t string,
	templatesJSON string,
) (*bot.LineNotification, error) {
	ln := &bot.LineNotification{
		Service:      "Messaging API",
		PushCycle:    pushCycle,
		PushDateTime: pushDateTime,
		Limit:        limit,
		UserID:       lineId.UserID,
		GroupID:      lineId.GroupID,
		RoomID:       lineId.RoomID,
		Type:         t,
		Template:     templatesJSON,
	}
	result, err := ln.Add()
	if err != nil {
		return nil, err
	}
	return result, nil
}
