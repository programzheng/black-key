package bot

import (
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/config"
	"github.com/programzheng/black-key/internal/helper"
	"github.com/programzheng/black-key/internal/service/proxy"
)

type TextParsingStrategy interface {
	Execute(lineId LineID, text string) (interface{}, error)
}

type HelpStrategy struct{}

func (s *HelpStrategy) Execute(lineId LineID, text string) (interface{}, error) {
	return getHelp(text)
}

type InfoStrategy struct{}

func (s *InfoStrategy) Execute(lineId LineID, text string) (interface{}, error) {
	if text != "資訊" {
		return nil, nil
	}
	return getLineId(lineId)
}

type MemberLineAvatarStrategy struct{}

func (s *MemberLineAvatarStrategy) Execute(lineId LineID, text string) (interface{}, error) {
	if text != "我的大頭貼" {
		return nil, nil
	}
	return getMemberLineAvatar(lineId)
}

type GroupMemberLineAvatarStrategy struct{}

func (g *GroupMemberLineAvatarStrategy) Execute(lineId LineID, text string) (interface{}, error) {
	if text != "我的大頭貼" {
		return nil, nil
	}
	return getGroupMemberLineAvatar(lineId)
}

type BillingStrategy struct{}

func (s *BillingStrategy) Execute(lineId LineID, text string) (interface{}, error) {
	match := strings.Split(text, "|")[0]
	switch match {
	case "c list", "記帳列表":
		return getLineBillings(lineId)
	case "c", "記帳":
		return createBilling(lineId, text)
	case "記帳結算", "結帳", "結算":
		return getBills(lineId, text)
	default:
		return nil, nil
	}
}

type RockPaperScissorStrategy struct{}

func (r *RockPaperScissorStrategy) Execute(lineId LineID, text string) (interface{}, error) {
	if text != "猜拳" && text != "石頭布剪刀" && text != "剪刀石頭布" && text != "rock-paper-scissors" {
		return nil, nil
	}
	return startRockPaperScissor(lineId)
}

type TodoStrategy struct{}

func (s *TodoStrategy) Execute(lineId LineID, text string) (interface{}, error) {
	match := strings.Split(text, "|")[0]
	switch match {
	case "所有提醒", "所有通知", "All TODO":
		return getTodo(lineId)
	case "提醒", "通知", "TODO":
		return todo(lineId, text)
	case "多提醒", "多通知":
		return todos(lineId, text)
	default:
		return nil, nil
	}
}

type ProxyStrategy struct{}

func (s *ProxyStrategy) Execute(lineId LineID, text string) (interface{}, error) {
	imageUrl := proxy.GetGrpcProxyResponse(nil, text)
	if imageUrl != "" {
		return linebot.NewImageMessage(imageUrl, imageUrl), nil
	}
	return nil, nil
}

type DefaultStrategy struct{}

func (s *DefaultStrategy) Execute(lineId LineID, text string) (interface{}, error) {
	if helper.ConvertToBool(config.Cfg.GetString("LINE_MESSAGING_DEBUG")) {
		return linebot.NewTextMessage(text), nil
	}
	return nil, nil
}
