package templates

import (
	"fmt"
	"html"
	"strings"

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
	//defaultOrgName         = "CAOS AG"
	//defaultOrgURL          = "http://www.caos.ch"
	//defaultFooter1         = "Teufener Strasse 19"
	//defaultFooter2         = "CH-9000 St. Gallen"
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
	IncludeLogo     bool
	LogoURL         string
	FontURL         string
	FontFamily      string

	IncludeFooter bool
	FooterText    string
}

func (data *TemplateData) Translate(i18n *i18n.Translator, args map[string]interface{}, langs ...string) {
	data.Title = i18n.Localize(data.Title, nil, langs...)
	data.PreHeader = i18n.Localize(data.PreHeader, nil, langs...)
	data.Subject = i18n.Localize(data.Subject, nil, langs...)
	data.Greeting = i18n.Localize(data.Greeting, args, langs...)
	data.Text = html.UnescapeString(i18n.Localize(data.Text, args, langs...))
	if data.Href != "" {
		data.Href = i18n.Localize(data.Href, nil, langs...)
	}
	data.ButtonText = i18n.Localize(data.ButtonText, nil, langs...)
}

func GetTemplateData(apiDomain, href string, text *iam_model.MailTextView, policy *iam_model.LabelPolicyView) TemplateData {
	templateData := TemplateData{
		Title:           text.Title,
		PreHeader:       text.PreHeader,
		Subject:         text.Subject,
		Greeting:        text.Greeting,
		Text:            html.UnescapeString(text.Text),
		Href:            href,
		ButtonText:      text.ButtonText,
		PrimaryColor:    defaultPrimaryColor,
		BackgroundColor: defaultBackgroundColor,
		FontColor:       defaultFontColor,
		LogoURL:         defaultLogo,
		FontURL:         defaultFont,
		FontFamily:      defaultFontFamily,
		IncludeFooter:   false,
	}
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
	if policy.LogoURL == "" {
		templateData.IncludeLogo = false
	} else {
		templateData.IncludeLogo = true
		templateData.LogoURL = fmt.Sprintf("%s/assets/v1/%s", apiDomain, policy.LogoURL)
	}
	if policy.FontURL != "" {
		split := strings.Split(policy.FontURL, "/")
		templateData.FontFamily = split[len(split)-1]
		templateData.FontURL = fmt.Sprintf("%s/assets/v1/%s", apiDomain, policy.FontURL)
	}
	return templateData
}
