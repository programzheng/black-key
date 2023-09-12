package bot

import (
	"encoding/json"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/internal/model/bot"
	"github.com/programzheng/black-key/internal/service/rent_house"
)

const (
	lineFeatureNotificationFeatureNewRentHomes = "new_rent_homes"
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

func GetFlexMessageByLineFeatureNotification(lfn *bot.LineFeatureNotification) ([]linebot.SendingMessage, error) {
	switch lfn.Feature {
	case lineFeatureNotificationFeatureNewRentHomes:
		grhcs, err := rent_house.GetGetRentHousesConditionsByJSONString(lfn.Request)
		if err != nil {
			return nil, err
		}
		grhr, err := rent_house.GetRentHousesByConditionsResponse(grhcs)
		if err != nil {
			return nil, err
		}
		rhs, err := rent_house.ConvertGetRentHousesResponseToRentHouses(grhr)
		if err != nil {
			return nil, err
		}
		return NewNewRentHousesFlexTemplate(grhcs.City, rhs), nil
	}

	return nil, nil
}
