package bot

import (
	"testing"

	"github.com/programzheng/black-key/config"
)

func TestPushLineFeatureNotificationNewRentHousesFlexTemplate(t *testing.T) {
	ft := NewNewRentHousesFlexTemplate("測試市", []*RentHouse{
		{
			Title:     "test1",
			PostID:    1,
			Price:     5000,
			PriceUnit: "元/月",
			PhotoList: []string{config.Cfg.GetString("TEST_LINE_PICTURE_URL")},
			DetailUrl: "https://localhost",
		},
		{
			Title:     "test2",
			PostID:    2,
			Price:     10000,
			PriceUnit: "元/月",
			PhotoList: []string{config.Cfg.GetString("TEST_LINE_PICTURE_URL")},
			DetailUrl: "https://localhost",
		},
	})

	err := LinePushMessage(config.Cfg.GetString("TEST_LINE_PUSH_ID"), ft)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	t.Log("TestPushLineFeatureNotificationNewRentHousesFlexTemplate succeeded!")
}
