package types

import (
	"context"
	"fmt"
	"strings"

	http_util "github.com/zitadel/zitadel/internal/api/http"

	"github.com/zitadel/zitadel/internal/api/assets"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/notification/templates"
	"github.com/zitadel/zitadel/internal/query"
)

func GetTemplateData(ctx context.Context, translator *i18n.Translator, translateArgs map[string]interface{}, href, msgType, lang string, policy *query.LabelPolicy) templates.TemplateData {
	assetsPrefix := http_util.ComposedOrigin(ctx) + assets.HandlerPrefix
	templateData := templates.TemplateData{
		URL:             href,
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
