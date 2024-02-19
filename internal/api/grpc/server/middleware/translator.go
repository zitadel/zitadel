package middleware

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type localizers interface {
	Localizers() []Localizer
}
type Localizer interface {
	LocalizationKey() string
	SetLocalizedMessage(string)
}

func translateFields(ctx context.Context, object localizers, translator *i18n.Translator) {
	if translator == nil || object == nil {
		return
	}
	for _, field := range object.Localizers() {
		field.SetLocalizedMessage(translator.LocalizeFromCtx(ctx, field.LocalizationKey(), nil))
	}
}

func translateError(ctx context.Context, err error, translator *i18n.Translator) error {
	if translator == nil || err == nil {
		return err
	}
	caosErr := new(zerrors.ZitadelError)
	if errors.As(err, &caosErr) {
		descErr := new(zerrors.DescriptiveError)
		errors.As(err, &descErr)
		caosErr.SetMessage(translator.LocalizeFromCtx(ctx, caosErr.GetMessage(), descErr.GetVars()))
	}
	return err
}
