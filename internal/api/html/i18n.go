package html

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/caos/logging"
	"github.com/ghodss/yaml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/api"
	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/errors"
)

type Translator struct {
	bundle        *i18n.Bundle
	cookieName    string
	cookieHandler *http_util.CookieHandler
}

type TranslatorConfig struct {
	Path            string
	DefaultLanguage language.Tag
	CookieName      string
}

func NewTranslator(config TranslatorConfig) (*Translator, error) {
	t := new(Translator)
	var err error
	t.bundle, err = newBundle(config.Path, config.DefaultLanguage)
	if err != nil {
		return nil, err
	}
	t.cookieHandler = http_util.NewCookieHandler()
	t.cookieName = config.CookieName
	return t, nil
}

func newBundle(i18nDir string, defaultLanguage language.Tag) (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(defaultLanguage)
	yamlUnmarshal := func(data []byte, v interface{}) error { return yaml.Unmarshal(data, v) }
	bundle.RegisterUnmarshalFunc("yaml", yamlUnmarshal)
	bundle.RegisterUnmarshalFunc("yml", yamlUnmarshal)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	files, err := ioutil.ReadDir(i18nDir)
	if err != nil {
		return nil, errors.ThrowNotFound(err, "HTML-MnXRie", "path not found")
	}
	for _, file := range files {
		bundle.MustLoadMessageFile(path.Join(i18nDir, file.Name()))
	}
	return bundle, nil
}

func (t *Translator) LocalizeFromRequest(r *http.Request, id string, args map[string]interface{}) string {
	s, err := t.localizerFromRequest(r).Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: args,
	})
	if err != nil {
		logging.Log("HTML-MsF5sx").WithError(err).Warnf("missing translation")
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
		langs = append(langs, r.Header.Get(api.AcceptLanguage))
	}
	return langs
}
