package eventstore

import (
	"context"
	"net/http"

	"github.com/caos/logging"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/query"

	"github.com/caos/zitadel/internal/iam/model"
)

type IAMRepository struct {
	IAMID    string
	LoginDir http.FileSystem

	IAMV2QuerySide *query.Queries
	supportedLangs []language.Tag
}

func (repo *IAMRepository) Languages(ctx context.Context) ([]language.Tag, error) {
	if len(repo.supportedLangs) == 0 {
		langs, err := i18n.SupportedLanguages(repo.LoginDir)
		logging.Log("ADMIN-tiMWs").OnError(err).Debug("unable to parse language")
		repo.supportedLangs = langs
	}
	return repo.supportedLangs, nil
}

func (repo *IAMRepository) GetIAM(ctx context.Context) (*model.IAM, error) {
	return repo.IAMV2QuerySide.IAMByID(ctx, repo.IAMID)
}
