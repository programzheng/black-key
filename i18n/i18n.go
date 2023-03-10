package i18n

import (
	"embed"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/programzheng/black-key/config"
	"golang.org/x/text/language"
)

type Translation struct {
	Package interface{}
}

var LocaleFSPathForLoad string
var LocaleFS embed.FS

func NewI18nBundle() (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.TraditionalChinese)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, lang := range []string{"zh-Hant"} {
		_, err := bundle.LoadMessageFileFS(LocaleFS, fmt.Sprintf(LocaleFSPathForLoad, lang))
		if err != nil {
			return nil, err
		}
	}

	return bundle, nil
}

func NewI18nLocalizer(lang string) (*i18n.Localizer, error) {
	b, err := NewI18nBundle()
	if err != nil {
		return nil, err
	}
	return i18n.NewLocalizer(b, lang), nil
}

func NewTranslation(lang string) (*Translation, error) {
	if lang == "" {
		lang = config.Cfg.GetString("DEFAULT_LANGUAGE")
	}
	l, err := NewI18nLocalizer(lang)
	if err != nil {
		log.Printf("i18.go NewTranslation error: %v", err)
		return nil, err
	}
	return &Translation{
		Package: l,
	}, nil
}

func (t *Translation) Translate(key string) string {
	if t.Package == nil {
		nt, err := NewTranslation("")
		if err != nil {
			return key
		}
		t = nt
	}

	localizer := t.Package.(*i18n.Localizer)
	r := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          key,
			Description: key,
			Other:       key,
		},
	})

	return r
}
