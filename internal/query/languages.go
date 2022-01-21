package query

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/i18n"
	"golang.org/x/text/language"
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
