package bot

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/programzheng/black-key/config"
	"github.com/programzheng/black-key/i18n"
	"github.com/programzheng/black-key/internal/service/rent_house"
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

	fm := linebot.NewFlexMessage(headText, &linebot.BubbleContainer{
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
	if config.Cfg.GetBool("LINE_MESSAGING_DEBUG") {
		dbs, err := fm.MarshalJSON()
		if err != nil {
			log.Printf("line_messaging_flex_message LINE_MESSAGING_DEBUG error:%v", err)
		}
		log.Printf(string(dbs))
	}
	return fm
}

func NewNewRentHousesFlexTemplate(altText string, rhs []*rent_house.RentHouse) []linebot.SendingMessage {
	if len(rhs) == 0 {
		return []linebot.SendingMessage{
			linebot.NewFlexMessage("目前無新上架租屋", &linebot.BubbleContainer{
				Type: linebot.FlexContainerTypeBubble,
				Header: &linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeBaseline,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Text: "目前無新上架租屋",
						},
					},
				},
			}),
		}
	}

	fms := []linebot.SendingMessage{}

	altTitle := fmt.Sprintf("%s新上架租屋通知", altText)
	altTitlefm := linebot.NewFlexMessage(altTitle, &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeBaseline,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Text: altTitle,
				},
			},
		},
	})
	fms = append(fms, altTitlefm)

	chunkSize := 12
	for i := 0; i < len(rhs); i += chunkSize {
		end := i + chunkSize
		if end > len(rhs) {
			end = len(rhs)
		}
		rhsCarousel := &linebot.CarouselContainer{}
		bodyTitleFlex := 2
		bodyContentFlex := 5
		for _, rh := range rhs[i:end] {
			bubble := &linebot.BubbleContainer{
				Type: linebot.FlexContainerTypeBubble,
				Header: &linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeBaseline,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Text:  rh.Title,
									Align: linebot.FlexComponentAlignTypeCenter,
								},
							},
						},
					},
				},
				Body: &linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeVertical,
					Spacing: linebot.FlexComponentSpacingTypeMd,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeHorizontal,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Text:  "價格",
									Flex:  &bodyTitleFlex,
									Align: linebot.FlexComponentAlignTypeCenter,
									Color: "#666666",
								},
								&linebot.TextComponent{
									Text:  fmt.Sprintf("%d", rh.Price),
									Flex:  &bodyContentFlex,
									Align: linebot.FlexComponentAlignTypeCenter,
									Color: "#666666",
								},
							},
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeHorizontal,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Text:  "單位",
									Flex:  &bodyTitleFlex,
									Align: linebot.FlexComponentAlignTypeCenter,
									Color: "#666666",
								},
								&linebot.TextComponent{
									Text:  rh.PriceUnit,
									Flex:  &bodyContentFlex,
									Align: linebot.FlexComponentAlignTypeCenter,
									Color: "#666666",
								},
							},
						},
					},
				},
				Footer: &linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeVertical,
							Contents: []linebot.FlexComponent{
								&linebot.ButtonComponent{
									Type: linebot.FlexComponentTypeButton,
									Action: linebot.NewURIAction(
										"開啟原始網頁連結",
										rh.DetailUrl,
									),
								},
							},
						},
					},
				},
			}
			if len(rh.PhotoList) > 0 {
				bubble.Hero = &linebot.ImageComponent{
					URL:         rh.PhotoList[0],
					Size:        linebot.FlexImageSizeTypeFull,
					AspectRatio: linebot.FlexImageAspectRatioType16to9,
					AspectMode:  linebot.FlexImageAspectModeTypeCover,
				}
			}
			rhsCarousel.Contents = append(rhsCarousel.Contents, bubble)
		}
		rhsfm := linebot.NewFlexMessage("新上架租屋通知", rhsCarousel)
		fms = append(fms, rhsfm)
	}

	return fms
}
