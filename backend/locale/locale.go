package locale

import (
	"embed"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed translates/*.json
var localesFS embed.FS

func GetNewLocalizer(lang string) *i18n.Localizer {
	bundle := i18n.NewBundle(language.Polish)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	loadBundleMessageFile(bundle, "en")
	if lang != "en" {
		loadBundleMessageFile(bundle, lang)
	}

	return i18n.NewLocalizer(bundle, lang, "en")
}

func loadBundleMessageFile(bundle *i18n.Bundle, lang string) error {
	data, err := localesFS.ReadFile("translates/" + lang + ".json")
	if err != nil {
		return nil // ignore if file not found
	}
	_, err = bundle.ParseMessageFileBytes(data, lang+".json")
	if err != nil {
		return errors.Wrapf(err, "load message file %s", lang+".json")
	}
	return nil
}
