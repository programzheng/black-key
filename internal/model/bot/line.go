package bot

import (
	"github.com/programzheng/black-key/internal/database"
)

type LineBotRequest struct {
	Type       string `bson:"type"`
	GroupID    string `bson:"group_id"`
	RoomID     string `bson:"room_id"`
	UserID     string `bson:"user_id"`
	ReplyToken string `bson:"reply_token"`
	Request    string `bson:"request"`
}

type LineBotRequestRepository struct {
	Model *LineBotRequest
	database.MongoBaseRepository
}

const mongoCollectionName = "line_bot_requests"

func NewLineBotRequestRepository() *LineBotRequestRepository {
	mbr := *database.NewMongoBaseRepository()
	mbr.CollectionName = mongoCollectionName
	lbrr := &LineBotRequestRepository{
		Model:               &LineBotRequest{},
		MongoBaseRepository: mbr,
	}
	return lbrr
}
