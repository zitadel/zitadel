package i18n

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/zitadel/logging"
	"golang.org/x/text/language"
	"sigs.k8s.io/yaml"

	"github.com/zitadel/zitadel/internal/domain"
)

const i18nPath = "/i18n"

var translationMessages = map[Namespace]map[language.Tag]*i18n.MessageFile{
	ZITADEL:      make(map[language.Tag]*i18n.MessageFile),
	LOGIN:        make(map[language.Tag]*i18n.MessageFile),
	NOTIFICATION: make(map[language.Tag]*i18n.MessageFile),
}

func init() {
	unmarshaler := map[string]i18n.UnmarshalFunc{
		"yaml": func(data []byte, v interface{}) error { return yaml.Unmarshal(data, v) },
		"json": json.Unmarshal,
		"toml": toml.Unmarshal,
	}
	for ns := range translationMessages {
		dir := LoadFilesystem(ns)
		i18nDir, err := dir.Open(i18nPath)
		logging.WithFields("namespace", ns).OnError(err).Panic("unable to open translation files")
		defer i18nDir.Close()
		files, err := i18nDir.Readdir(0)
		logging.WithFields("namespace", ns).OnError(err).Panic("unable to read translation files")
		for _, file := range files {
			f, err := dir.Open("/i18n/" + file.Name())
			logging.WithFields("namespace", ns, "file", file.Name()).OnError(err).Panic("unable to open translation file")
			defer f.Close()

			content, err := io.ReadAll(f)
			logging.WithFields("namespace", ns, "file", file.Name()).OnError(err).Panic("unable to read translation file")

			messageFile, err := i18n.ParseMessageFileBytes(content, file.Name(), unmarshaler)
			logging.WithFields("namespace", ns, "file", file.Name()).OnError(err).Panic("unable to parse translation file")

			fileLang, _ := strings.CutSuffix(file.Name(), filepath.Ext(file.Name()))
			lang := language.Make(fileLang)

			translationMessages[ns][lang] = messageFile
		}
	}
}

func newBundle(ns Namespace, defaultLanguage language.Tag, allowedLanguages []language.Tag) (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(defaultLanguage)

	for lang, file := range translationMessages[ns] {
		if err := domain.LanguageIsAllowed(false, allowedLanguages, lang); err != nil {
			continue
		}
		bundle.MustAddMessages(lang, file.Messages...)
	}

	return bundle, nil
}
