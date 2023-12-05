package i18n

import (
	"errors"
	"strings"

	"golang.org/x/text/language"
)

var supportedLanguages []language.Tag

func SupportedLanguages() []language.Tag {
	if supportedLanguages == nil {
		panic("supported languages not loaded")
	}
	return supportedLanguages
}

func SupportLanguages(languages ...language.Tag) {
	supportedLanguages = languages
}

func MustLoadSupportedLanguagesFromDir() {
	var err error
	defer func() {
		if err != nil {
			panic("failed to load supported languages: " + err.Error())
		}
	}()
	if supportedLanguages != nil {
		return
	}
	i18nDir, err := LoadFilesystem(LOGIN).Open(i18nPath)
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, i18nDir.Close())
	}()
	files, err := i18nDir.Readdir(0)
	if err != nil {
		return
	}
	supportedLanguages = make([]language.Tag, 0, len(files))
	for _, file := range files {
		lang := language.Make(strings.TrimSuffix(file.Name(), ".yaml"))
		if lang != language.Und {
			supportedLanguages = append(supportedLanguages, lang)
		}
	}
}
