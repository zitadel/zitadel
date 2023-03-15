package saml

import (
	"context"
	"time"

	"github.com/zitadel/saml/pkg/provider"
	"github.com/zitadel/saml/pkg/provider/key"
	"github.com/zitadel/saml/pkg/provider/models"
	"github.com/zitadel/saml/pkg/provider/serviceprovider"
	"github.com/zitadel/saml/pkg/provider/xml/samlp"

	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/auth/repository"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var _ provider.EntityStorage = &Storage{}
var _ provider.IdentityProviderStorage = &Storage{}
var _ provider.AuthStorage = &Storage{}
var _ provider.UserStorage = &Storage{}

type Storage struct {
	certChan                   <-chan interface{}
	defaultCertificateLifetime time.Duration

	currentCACertificate       query.Certificate
	currentMetadataCertificate query.Certificate
	currentResponseCertificate query.Certificate

	locker               crdb.Locker
	certificateAlgorithm string
	encAlg               crypto.EncryptionAlgorithm
	certEncAlg           crypto.EncryptionAlgorithm

	eventstore *eventstore.Eventstore
	repo       repository.Repository
	command    *command.Commands
	query      *query.Queries

	defaultLoginURL string
}

func (p *Storage) GetEntityByID(ctx context.Context, entityID string) (*serviceprovider.ServiceProvider, error) {
	app, err := p.query.AppBySAMLEntityID(ctx, entityID, false)
	if err != nil {
		return nil, err
	}
	return serviceprovider.NewServiceProvider(
		app.ID,
		&serviceprovider.Config{
			Metadata: app.SAMLConfig.Metadata,
		},
		p.defaultLoginURL,
	)
}

func (p *Storage) GetEntityIDByAppID(ctx context.Context, appID string) (string, error) {
	app, err := p.query.AppByID(ctx, appID, false)
	if err != nil {
		return "", err
	}
	return app.SAMLConfig.EntityID, nil
}

func (p *Storage) Health(context.Context) error {
	return nil
}

func (p *Storage) GetCA(ctx context.Context) (*key.CertificateAndKey, error) {
	return p.GetCertificateAndKey(ctx, domain.KeyUsageSAMLCA)
}

func (p *Storage) GetMetadataSigningKey(ctx context.Context) (*key.CertificateAndKey, error) {
	return p.GetCertificateAndKey(ctx, domain.KeyUsageSAMLMetadataSigning)
}

func (p *Storage) GetResponseSigningKey(ctx context.Context) (*key.CertificateAndKey, error) {
	return p.GetCertificateAndKey(ctx, domain.KeyUsageSAMLResponseSinging)
}

func (p *Storage) CreateAuthRequest(ctx context.Context, req *samlp.AuthnRequestType, acsUrl, protocolBinding, relayState, applicationID string) (_ models.AuthRequestInt, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, errors.ThrowPreconditionFailed(nil, "SAML-sd436", "no user agent id")
	}

	authRequest := CreateAuthRequestToBusiness(ctx, req, acsUrl, protocolBinding, applicationID, relayState, userAgentID)

	resp, err := p.repo.CreateAuthRequest(ctx, authRequest)
	if err != nil {
		return nil, err
	}

	return AuthRequestFromBusiness(resp)
}

func (p *Storage) AuthRequestByID(ctx context.Context, id string) (_ models.AuthRequestInt, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, errors.ThrowPreconditionFailed(nil, "SAML-D3g21", "no user agent id")
	}
	resp, err := p.repo.AuthRequestByIDCheckLoggedIn(ctx, id, userAgentID)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(resp)
}

func (p *Storage) SetUserinfoWithUserID(ctx context.Context, userinfo models.AttributeSetter, userID string, attributes []int) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	user, err := p.query.GetUserByID(ctx, true, userID, false)
	if err != nil {
		return err
	}

	setUserinfo(user, userinfo, attributes)
	return nil
}

func (p *Storage) SetUserinfoWithLoginName(ctx context.Context, userinfo models.AttributeSetter, loginName string, attributes []int) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	loginNameSQ, err := query.NewUserLoginNamesSearchQuery(loginName)
	if err != nil {
		return err
	}
	user, err := p.query.GetUser(ctx, true, false, loginNameSQ)
	if err != nil {
		return err
	}

	setUserinfo(user, userinfo, attributes)
	return nil
}

func setUserinfo(user *query.User, userinfo models.AttributeSetter, attributes []int) {
	if len(attributes) == 0 {
		userinfo.SetUsername(user.PreferredLoginName)
		userinfo.SetUserID(user.ID)
		if user.Human == nil {
			return
		}
		userinfo.SetEmail(string(user.Human.Email))
		userinfo.SetSurname(user.Human.LastName)
		userinfo.SetGivenName(user.Human.FirstName)
		userinfo.SetFullName(user.Human.DisplayName)
		return
	}
	for _, attribute := range attributes {
		switch attribute {
		case provider.AttributeEmail:
			if user.Human != nil {
				userinfo.SetEmail(string(user.Human.Email))
			}
		case provider.AttributeSurname:
			if user.Human != nil {
				userinfo.SetSurname(user.Human.LastName)
			}
		case provider.AttributeFullName:
			if user.Human != nil {
				userinfo.SetFullName(user.Human.DisplayName)
			}
		case provider.AttributeGivenName:
			if user.Human != nil {
				userinfo.SetGivenName(user.Human.FirstName)
			}
		case provider.AttributeUsername:
			userinfo.SetUsername(user.PreferredLoginName)
		case provider.AttributeUserID:
			userinfo.SetUserID(user.ID)
		}
	}
}
