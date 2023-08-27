package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	log "github.com/sirupsen/logrus"

	"github.com/programzheng/black-key/config"
	"github.com/programzheng/black-key/internal/filesystem"
	"github.com/programzheng/black-key/internal/helper"

	"github.com/line/line-bot-sdk-go/linebot"
)

type LineID struct {
	GroupID string
	RoomID  string
	UserID  string
}

type LinePostBackAction struct {
	Action string
	Data   map[string]interface{}
	Params LinePostBackActionParams
}

type LinePostBackActionParams struct {
	Date     string `json:"date,omitempty"`
	Time     string `json:"time,omitempty"`
	Datetime string `json:"datetime,omitempty"`
}

type LineBotPushMessage struct {
	PushID string `json:"pushId"`
	Token  string `json:"token"`
	Text   string `json:"text"`
}

var BotClient = SetLineBot()

func SetLineBot() *linebot.Client {
	channelSecret := config.Cfg.GetString("LINE_CHANNEL_SECRET")
	channelAccessToken := config.Cfg.GetString("LINE_CHANNEL_ACCESS_TOKEN")
	BotClient, err := linebot.New(channelSecret, channelAccessToken)
	if err != nil {
		log.Println("LINE bot error:", err)
	}
	return BotClient
}

func LineReplyMessage(replyToken string, messages interface{}) {
	var sendMessages []linebot.SendingMessage
	rv := reflect.ValueOf(messages)
	if rv.Kind() == reflect.Slice {
		sendMessages = messages.([]linebot.SendingMessage)
	} else {
		sendMessages = append(sendMessages, messages.(linebot.SendingMessage))
	}
	if helper.ConvertToBool(config.Cfg.GetString("LINE_MESSAGING_DEBUG")) {
		fmt.Printf("LineReplyMessage:\nreplyToken:%s\nmessages: %v\n", replyToken, helper.GetJSON(messages))
	}
	basicResponse, err := BotClient.ReplyMessage(replyToken, sendMessages...).Do()
	if err != nil {
		log.Println("LINE Message API reply message Request error:", err)
	}
	log.Printf("LINE Message API reply message Request response:%v\n", basicResponse)
}

func GetMessageContent(messageId string) (*linebot.MessageContentResponse, error) {
	return BotClient.GetMessageContent(messageId).Do()
}

func LinePushMessage(toID string, messages interface{}) error {
	var sendMessages []linebot.SendingMessage
	rv := reflect.ValueOf(messages)
	if rv.Kind() == reflect.Slice {
		sendMessages = messages.([]linebot.SendingMessage)
	} else {
		sendMessages = append(sendMessages, messages.(linebot.SendingMessage))
	}
	if helper.ConvertToBool(config.Cfg.GetString("LINE_MESSAGING_DEBUG")) {
		fmt.Printf("LinePushMessage:\ntoID: %s\nmessages: %v\n", toID, helper.GetJSON(messages))
	}
	response, err := BotClient.PushMessage(toID, sendMessages...).Do()
	if err != nil {
		log.Println("pkg/service/bot/line LinePushMessage Request error:", err)
		return err
	}
	log.Printf("pkg/service/bot/line LinePushMessage response:%v\n", response)
	return nil
}

func (lineId *LineID) getHashKey() string {

	b, err := json.Marshal(lineId)
	if err != nil {
		log.Errorf(
			`LineID getHashKey error:%v,
			UserID: %s,
			GroupID: %s,
			RoomID: %s
			`,
			err,
			lineId.UserID,
			lineId.GroupID,
			lineId.RoomID)
	}
	j := string(b)

	return helper.CreateMD5(j)
}

func (lineId *LineID) getTodosCacheKey() string {
	return fmt.Sprintf("%s|%s", "TODOS", lineId.getHashKey())
}

func messageContentResponseToStaticFile(
	messageContentResponse *linebot.MessageContentResponse,
) (filesystem.FileSystem, *filesystem.StaticFile) {
	ctx := context.Background()
	tmpFileName := "tmp" + helper.GetFileExtensionByContentType(messageContentResponse.ContentType)
	fs := filesystem.Create("")
	return fs, fs.Upload(ctx, tmpFileName, messageContentResponse.Content)
}
