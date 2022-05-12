package types

import (
	"fmt"
	"strings"

	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/notification/templates"
	"github.com/zitadel/zitadel/internal/query"
)

func GetTemplateData(translator *i18n.Translator, translateArgs map[string]interface{}, assetsPrefix, href, msgType, lang string, policy *query.LabelPolicy) templates.TemplateData {
	templateData := templates.TemplateData{
		Href:            href,
		PrimaryColor:    templates.DefaultPrimaryColor,
		BackgroundColor: templates.DefaultBackgroundColor,
		FontColor:       templates.DefaultFontColor,
		FontFamily:      templates.DefaultFontFamily,
		IncludeFooter:   false,
	}
	templateData.Translate(translator, msgType, translateArgs, lang)
	if policy.Light.PrimaryColor != "" {
		templateData.PrimaryColor = policy.Light.PrimaryColor
	}
	if policy.Light.BackgroundColor != "" {
		templateData.BackgroundColor = policy.Light.BackgroundColor
	}
	if policy.Light.FontColor != "" {
		templateData.FontColor = policy.Light.FontColor
	}
	if assetsPrefix == "" {
		return templateData
	}
	if policy.Light.LogoURL != "" {
		templateData.LogoURL = fmt.Sprintf("%s/%s/%s", assetsPrefix, policy.ID, policy.Light.LogoURL)
	}
	if policy.FontURL != "" {
		split := strings.Split(policy.FontURL, "/")
		templateData.FontFaceFamily = split[len(split)-1]
		templateData.FontURL = fmt.Sprintf("%s/%s/%s", assetsPrefix, policy.ID, policy.FontURL)
		templateData.FontFamily = templateData.FontFaceFamily + "," + templates.DefaultFontFamily
	}
	return templateData
}
