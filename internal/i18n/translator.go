package i18n

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

type Translator struct {
	bundle             *i18n.Bundle
	cookieName         string
	cookieHandler      *http_util.CookieHandler
	preferredLanguages []string
	allowedLanguages   []language.Tag
}

type TranslatorConfig struct {
	DefaultLanguage language.Tag
	CookieName      string
}

type Message struct {
	ID   string
	Text string
}

// NewZitadelTranslator translates to all supported languages, as the ZITADEL texts are not customizable.
func NewZitadelTranslator(defaultLanguage language.Tag) *Translator {
	return newTranslator(ZITADEL, defaultLanguage, SupportedLanguages(), "")
}

func NewNotificationTranslator(defaultLanguage language.Tag, allowedLanguages []language.Tag) *Translator {
	return newTranslator(NOTIFICATION, defaultLanguage, allowedLanguages, "")
}

func NewLoginTranslator(defaultLanguage language.Tag, allowedLanguages []language.Tag, cookieName string) *Translator {
	return newTranslator(LOGIN, defaultLanguage, allowedLanguages, cookieName)
}

func newTranslator(ns Namespace, defaultLanguage language.Tag, allowedLanguages []language.Tag, cookieName string) *Translator {
	t := new(Translator)
	t.allowedLanguages = allowedLanguages
	if len(t.allowedLanguages) == 0 {
		t.allowedLanguages = SupportedLanguages()
	}
	t.bundle = newBundle(ns, defaultLanguage, t.allowedLanguages)
	t.cookieHandler = http_util.NewCookieHandler()
	t.cookieName = cookieName
	return t
}

func (t *Translator) SupportedLanguages() []language.Tag {
	return t.allowedLanguages
}

// AddMessages adds messages to the translator for the given language tag.
// If the tag is not in the allowed languages, the messages are not added.
func (t *Translator) AddMessages(tag language.Tag, messages ...Message) error {
	if len(messages) == 0 {
		return nil
	}
	var isAllowed bool
	for _, allowed := range t.allowedLanguages {
		if allowed == tag {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
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

func (t *Translator) LocalizeWithoutArgs(id string, langs ...string) string {
	return localize(t.localizer(langs...), id, map[string]interface{}{})
}

func (t *Translator) Lang(r *http.Request) language.Tag {
	matcher := language.NewMatcher(t.allowedLanguages)
	tag, _ := language.MatchStrings(matcher, t.langsFromRequest(r)...)
	return tag
}

func (t *Translator) SetLangCookie(w http.ResponseWriter, r *http.Request, lang language.Tag) {
	t.cookieHandler.SetCookie(w, t.cookieName, r.Host, lang.String())
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
		if ctxData.PreferredLanguage != language.Und.String() {
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
	acceptLanguage := metautils.ExtractIncoming(ctx).Get("accept-language")
	if acceptLanguage != "" {
		return acceptLanguage
	}
	return metautils.ExtractIncoming(ctx).Get("grpcgateway-accept-language")
}

func localize(localizer *i18n.Localizer, id string, args map[string]interface{}) string {
	s, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: args,
	})
	if err != nil {
		logging.WithFields("id", id, "args", args).WithError(err).Warnf("missing translation")
		return id
	}
	return s
}
