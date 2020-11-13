package templates

import (
	"html"

	"github.com/caos/zitadel/internal/i18n"
)

type TemplateData struct {
	Title          string
	PreHeader      string
	Subject        string
	Greeting       string
	Text           string
	Href           string
	ButtonText     string
	PrimaryColor   string
	SecondaryColor string
}

func (data *TemplateData) Translate(i18n *i18n.Translator, args map[string]interface{}, langs ...string) {
	data.Title = i18n.Localize(data.Title, nil, langs...)
	data.PreHeader = i18n.Localize(data.PreHeader, nil, langs...)
	data.Subject = i18n.Localize(data.Subject, nil, langs...)
	data.Greeting = i18n.Localize(data.Greeting, args, langs...)
	data.Text = html.UnescapeString(i18n.Localize(data.Text, args, langs...))
	data.Href = i18n.Localize(data.Href, nil, langs...)
	data.ButtonText = i18n.Localize(data.ButtonText, nil, langs...)
}
