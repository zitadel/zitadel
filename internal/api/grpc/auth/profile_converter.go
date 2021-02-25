package auth

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/grpc/user"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/auth"
	"golang.org/x/text/language"
)

func UpdateProfileToDomain(ctx context.Context, profile *auth.UpdateMyProfileRequest) *domain.Profile {
	lang, err := language.Parse(profile.PreferredLanguage)
	logging.Log("AUTH-x19v6").OnError(err).Debug("unable to parse preferred language")

	return &domain.Profile{
		ObjectRoot:        ctxToObjectRoot(ctx),
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		NickName:          profile.NickName,
		PreferredLanguage: lang,
		Gender:            user.GenderToDomain(profile.Gender),
	}
}
