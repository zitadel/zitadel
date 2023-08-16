package templates

import (
	"fmt"
	"html"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/i18n"
)

const (
	DefaultFontFamily      = "-apple-system, BlinkMacSystemFont, Segoe UI, Lato, Arial, Helvetica, sans-serif"
	DefaultFontColor       = "#22292f"
	DefaultBackgroundColor = "#fafafa"
	DefaultPrimaryColor    = "#5282C1"
)

type TemplateData struct {
	Title           string
	PreHeader       string
	Subject         string
	Greeting        string
	Text            string
	URL             string
	ButtonText      string
	PrimaryColor    string
	BackgroundColor string
	FontColor       string
	LogoURL         string
	FontURL         string
	FontFaceFamily  string
	FontFamily      string

	IncludeFooter bool
	FooterText    string
}

func (data *TemplateData) Translate(translator *i18n.Translator, msgType string, args map[string]interface{}, langs ...string) {
	data.Title = translator.Localize(fmt.Sprintf("%s.%s", msgType, domain.MessageTitle), args, langs...)
	data.PreHeader = translator.Localize(fmt.Sprintf("%s.%s", msgType, domain.MessagePreHeader), args, langs...)
	data.Subject = translator.Localize(fmt.Sprintf("%s.%s", msgType, domain.MessageSubject), args, langs...)
	data.Greeting = translator.Localize(fmt.Sprintf("%s.%s", msgType, domain.MessageGreeting), args, langs...)
	data.Text = html.UnescapeString(translator.Localize(fmt.Sprintf("%s.%s", msgType, domain.MessageText), args, langs...))
	data.ButtonText = translator.Localize(fmt.Sprintf("%s.%s", msgType, domain.MessageButtonText), args, langs...)
	// Footer text is neither included in i18n files nor defaults.yaml
	footerText := fmt.Sprintf("%s.%s", msgType, domain.MessageFooterText)
	data.FooterText = translator.Localize(footerText, args, langs...)
	// translator returns the id of the string to be translated if no translation is found for that id
	// we'll include the footer if we have a custom non-empty string and if the string doesn't include the
	// id of the string that could not be translated example InitCode.Footer
	data.IncludeFooter = len(data.FooterText) > 0 && data.FooterText != footerText
}
