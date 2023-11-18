package admin

import (
	"errors"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"golang.org/x/text/language"
)

func selectLanguagesToCommand(languages *admin.SelectLanguages) (tags []language.Tag, err error) {
	allowedLanguages := languages.GetList()
	if allowedLanguages == nil && languages != nil {
		allowedLanguages = make([]string, 0)
	}
	if allowedLanguages == nil {
		return nil, nil
	}
	tags = make([]language.Tag, len(allowedLanguages))
	for i, lang := range allowedLanguages {
		var parseErr error
		tags[i], parseErr = language.Parse(lang)
		err = errors.Join(err, parseErr)
	}
	return tags, err
}
