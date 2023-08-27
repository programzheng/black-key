package bot

import (
	"encoding/json"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/i18n"
	log "github.com/sirupsen/logrus"
)

type DatetimePickerActionMode string

const (
	DatetimePickerActionModeDateTime DatetimePickerActionMode = "datetime"
)

func newTodoFlexTemplate(text string) *linebot.FlexMessage {
	t := &(i18n.Translation{})

	nt := time.Now()

	headText := t.Translate("LINE_Messaging_Todo_FlexTemplate_Header_Text")
	bodyTitleFlex := 2
	bodyContentFlex := 5

	lpba := LinePostBackAction{
		Action: "todo",
		Data: map[string]interface{}{
			"Text": text,
		},
	}
	lpbaJson, err := json.Marshal(lpba)
	if err != nil {
		log.Fatalf("newTodoFlexTemplate json failed: %v", err)
	}
	data := string(lpbaJson)

	return linebot.NewFlexMessage(headText, &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: headText,
				},
			},
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeBaseline,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  t.Translate("LINE_Messaging_Todo_FlexTemplate_Body_Text"),
							Size:  linebot.FlexTextSizeTypeSm,
							Color: "#aaaaaa",
							Flex:  &bodyTitleFlex,
						},
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  text,
							Size:  linebot.FlexTextSizeTypeSm,
							Flex:  &bodyContentFlex,
							Color: "#666666",
							Wrap:  true,
						},
					},
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeButton,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: linebot.NewDatetimePickerAction(
						t.Translate("LINE_Messaging_Todo_FlexTemplate_DatetimePickerAction_Label"),
						data,
						string(DatetimePickerActionModeDateTime),
						nt.Add(time.Hour).Format("2006-01-02T15:04"),
						"",
						nt.Format("2006-01-02T15:04"),
					),
				},
			},
		},
	})
}
