package bot

import (
	"time"

	"github.com/programzheng/black-key/config"
	"github.com/programzheng/black-key/pkg/helper"
	"github.com/programzheng/black-key/pkg/library/line/bot/template"
	"github.com/programzheng/black-key/pkg/service/bot"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"
)

const lineOfficalID = "Udeadbeefdeadbeefdeadbeefdeadbeef"

var botClient = bot.SetLineBot()

func LineWebHook(ctx *gin.Context) {
	events, err := botClient.ParseRequest(ctx.Request)
	if err != nil {
		log.Println("LINE Message API parse Request error:", err)
	}

	for _, event := range events {
		request, err := event.MarshalJSON()
		if err != nil {
			log.Println("LINE Message API event to json error:", err)
		}
		if event.Source.UserID == lineOfficalID {
			helper.Success(ctx, nil, nil)
			return
		}
		requestString := string(request)
		lineBotRequest := bot.LineBotRequest{
			Type:       string(event.Source.Type),
			GroupID:    event.Source.GroupID,
			RoomID:     event.Source.RoomID,
			UserID:     event.Source.UserID,
			ReplyToken: event.ReplyToken,
			Request:    requestString,
		}
		if _, err := lineBotRequest.Add(); err != nil {
			helper.Fail(ctx, err)
			return
		}
		lineId := bot.LineID{
			GroupID: event.Source.GroupID,
			RoomID:  event.Source.RoomID,
			UserID:  event.Source.UserID,
		}
		switch event.Source.Type {
		case "user":
			switch event.Type {
			case linebot.EventTypeMessage:
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					replyTemplateMessage, err := bot.UserParseTextGenTemplate(lineId, message.Text)
					if err != nil {
						log.Printf("UserParseTextGenTemplate error: %v", err)
					}
					if replyTemplateMessage != nil {
						bot.LineReplyMessage(event.ReplyToken, replyTemplateMessage)
					}
				}
			case linebot.EventTypePostback:
			}
		case "group":
			switch event.Type {
			case linebot.EventTypeMessage:
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					replyTemplateMessage, err := bot.GroupParseTextGenTemplate(lineId, message.Text)
					if err != nil {
						log.Printf("GroupParseTextGenTemplate error: %v", err)
					}
					if replyTemplateMessage != nil {
						bot.LineReplyMessage(event.ReplyToken, replyTemplateMessage)
					}
				}
			case linebot.EventTypePostback:
				replyTemplateMessage := bot.GroupParsePostBackGenTemplate(lineId, event.Postback)
				if replyTemplateMessage != nil {
					bot.LineReplyMessage(event.ReplyToken, replyTemplateMessage)
				}
			}
		}

	}
}

func LinePush(ctx *gin.Context) {
	var pushMessage bot.LineBotPushMessage

	if err := ctx.BindJSON(&pushMessage); err != nil {
		helper.BadRequest(ctx, err)
		return
	}
	token := helper.CreateMD5(time.Now().Format(helper.Iso8601))
	if pushMessage.Token != token {
		helper.Unauthorized(ctx, nil)
		return
	}
	pushId := config.Cfg.GetString("LINE_DEFAULT_PUSH_ID")

	if pushMessage.PushID != "" {
		pushId = pushMessage.PushID
	}

	err := bot.LinePushMessage(pushId, template.Text(pushMessage.Text))
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
}

func defaultTemplateMessage() *linebot.TemplateMessage {
	leftBtn := linebot.NewMessageAction("left", "left clicked")
	rightBtn := linebot.NewMessageAction("right", "right clicked")
	template := linebot.NewConfirmTemplate("Hello World", leftBtn, rightBtn)
	message := linebot.NewTemplateMessage("Reply", template)
	return message
}
