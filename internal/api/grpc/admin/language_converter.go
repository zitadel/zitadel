package admin

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/pkg/grpc/admin"
)

func selectLanguagesToCommand(languages *admin.SelectLanguages) (tags []language.Tag, err error) {
	allowedLanguages := languages.GetList()
	if allowedLanguages == nil && languages != nil {
		allowedLanguages = make([]string, 0)
	}
	if allowedLanguages == nil {
		return nil, nil
	}
	return domain.ParseLanguage(allowedLanguages...)
}
