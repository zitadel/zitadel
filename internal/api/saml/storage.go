package saml

import (
	"context"
	"github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/api/saml/key"
	"github.com/caos/zitadel/internal/api/saml/xml"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/samlp"
	"github.com/caos/zitadel/internal/auth/repository"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"time"
)

type StorageConfig struct {
	DefaultLoginURL string
}

type Storage interface {
	EntityStorage
	Health(context.Context) error
}

type EntityStorage interface {
	GetEntityByID(ctx context.Context, entityID string) (*ServiceProvider, error)
	GetEntityIDByAppID(ctx context.Context, entityID string) (string, error)
	GetCA(context.Context, chan<- key.CertificateAndKey)
	GetMetadataSigningKey(context.Context, chan<- key.CertificateAndKey)
	GetResponseSigningKey(context.Context, chan<- key.CertificateAndKey)
}

type UserStorage interface {
	SetUserinfo(ctx context.Context, userinfo AttributeSetter, userID, applicationID string, attributes []int) (err error)
}

type ProviderStorage struct {
	certChan                  <-chan interface{}
	certificateRotationCheck  time.Duration
	certificateGracefulPeriod time.Duration

	currentCACertificate       query.Certificate
	currentMetadataCertificate query.Certificate
	currentResponseCertificate query.Certificate

	locker               crdb.Locker
	certificateAlgorithm string
	encAlg               crypto.EncryptionAlgorithm

	eventstore *eventstore.Eventstore
	repo       repository.Repository
	command    *command.Commands
	query      *query.Queries

	SignAlgorithm   string
	defaultLoginURL string
}

func (p *ProviderStorage) GetEntityByID(ctx context.Context, entityID string) (*ServiceProvider, error) {
	app, err := p.query.AppBySAMLEntityID(ctx, entityID)
	if err != nil {
		return nil, err
	}
	metadata := app.SAMLConfig.Metadata

	return NewServiceProvider(
		app.ID,
		&ServiceProviderConfig{
			Metadata: metadata,
			URL:      app.SAMLConfig.MetadataURL,
		},
		p.defaultLoginURL,
	)
}

func (p *ProviderStorage) GetEntityIDByAppID(ctx context.Context, appID string) (string, error) {
	app, err := p.query.AppByID(ctx, appID)
	if err != nil {
		return "", err
	}
	metadata, err := xml.ParseMetadataXmlIntoStruct([]byte(app.SAMLConfig.Metadata))
	if err != nil {
		return "", err
	}
	return string(metadata.EntityID), nil
}

func (p *ProviderStorage) Health(context.Context) error {
	return nil
}

func (p *ProviderStorage) GetCA(ctx context.Context, certAndKeyChan chan<- key.CertificateAndKey) {
	p.GetCertificateAndKey(ctx, certAndKeyChan, key_model.KeyUsageSAMLCA)
}

func (p *ProviderStorage) GetMetadataSigningKey(ctx context.Context, certAndKeyChan chan<- key.CertificateAndKey) {
	p.GetCertificateAndKey(ctx, certAndKeyChan, key_model.KeyUsageSAMLMetadataSigning)
}

func (p *ProviderStorage) GetResponseSigningKey(ctx context.Context, certAndKeyChan chan<- key.CertificateAndKey) {
	p.GetCertificateAndKey(ctx, certAndKeyChan, key_model.KeyUsageSAMLResponseSinging)
}

func (p *ProviderStorage) CreateAuthRequest(ctx context.Context, req *samlp.AuthnRequest, acsUrl, protocolBinding, relayState, issuerID string) (_ AuthRequestInt, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-sd436", "no user agent id")
	}

	authRequest := CreateAuthRequestToBusiness(ctx, req, acsUrl, protocolBinding, issuerID, relayState, userAgentID)

	resp, err := p.repo.CreateAuthRequest(ctx, authRequest)
	if err != nil {
		return nil, err
	}

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

	user, err := p.query.GetUserByID(ctx, nameID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Email":             user.Human.Email,
		"FirstName":         user.Human.FirstName,
		"LastName":          user.Human.LastName,
		"PreferredUsername": user.PreferredLoginName,
	}, nil
}

func (p *ProviderStorage) SetUserinfo(ctx context.Context, userinfo AttributeSetter, userID, applicationID string, attributes []int) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	user, err := p.query.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	for _, attribute := range attributes {
		switch attribute {
		case AttributeEmail:
			userinfo.SetEmail(user.Human.Email)
		case AttributeSurname:
			userinfo.SetSurname(user.Human.LastName)
		case AttributeFullName:
			userinfo.SetFullName(user.Human.DisplayName)
		case AttributeGivenName:
			userinfo.SetGivenName(user.Human.FirstName)
		case AttributeUsername:
			userinfo.SetUsername(user.PreferredLoginName)
		case AttributeUserID:
			userinfo.SetUserID(userID)
		case AttributeApplicationID:
			userinfo.SetApplicationID(applicationID)
		}
	}
	if attributes == nil || len(attributes) == 0 {
		userinfo.SetEmail(user.Human.Email)
		userinfo.SetSurname(user.Human.LastName)
		userinfo.SetGivenName(user.Human.FirstName)
		userinfo.SetFullName(user.Human.DisplayName)
		userinfo.SetUsername(user.PreferredLoginName)
		userinfo.SetUserID(userID)
	}
	return nil
}
