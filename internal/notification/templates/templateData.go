package templates

import (
	"fmt"
	"html"
	"strings"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/i18n"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

const (
	defaultFont            = "http://fonts.googleapis.com/css?family=Lato:200,300,400,600"
	defaultFontFamily      = "-apple-system, BlinkMacSystemFont, Segoe UI, Lato, Arial, Helvetica, sans-serif"
	defaultLogo            = "https://static.zitadel.ch/zitadel-logo-dark@3x.png"
	defaultFontColor       = "#22292f"
	defaultBackgroundColor = "#fafafa"
	defaultPrimaryColor    = "#5282C1"
)

type TemplateData struct {
	Title           string
	PreHeader       string
	Subject         string
	Greeting        string
	Text            string
	Href            string
	ButtonText      string
	PrimaryColor    string
	BackgroundColor string
	FontColor       string
	LogoURL         string
	FontURL         string
	FontFamily      string

	IncludeFooter bool
	FooterText    string
}

func (data *TemplateData) Translate(translator *i18n.Translator, msgType string, args map[string]interface{}, langs ...string) {
	data.Title = translator.Localize(fmt.Sprintf("%s.%s", msgType, domain.MessageTitle), nil, langs...)
	data.PreHeader = translator.Localize(fmt.Sprintf("%s.%s", msgType, domain.MessagePreHeader), nil, langs...)
	data.Subject = translator.Localize(fmt.Sprintf("%s.%s", msgType, domain.MessageSubject), nil, langs...)
	data.Greeting = translator.Localize(fmt.Sprintf("%s.%s", msgType, domain.MessageGreeting), args, langs...)
	data.Text = html.UnescapeString(translator.Localize(fmt.Sprintf("%s.%s", msgType, domain.MessageText), args, langs...))
	data.ButtonText = translator.Localize(fmt.Sprintf("%s.%s", msgType, domain.MessageButtonText), nil, langs...)
	data.FooterText = translator.Localize(fmt.Sprintf("%s.%s", msgType, domain.MessageFooterText), nil, langs...)
}

func GetTemplateData(translator *i18n.Translator, translateArgs map[string]interface{}, apiDomain, href, msgType, lang string, policy *iam_model.LabelPolicyView) TemplateData {
	templateData := TemplateData{
		Href:            href,
		PrimaryColor:    defaultPrimaryColor,
		BackgroundColor: defaultBackgroundColor,
		FontColor:       defaultFontColor,
		LogoURL:         defaultLogo,
		FontURL:         defaultFont,
		FontFamily:      defaultFontFamily,
		IncludeFooter:   false,
	}
	templateData.Translate(translator, msgType, translateArgs, lang)
	if policy.PrimaryColor != "" {
		templateData.PrimaryColor = policy.PrimaryColor
	}
	if policy.BackgroundColor != "" {
		templateData.BackgroundColor = policy.BackgroundColor
	}
	if policy.FontColor != "" {
		templateData.FontColor = policy.FontColor
	}
	if apiDomain == "" {
		return templateData
	}
	templateData.LogoURL = ""
	if policy.LogoURL != "" {
		templateData.LogoURL = fmt.Sprintf("%s/assets/v1/%s/%s", apiDomain, policy.AggregateID, policy.LogoURL)
	}
	if policy.FontURL != "" {
		split := strings.Split(policy.FontURL, "/")
		templateData.FontFamily = split[len(split)-1] + "," + defaultFontFamily
		templateData.FontURL = fmt.Sprintf("%s/assets/v1/%s/%s", apiDomain, policy.AggregateID, policy.FontURL)
	}
	return templateData
}
