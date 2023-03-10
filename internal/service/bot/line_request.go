package bot

import (
	"time"

	modelBot "github.com/programzheng/black-key/internal/model/bot"
	"go.mongodb.org/mongo-driver/bson"
)

type LineBotRequestService struct {
}

func (lbrService *LineBotRequestService) CreateOne(attributes map[string]interface{}) (string, error) {
	m := modelBot.LineBotRequest{
		Type:       attributes["Type"].(string),
		GroupID:    attributes["GroupID"].(string),
		RoomID:     attributes["RoomID"].(string),
		UserID:     attributes["UserID"].(string),
		ReplyToken: attributes["ReplyToken"].(string),
		Request:    attributes["Request"].(string),
		CreatedAt:  time.Now(),
	}
	ID, err := modelBot.NewLineBotRequestRepository().CreateOne(m)
	if err != nil {
		return "", err
	}
	return *ID, nil
}

func (lbrService *LineBotRequestService) Get(f map[string]interface{}) ([]modelBot.LineBotRequest, error) {
	filter := bson.D{}
	for k, v := range f {
		be := bson.E{
			Key:   k,
			Value: v,
		}
		filter = append(filter, be)
	}
	lbris, err := modelBot.NewLineBotRequestRepository().Find(filter)
	if err != nil {
		return nil, err
	}
	lbrs := lbris.([]modelBot.LineBotRequest)

	return lbrs, nil
}
