package getopt

import (
	"github.com/daforester/getopt-golang/getopt/errors"
	"fmt"
	"github.com/daforester/getopt-golang/getopt/localization"
	"strings"
)

var fallbackTranslator *Translator

type Translator struct {
	languageFile string
	translations map[string]string
}

func NewTranslator(language string, asFallback ...bool) (*Translator, error) {
	t := new(Translator)

	if language == "" {
		language = "en"
	}

	if !t.SetLanguage(language) {
		return nil, errors.New(fmt.Sprintf("language %s is not available", language))
	}

	if fallbackTranslator == nil && (len(asFallback) == 0 || asFallback[0] == false) {
		fb, err := NewTranslator("en", true)

		if err != nil {
			return nil, err
		}

		fallbackTranslator = fb
	}

	return t, nil
}

func (t *Translator) Translate(key string) string {
	if t.translations == nil {
		t.LoadTranslations()
	}

	v, isset := t.translations[key]

	if !isset {
		if t != fallbackTranslator {
			return fallbackTranslator.Translate(key)
		}

		return key
	}

	return v
}

func (t *Translator) SetLanguage(language string) bool {
	language = strings.ToLower(language)
	switch language {
	case "en":
		fallthrough
	case "engb":
		fallthrough
	case "en-gb":
		t.languageFile = "LangEN"
		return true
	case "fr":
		t.languageFile = "LangFR"
		return true
	case "de":
		t.languageFile = "LangDE"
		return true
	}

	return false
}

func (t *Translator) LoadTranslations() {
	switch t.languageFile {
	case "LangEN":
		t.translations = localization.LangEN
	case "LangFR":
		t.translations = localization.LangFR
	case "LangDE":
		t.translations = localization.LangDE
	}
}