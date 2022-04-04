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
	Aggregate project_repo.Aggregate
	ID        string
	Name      string
}

type addOIDCApp struct {
	AddApp
	Version                  domain.OIDCVersion
	RedirectUris             []string
	ResponseTypes            []domain.OIDCResponseType
	GrantTypes               []domain.OIDCGrantType
	ApplicationType          domain.OIDCApplicationType
	AuthMethodType           domain.OIDCAuthMethodType
	PostLogoutRedirectUris   []string
	DevMode                  bool
	AccessTokenType          domain.OIDCTokenType
	AccessTokenRoleAssertion bool
	IDTokenRoleAssertion     bool
	IDTokenUserinfoAssertion bool
	ClockSkew                time.Duration
	AdditionalOrigins        []string

	ClientID          string
	ClientSecret      *crypto.CryptoValue
	ClientSecretPlain string
}

//AddOIDCApp prepares the commands to add an oidc app. The ClientID will be set during the CreateCommands
func AddOIDCApp(app *addOIDCApp, clientSecretAlg crypto.HashAlgorithm) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if app.ID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-NnavI", "Errors.Invalid.Argument")
		}

		if app.Name = strings.TrimSpace(app.Name); app.Name == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-Fef31", "Errors.Invalid.Argument")
		}

		if app.ClockSkew > time.Second*5 || app.ClockSkew < 0 {
			return nil, errors.ThrowInvalidArgument(nil, "V2-PnCMS", "Errors.Invalid.Argument")
		}

		for _, origin := range app.AdditionalOrigins {
			if !http_util.IsOrigin(origin) {
				return nil, errors.ThrowInvalidArgument(nil, "V2-DqWPX", "Errors.Invalid.Argument")
			}
		}

		if !domain.ContainsRequiredGrantTypes(app.ResponseTypes, app.GrantTypes) {
			return nil, errors.ThrowInvalidArgument(nil, "V2-sLpW1", "Errors.Invalid.Argument")
		}

		return func(ctx context.Context, filter preparation.FilterToQueryReducer) (_ []eventstore.Command, err error) {
			project, err := projectWriteModel(ctx, filter, app.Aggregate.ID, app.Aggregate.ResourceOwner)
			if err != nil || !project.State.Valid() {
				return nil, errors.ThrowNotFound(err, "PROJE-6swVG", "Errors.Project.NotFound")
			}

			app.ClientID, err = domain.NewClientID(id.SonyFlakeGenerator, project.Name)
			if err != nil {
				return nil, errors.ThrowInternal(err, "V2-VMSQ1", "Errors.Internal")
			}

			//requires client secret
			// TODO(release blocking):we have to return the secret
			if app.AuthMethodType == domain.OIDCAuthMethodTypeBasic || app.AuthMethodType == domain.OIDCAuthMethodTypePost {
				app.ClientSecret, app.ClientSecretPlain, err = newAppClientSecret(ctx, filter, clientSecretAlg)
				if err != nil {
					return nil, err
				}
			}

			return []eventstore.Command{
				project_repo.NewApplicationAddedEvent(
					ctx,
					&app.Aggregate.Aggregate,
					app.ID,
					app.Name,
				),
				project_repo.NewOIDCConfigAddedEvent(
					ctx,
					&app.Aggregate.Aggregate,
					app.Version,
					app.ID,
					app.ClientID,
					app.ClientSecret,
					app.RedirectUris,
					app.ResponseTypes,
					app.GrantTypes,
					app.ApplicationType,
					app.AuthMethodType,
					app.PostLogoutRedirectUris,
					app.DevMode,
					app.AccessTokenType,
					app.AccessTokenRoleAssertion,
					app.IDTokenRoleAssertion,
					app.IDTokenUserinfoAssertion,
					app.ClockSkew,
					app.AdditionalOrigins,
				),
			}, nil
		}, nil
	}
}

type addAPIApp struct {
	AddApp
	AuthMethodType domain.APIAuthMethodType

	ClientID          string
	ClientSecret      *crypto.CryptoValue
	ClientSecretPlain string
}

func AddAPIApp(app *addAPIApp, clientSecretAlg crypto.HashAlgorithm) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if app.ID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-XHsKt", "Errors.Invalid.Argument")
		}
		if app.Name = strings.TrimSpace(app.Name); app.Name == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-F7g21", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			project, err := projectWriteModel(ctx, filter, app.Aggregate.ID, app.Aggregate.ResourceOwner)
			if err != nil || !project.State.Valid() {
				return nil, errors.ThrowNotFound(err, "PROJE-Sf2gb", "Errors.Project.NotFound")
			}

			app.ClientID, err = domain.NewClientID(id.SonyFlakeGenerator, project.Name)
			if err != nil {
				return nil, errors.ThrowInternal(err, "V2-f0pgP", "Errors.Internal")
			}

			//requires client secret
			// TODO(release blocking):we have to return the secret
			if app.AuthMethodType == domain.APIAuthMethodTypeBasic {
				app.ClientSecret, app.ClientSecretPlain, err = newAppClientSecret(ctx, filter, clientSecretAlg)
				if err != nil {
					return nil, err
				}
			}

			return []eventstore.Command{
				project_repo.NewApplicationAddedEvent(
					ctx,
					&app.Aggregate.Aggregate,
					app.ID,
					app.Name,
				),
				project_repo.NewAPIConfigAddedEvent(
					ctx,
					&app.Aggregate.Aggregate,
					app.ID,
					app.ClientID,
					app.ClientSecret,
					app.AuthMethodType,
				),
			}, nil
		}, nil
	}
}

func newAppClientSecret(ctx context.Context, filter preparation.FilterToQueryReducer, alg crypto.HashAlgorithm) (value *crypto.CryptoValue, plain string, err error) {
	return newCryptoCodeWithPlain(ctx, filter, domain.SecretGeneratorTypeAppSecret, alg)
}
