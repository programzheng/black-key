package bot

import (
	"time"

	"github.com/programzheng/black-key/internal/database"
)

type LineBotRequest struct {
	Type       string    `bson:"type"`
	GroupID    string    `bson:"group_id"`
	RoomID     string    `bson:"room_id"`
	UserID     string    `bson:"user_id"`
	ReplyToken string    `bson:"reply_token"`
	Request    string    `bson:"request"`
	CreatedAt  time.Time `bson:"created_at"`
}

type LineBotRequestRepository struct {
	database.MongoBaseRepository
}

const mongoCollectionName = "line_bot_requests"

func NewLineBotRequestRepository() *LineBotRequestRepository {
	mbr := database.NewMongoBaseRepository(mongoCollectionName)
	mbr.Models = []LineBotRequest{}
	lbrr := &LineBotRequestRepository{
		MongoBaseRepository: *mbr,
	}
	return lbrr
}
