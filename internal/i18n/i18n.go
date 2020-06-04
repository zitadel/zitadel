package i18n

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/errors"

	"github.com/caos/logging"
	"github.com/ghodss/yaml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	i18nPath = "/i18n"
)

type Translator struct {
	bundle        *i18n.Bundle
	cookieName    string
	cookieHandler *http_util.CookieHandler
}

type TranslatorConfig struct {
	DefaultLanguage language.Tag
	CookieName      string
}

func NewTranslator(dir http.FileSystem, config TranslatorConfig) (*Translator, error) {
	t := new(Translator)
	var err error
	t.bundle, err = newBundle(dir, config.DefaultLanguage)
	if err != nil {
		return nil, err
	}
	t.cookieHandler = http_util.NewCookieHandler()
	t.cookieName = config.CookieName
	return t, nil
}

func newBundle(dir http.FileSystem, defaultLanguage language.Tag) (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(defaultLanguage)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	i18nDir, err := dir.Open(i18nPath)
	if err != nil {
		return nil, errors.ThrowNotFound(err, "I18N-MnXRie", "path not found")
	}
	defer i18nDir.Close()
	files, err := i18nDir.Readdir(0)
	if err != nil {
		return nil, errors.ThrowNotFound(err, "I18N-Gew23", "cannot read dir")
	}
	for _, file := range files {
		if err := addFileToBundle(dir, bundle, file); err != nil {
			return nil, errors.ThrowNotFound(err, "I18N-ZS2AW", "cannot append file to bundle")
		}
	}
	return bundle, nil
}

func addFileToBundle(dir http.FileSystem, bundle *i18n.Bundle, file os.FileInfo) error {
	f, err := dir.Open("/i18n/" + file.Name())
	if err != nil {
		return err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	bundle.MustParseMessageFileBytes(content, file.Name())
	return nil
}

func (t *Translator) LocalizeFromRequest(r *http.Request, id string, args map[string]interface{}) string {
	s, err := t.localizerFromRequest(r).Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: args,
	})
	if err != nil {
		logging.Log("I18N-MsF5sx").WithError(err).Warnf("missing translation")
		return id
	}
	return s
}

func (t *Translator) Localize(id string, args map[string]interface{}) string {
	s, _ := t.localizer().Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: args,
	})
	return s
}

func (t *Translator) Lang(r *http.Request) language.Tag {
	matcher := language.NewMatcher(t.bundle.LanguageTags())
	tag, _ := language.MatchStrings(matcher, t.langsFromRequest(r)...)
	return tag
}

func (t *Translator) SetLangCookie(w http.ResponseWriter, lang language.Tag) {
	t.cookieHandler.SetCookie(w, t.cookieName, lang.String())
}

func (t *Translator) localizerFromRequest(r *http.Request) *i18n.Localizer {
	return t.localizer(t.langsFromRequest(r)...)
}

func (t *Translator) localizer(langs ...string) *i18n.Localizer {
	return i18n.NewLocalizer(t.bundle, langs...)
}

func (t *Translator) langsFromRequest(r *http.Request) []string {
	langs := make([]string, 0)
	if r != nil {
		lang, err := t.cookieHandler.GetCookieValue(r, t.cookieName)
		if err == nil {
			langs = append(langs, lang)
		}
		langs = append(langs, r.Header.Get("Accept-Language"))
	}
	return langs
}
