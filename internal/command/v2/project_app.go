package command

import (
	"context"
	"strings"
	"time"

	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	project_repo "github.com/caos/zitadel/internal/repository/project"
)

type AddApp struct {
}

type addOIDCApp struct {
	App     AddApp
	Version domain.OIDCVersion
}

func AddOIDCApp(
	a project_repo.Aggregate,
	version domain.OIDCVersion,
	appID,
	name string,
	redirectUris []string,
	responseTypes []domain.OIDCResponseType,
	grantTypes []domain.OIDCGrantType,
	applicationType domain.OIDCApplicationType,
	authMethodType domain.OIDCAuthMethodType,
	postLogoutRedirectUris []string,
	devMode bool,
	accessTokenType domain.OIDCTokenType,
	accessTokenRoleAssertion bool,
	idTokenRoleAssertion bool,
	idTokenUserinfoAssertion bool,
	clockSkew time.Duration,
	additionalOrigins []string,
	clientSecretAlg crypto.HashAlgorithm,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if appID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-NnavI", "Errors.Invalid.Argument")
		}

		if name = strings.TrimSpace(name); name == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-Fef31", "Errors.Invalid.Argument")
		}

		if clockSkew > time.Second*5 || clockSkew < 0 {
			return nil, errors.ThrowInvalidArgument(nil, "V2-PnCMS", "Errors.Invalid.Argument")
		}

		for _, origin := range additionalOrigins {
			if !http_util.IsOrigin(origin) {
				return nil, errors.ThrowInvalidArgument(nil, "V2-DqWPX", "Errors.Invalid.Argument")
			}
		}

		if !domain.ContainsRequiredGrantTypes(responseTypes, grantTypes) {
			return nil, errors.ThrowInvalidArgument(nil, "V2-sLpW1", "Errors.Invalid.Argument")
		}

		return func(ctx context.Context, filter preparation.FilterToQueryReducer) (_ []eventstore.Command, err error) {
			project, err := projectWriteModel(ctx, filter, a.ID, a.ResourceOwner)
			if err != nil || !project.State.Valid() {
				return nil, errors.ThrowNotFound(err, "PROJE-6swVG", "Errors.Project.NotFound")
			}

			clientID, err := domain.NewClientID(id.SonyFlakeGenerator, project.Name)
			if err != nil {
				return nil, errors.ThrowInternal(err, "V2-VMSQ1", "Errors.Internal")
			}

			var (
				clientSecret      *crypto.CryptoValue
				clientSecretPlain string
			)
			//requires client secret
			// TODO(release blocking):we have to return the secret
			if authMethodType == domain.OIDCAuthMethodTypeBasic || authMethodType == domain.OIDCAuthMethodTypePost {
				clientSecret, clientSecretPlain, err = newAppClientSecret(ctx, filter, clientSecretAlg)
				if err != nil {
					return nil, err
				}
			}

			return []eventstore.Command{
				project_repo.NewApplicationAddedEvent(
					ctx,
					&a.Aggregate,
					appID,
					name,
				),
				project_repo.NewOIDCConfigAddedEvent(
					ctx,
					&a.Aggregate,
					version,
					appID,
					clientID,
					clientSecret,
					redirectUris,
					responseTypes,
					grantTypes,
					applicationType,
					authMethodType,
					postLogoutRedirectUris,
					devMode,
					accessTokenType,
					accessTokenRoleAssertion,
					idTokenRoleAssertion,
					idTokenUserinfoAssertion,
					clockSkew,
					additionalOrigins,
				),
			}, nil
		}, nil
	}
}

func AddAPIApp(
	a project_repo.Aggregate,
	appID,
	name string,
	authMethodType domain.APIAuthMethodType,
	clientSecretAlg crypto.HashAlgorithm,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if appID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-XHsKt", "Errors.Invalid.Argument")
		}
		if name = strings.TrimSpace(name); name == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-F7g21", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			project, err := projectWriteModel(ctx, filter, a.ID, a.ResourceOwner)
			if err != nil || !project.State.Valid() {
				return nil, errors.ThrowNotFound(err, "PROJE-Sf2gb", "Errors.Project.NotFound")
			}

			clientID, err := domain.NewClientID(id.SonyFlakeGenerator, project.Name)
			if err != nil {
				return nil, errors.ThrowInternal(err, "V2-f0pgP", "Errors.Internal")
			}

			var (
				clientSecret      *crypto.CryptoValue
				clientSecretPlain string
			)
			//requires client secret
			// TODO(release blocking):we have to return the secret
			if authMethodType == domain.APIAuthMethodTypeBasic {
				clientSecret, clientSecretPlain, err = newAppClientSecret(ctx, filter, clientSecretAlg)
				if err != nil {
					return nil, err
				}
			}

			return []eventstore.Command{
				project_repo.NewApplicationAddedEvent(
					ctx,
					&a.Aggregate,
					appID,
					name,
				),
				project_repo.NewAPIConfigAddedEvent(
					ctx,
					&a.Aggregate,
					appID,
					clientID,
					clientSecret,
					authMethodType,
				),
			}, nil
		}, nil
	}
}

func newAppClientSecret(ctx context.Context, filter preparation.FilterToQueryReducer, alg crypto.HashAlgorithm) (value *crypto.CryptoValue, plain string, err error) {
	return newCryptoCodeWithPlain(ctx, filter, domain.SecretGeneratorTypeAppSecret, alg)
}
