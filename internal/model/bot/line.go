package bot

import (
	"github.com/programzheng/black-key/internal/database"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"gorm.io/gorm"
)

type LineBotRequest struct {
	gorm.Model
	Type       string `bson:"type"`
	GroupID    string `bson:"group_id"`
	RoomID     string `bson:"room_id"`
	UserID     string `bson:"user_id"`
	ReplyToken string `bson:"reply_token"`
	Request    string `bson:"request"`
}

const MongoCollection = "line_bot_request"

func (lbr *LineBotRequest) Create() (*string, error) {
	r, err := database.NewMongoInstance().CreateOne(MongoCollection, lbr)
	if err != nil {
		return nil, err
	}
	id := r.InsertedID.(primitive.ObjectID).Hex()
	return &id, nil
}
