package bot

import (
	"github.com/programzheng/black-key/internal/model/bot"
	"go.mongodb.org/mongo-driver/bson"
)

type LineBotRequestService struct {
}

func (lbrService *LineBotRequestService) CreateOne(attributes map[string]interface{}) (string, error) {
	m := bot.LineBotRequest{
		Type:       attributes["Type"].(string),
		GroupID:    attributes["GroupID"].(string),
		RoomID:     attributes["RoomID"].(string),
		UserID:     attributes["UserID"].(string),
		ReplyToken: attributes["ReplyToken"].(string),
		Request:    attributes["Request"].(string),
	}
	ID, err := bot.NewLineBotRequestRepository().CreateOne(m)
	if err != nil {
		return "", err
	}
	return *ID, nil
}

func (lbrService *LineBotRequestService) Get(f map[string]interface{}) ([]bot.LineBotRequest, error) {
	filter := bson.D{}
	for k, v := range f {
		be := bson.E{
			Key:   k,
			Value: v,
		}
		filter = append(filter, be)
	}
	lbrsis, err := bot.NewLineBotRequestRepository().Find(filter, []bot.LineBotRequest{})
	if err != nil {
		return nil, err
	}
	lbrs := lbrsis.([]bot.LineBotRequest)
	return lbrs, nil
}
