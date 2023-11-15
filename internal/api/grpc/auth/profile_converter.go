package auth

import (
	"context"

	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/v2/internal/api/grpc/user"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/pkg/grpc/auth"
)

func UpdateProfileToDomain(ctx context.Context, profile *auth.UpdateMyProfileRequest) *domain.Profile {
	lang, err := language.Parse(profile.PreferredLanguage)
	logging.Log("AUTH-x19v6").OnError(err).Debug("unable to parse preferred language")

	return &domain.Profile{
		ObjectRoot:        ctxToObjectRoot(ctx),
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		NickName:          profile.NickName,
		DisplayName:       profile.DisplayName,
		PreferredLanguage: lang,
		Gender:            user.GenderToDomain(profile.Gender),
	}
}
