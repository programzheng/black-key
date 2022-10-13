package i18n

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func TestNewI18nBundle(t *testing.T) {
	_, err := newI18nBundle()
	if err != nil {
		t.Fatalf("TestNewI18nBundle %v", err)
		return
	}
	t.Log("TestNewI18nBundle success")
}

func TestNewZhHantI18nLocalizer(t *testing.T) {
	_, err := NewI18nLocalizer(language.TraditionalChinese.String())
	if err != nil {
		t.Fatalf("TestNewZhHantI18nLocalizer %v", err)
		return
	}
	t.Log("TestNewZhHantI18nLocalizer success")
}

func TestNewI18nLocalizerAndTranslate(t *testing.T) {
	localizer, err := NewI18nLocalizer(language.TraditionalChinese.String())
	if err != nil {
		t.Fatalf("TestNewZhHantI18nLocalizer %v", err)
		return
	}
	test := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "Test",
			Description: "test",
			Other:       "test",
		},
	})
	t.Logf("TestNewZhHantI18nLocalizer %s", test)
}

func TestTranslate(t *testing.T) {
	test := (&Translation{}).Translate("Test")
	t.Logf("TestTranslate %s", test)
}
