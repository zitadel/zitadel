package i18n

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"sigs.k8s.io/yaml"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const i18nPath = "/i18n"

func newBundle(dir http.FileSystem, defaultLanguage language.Tag, allowedLanguages []language.Tag) (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(defaultLanguage)
	bundle.RegisterUnmarshalFunc("yaml", func(data []byte, v interface{}) error { return yaml.Unmarshal(data, v) })
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	i18nDir, err := dir.Open(i18nPath)
	if err != nil {
		return nil, zerrors.ThrowNotFound(err, "I18N-MnXRie", "path not found")
	}
	defer i18nDir.Close()
	files, err := i18nDir.Readdir(0)
	if err != nil {
		return nil, zerrors.ThrowNotFound(err, "I18N-Gew23", "cannot read dir")
	}
	for _, file := range files {
		fileLang, _ := strings.CutSuffix(file.Name(), filepath.Ext(file.Name()))
		if err = domain.LanguageIsAllowed(false, allowedLanguages, language.Make(fileLang)); err != nil {
			continue
		}
		if err := addFileFromFileSystemToBundle(dir, bundle, file); err != nil {
			return nil, zerrors.ThrowNotFoundf(err, "I18N-ZS2AW", "cannot append file %s to Bundle", file.Name())
		}
	}
	return bundle, nil
}

func addFileFromFileSystemToBundle(dir http.FileSystem, bundle *i18n.Bundle, file os.FileInfo) error {
	f, err := dir.Open("/i18n/" + file.Name())
	if err != nil {
		return err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	_, err = bundle.ParseMessageFileBytes(content, file.Name())
	return err
}
