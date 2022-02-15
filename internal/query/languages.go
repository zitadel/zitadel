package query

import (
	"context"

	"github.com/caos/logging"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/i18n"
)

func (q *Queries) Languages(ctx context.Context) ([]language.Tag, error) {
	if len(q.supportedLangs) == 0 {
		langs, err := i18n.SupportedLanguages(q.LoginDir)
		if err != nil {
			logging.Log("ADMIN-tiMWs").WithError(err).Debug("unable to parse language")
			return nil, err
		}
		q.supportedLangs = langs
	}
	return q.supportedLangs, nil
}
