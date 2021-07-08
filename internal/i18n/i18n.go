package i18n

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/caos/logging"
	"github.com/ghodss/yaml"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/api/authz"
	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/errors"
)

const (
	i18nPath = "/i18n"
)

type Translator struct {
	bundle             *i18n.Bundle
	cookieName         string
	cookieHandler      *http_util.CookieHandler
	preferredLanguages []string
}

type TranslatorConfig struct {
	DefaultLanguage language.Tag
	CookieName      string
}

type Message struct {
	ID   string
	Text string
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
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
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
		if err := addFileFromFileSystemToBundle(dir, bundle, file); err != nil {
			return nil, errors.ThrowNotFound(err, "I18N-ZS2AW", "cannot append file to Bundle")
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
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	bundle.MustParseMessageFileBytes(content, file.Name())
	return nil
}

func (t *Translator) AddMessages(tag language.Tag, messages ...Message) error {
	if len(messages) == 0 {
		return nil
	}
	i18nMessages := make([]*i18n.Message, len(messages))
	for i, message := range messages {
		i18nMessages[i] = &i18n.Message{
			ID:    message.ID,
			Other: message.Text,
		}
	}
	return t.bundle.AddMessages(tag, i18nMessages...)
}

func (t *Translator) LocalizeFromRequest(r *http.Request, id string, args map[string]interface{}) string {
	return localize(t.localizerFromRequest(r), id, args)
}

func (t *Translator) LocalizeFromCtx(ctx context.Context, id string, args map[string]interface{}) string {
	return localize(t.localizerFromCtx(ctx), id, args)
}

func (t *Translator) Localize(id string, args map[string]interface{}, langs ...string) string {
	return localize(t.localizer(langs...), id, args)
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

func (t *Translator) localizerFromCtx(ctx context.Context) *i18n.Localizer {
	return t.localizer(t.langsFromCtx(ctx)...)
}

func (t *Translator) localizer(langs ...string) *i18n.Localizer {
	return i18n.NewLocalizer(t.bundle, langs...)
}

func (t *Translator) langsFromRequest(r *http.Request) []string {
	langs := t.preferredLanguages
	if r != nil {
		lang, err := t.cookieHandler.GetCookieValue(r, t.cookieName)
		if err == nil {
			langs = append(langs, lang)
		}
		langs = append(langs, r.Header.Get("Accept-Language"))
	}
	return langs
}

func (t *Translator) langsFromCtx(ctx context.Context) []string {
	langs := t.preferredLanguages
	if ctx != nil {
		ctxData := authz.GetCtxData(ctx)
		if ctxData.PreferredLanguage != "" {
			langs = append(langs, authz.GetCtxData(ctx).PreferredLanguage)
		}
		langs = append(langs, getAcceptLanguageHeader(ctx))
	}
	return langs
}

func (t *Translator) SetPreferredLanguages(langs ...string) {
	t.preferredLanguages = langs
}

func getAcceptLanguageHeader(ctx context.Context) string {
	return metautils.ExtractIncoming(ctx).Get("grpcgateway-accept-language")
}

func localize(localizer *i18n.Localizer, id string, args map[string]interface{}) string {
	s, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: args,
	})
	if err != nil {
		logging.Log("I18N-MsF5sx").WithError(err).Warnf("missing translation")
		return id
	}
	return s
}
