package templates

import (
	"github.com/caos/zitadel/internal/i18n"
)

type TemplateData struct {
	Title      string
	PreHeader  string
	Subject    string
	Greeting   string
	Text       string
	Href       string
	ButtonText string
}

func (data *TemplateData) Translate(i18n *i18n.Translator, args map[string]interface{}) {
	data.Title = i18n.Localize(data.Title, nil)
	data.PreHeader = i18n.Localize(data.PreHeader, nil)
	data.Subject = i18n.Localize(data.Subject, nil)
	data.Greeting = i18n.Localize(data.Greeting, args)
	data.Text = i18n.Localize(data.Text, args)
	data.Href = i18n.Localize(data.Href, nil)
	data.ButtonText = i18n.Localize(data.ButtonText, nil)
}
