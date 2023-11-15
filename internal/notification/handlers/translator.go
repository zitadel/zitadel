package handlers

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/i18n"
)

func (n *NotificationQueries) GetTranslatorWithOrgTexts(ctx context.Context, orgID, textType string) (*i18n.Translator, error) {
	translator, err := i18n.NewTranslator(n.statikDir, n.GetDefaultLanguage(ctx), "")
	if err != nil {
		return nil, err
	}

	allCustomTexts, err := n.CustomTextListByTemplate(ctx, authz.GetInstance(ctx).InstanceID(), textType, false)
	if err != nil {
		return translator, nil
	}
	customTexts, err := n.CustomTextListByTemplate(ctx, orgID, textType, false)
	if err != nil {
		return translator, nil
	}
	allCustomTexts.CustomTexts = append(allCustomTexts.CustomTexts, customTexts.CustomTexts...)

	for _, text := range allCustomTexts.CustomTexts {
		msg := i18n.Message{
			ID:   text.Template + "." + text.Key,
			Text: text.Text,
		}
		err = translator.AddMessages(text.Language, msg)
		logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID(), "orgID", orgID, "messageType", textType, "messageID", msg.ID).
			OnError(err).
			Warn("could not add translation message")
	}
	return translator, nil
}
