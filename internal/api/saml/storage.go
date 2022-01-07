package saml

import (
	"context"
	"github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/samlp"
	"github.com/caos/zitadel/internal/auth/repository"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type StorageConfig struct{}

type Storage interface {
	EntityStorage
	Health(context.Context) error
}

type EntityStorage interface {
	GetEntityByID(ctx context.Context, entityID string)
}

type ProviderStorage struct {
	repo    repository.Repository
	command *command.Commands
	query   *query.Queries
}

func (p *ProviderStorage) GetEntityByID(ctx context.Context, entityID string) {

}
func (p *ProviderStorage) Health(context.Context) error {
	return nil
}

func (p *ProviderStorage) CreateAuthRequest(ctx context.Context, req *samlp.AuthnRequest, relayState, issuerID string) (_ AuthRequestInt, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-sd436", "no user agent id")
	}

	authRequest := CreateAuthRequestToBusiness(ctx, req, issuerID, relayState, userAgentID)

	resp, err := p.repo.CreateAuthRequest(ctx, authRequest)

	return AuthRequestFromBusiness(resp)
}

func (p *ProviderStorage) AuthRequestByID(ctx context.Context, id string) (_ AuthRequestInt, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-D3g21", "no user agent id")
	}
	resp, err := p.repo.AuthRequestByIDCheckLoggedIn(ctx, id, userAgentID)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(resp)
}

func (p *ProviderStorage) AuthRequestByCode(ctx context.Context, code string) (_ AuthRequestInt, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	resp, err := p.repo.AuthRequestByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(resp)
}

func (p *ProviderStorage) GetAttributesFromNameID(ctx context.Context, nameID string) (_ map[string]interface{}, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	user, err := p.repo.UserByID(ctx, nameID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Email":             user.Email,
		"FirstName":         user.FirstName,
		"LastName":          user.LastName,
		"PreferredUsername": user.PreferredLoginName,
	}, nil
}
