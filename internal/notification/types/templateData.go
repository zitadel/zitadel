package types

import (
	"fmt"
	"strings"

	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/notification/templates"
	"github.com/caos/zitadel/internal/query"
)

func GetTemplateData(translator *i18n.Translator, translateArgs map[string]interface{}, apiDomain, href, msgType, lang string, policy *query.LabelPolicy) templates.TemplateData {
	templateData := templates.TemplateData{
		Href:            href,
		PrimaryColor:    templates.DefaultPrimaryColor,
		BackgroundColor: templates.DefaultBackgroundColor,
		FontColor:       templates.DefaultFontColor,
		LogoURL:         templates.DefaultLogo,
		FontURL:         templates.DefaultFont,
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
	if apiDomain == "" {
		return templateData
	}
	templateData.LogoURL = ""
	if policy.Light.LogoURL != "" {
		templateData.LogoURL = fmt.Sprintf("%s/assets/v1/%s/%s", apiDomain, policy.ID, policy.Light.LogoURL)
	}
	if policy.FontURL != "" {
		split := strings.Split(policy.FontURL, "/")
		templateData.FontFamily = split[len(split)-1] + "," + templates.DefaultFontFamily
		templateData.FontURL = fmt.Sprintf("%s/assets/v1/%s/%s", apiDomain, policy.ID, policy.FontURL)
	}
	return templateData
}
