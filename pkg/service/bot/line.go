package bot

import (
	"fmt"
	"reflect"

	log "github.com/sirupsen/logrus"

	"github.com/programzheng/black-key/config"
	"github.com/programzheng/black-key/pkg/helper"
	"github.com/programzheng/black-key/pkg/model/bot"
	"github.com/programzheng/black-key/pkg/service/billing"

	"github.com/line/line-bot-sdk-go/linebot"
)

type LineBotRequest struct {
	Type       string
	GroupID    string
	RoomID     string
	UserID     string
	ReplyToken string
	Request    string
}

type LineID struct {
	GroupID string
	RoomID  string
	UserID  string
}

type LinePostBackAction struct {
	Action string
	Data   map[string]interface{}
}

type LineBotPushMessage struct {
	PushID string `json:"pushId"`
	Token  string `json:"token"`
	Text   string `json:"text"`
}

var botClient = SetLineBot()

func SetLineBot() *linebot.Client {
	channelSecret := config.Cfg.GetString("LINE_CHANNEL_SECRET")
	channelAccessToken := config.Cfg.GetString("LINE_CHANNEL_ACCESS_TOKEN")
	botClient, err := linebot.New(channelSecret, channelAccessToken)
	if err != nil {
		log.Println("LINE bot error:", err)
	}
	return botClient
}

func (lineBotRequest *LineBotRequest) Add() (uint, error) {
	model := bot.LineBotRequest{
		Type:       lineBotRequest.Type,
		GroupID:    lineBotRequest.GroupID,
		RoomID:     lineBotRequest.RoomID,
		UserID:     lineBotRequest.UserID,
		ReplyToken: lineBotRequest.ReplyToken,
		Request:    lineBotRequest.Request,
	}
	ID, err := model.Add()
	if err != nil {
		return 0, err
	}
	return ID, nil
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
	basicResponse, err := botClient.ReplyMessage(replyToken, sendMessages...).Do()
	if err != nil {
		log.Println("LINE Message API reply message Request error:", err)
	}
	log.Printf("LINE Message API reply message Request response:%v\n", basicResponse)
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
	response, err := botClient.PushMessage(toID, sendMessages...).Do()
	if err != nil {
		log.Println("pkg/service/bot/line LinePushMessage Request error:", err)
		return err
	}
	log.Printf("pkg/service/bot/line LinePushMessage response:%v\n", response)
	return nil
}

func billingAction(lineId LineID, amount int, title string, note string) (billing.Billing, bot.LineBilling) {
	b := billing.Billing{
		Title:  title,
		Amount: amount,
		Note:   note,
	}
	billing, err := b.Add()
	if err != nil {
		log.Fatal("billingAction Billing add error:", err)
	}
	lb := bot.LineBilling{
		BillingID: billing.ID,
		GroupID:   lineId.GroupID,
		RoomID:    lineId.RoomID,
		UserID:    lineId.UserID,
	}
	_, err = lb.Add()
	if err != nil {
		log.Fatal("billingAction LineBilling add error:", err)
	}
	return b, lb
}

func generateErrorTextMessage() linebot.Message {
	return linebot.NewTextMessage("系統錯誤，請重新再試或是通知管理員")
}
