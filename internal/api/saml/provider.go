package saml

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/amdonov/xmlsig"
	"github.com/caos/zitadel/internal/api/saml/xml/metadata/md"
	"github.com/caos/zitadel/internal/auth/repository"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/query"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

type ProviderConfig struct {
	EntityID string

	MetadataCertificate *Certificate
	SignatureAlgorithm  string
	DigestAlgorithm     string
	EncryptionAlgorithm string

	Organisation  *Organisation
	ContactPerson *ContactPerson

	ValidUntil    string
	CacheDuration string
	ErrorURL      string

	IDP           *IdentityProviderConfig
	StorageConfig *StorageConfig
}

type Certificate struct {
	Path           string
	PrivateKeyPath string
	CaPath         string
}

type Organisation struct {
	Name        string
	DisplayName string
	URL         string
}

type ContactPerson struct {
	ContactType     md.ContactTypeType
	Company         string
	GivenName       string
	SurName         string
	EmailAddress    string
	TelephoneNumber string
}

func NewID() string {
	return fmt.Sprintf("_%s", uuid.New())
}

const (
	healthEndpoint    = "/healthz"
	readinessEndpoint = "/ready"
	metadataEndpoint  = "/metadata"
)

type Provider struct {
	storage Storage

	Metadata  *md.EntityDescriptor
	Signer    xmlsig.Signer
	TlsConfig *tls.Config

	IdentityProvider *IdentityProvider
}

func NewProvider(
	conf *ProviderConfig,
	command *command.Commands,
	query *query.Queries,
	repo repository.Repository,
) (*Provider, error) {
	tlsConfig, err := ConfigureTLS(conf.MetadataCertificate.Path, conf.MetadataCertificate.PrivateKeyPath, conf.MetadataCertificate.CaPath)
	if err != nil {
		return nil, err
	}

	signer, err := xmlsig.NewSignerWithOptions(tlsConfig.Certificates[0], xmlsig.SignerOptions{
		SignatureAlgorithm: conf.SignatureAlgorithm,
		DigestAlgorithm:    conf.DigestAlgorithm,
	})
	if err != nil {
		return nil, err
	}

	storage := &ProviderStorage{
		repo:    repo,
		command: command,
		query:   query,
	}

	idp, err := NewIdentityProvider(
		conf.EntityID,
		conf.IDP,
		storage,
	)
	if err != nil {
		return nil, err
	}

	return &Provider{
		Metadata:         conf.getMetadata(idp),
		TlsConfig:        tlsConfig,
		Signer:           signer,
		storage:          storage,
		IdentityProvider: idp,
	}, nil
}

func (p *Provider) Storage() Storage {
	return p.storage
}

func (p *Provider) Health(ctx context.Context) error {
	return p.Storage().Health(ctx)
}

func (p *Provider) Probes() []ProbesFn {
	return []ProbesFn{
		ReadyStorage(p.Storage()),
	}
}

type HttpInterceptor func(http.Handler) http.Handler

func CreateRouter(p Provider, interceptors ...HttpInterceptor) *mux.Router {
	intercept := buildInterceptor(interceptors...)
	router := mux.NewRouter()
	router.Use(handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"authorization", "content-type"}),
		handlers.AllowedOriginValidator(allowAllOrigins),
	))
	router.HandleFunc(healthEndpoint, healthHandler)
	router.HandleFunc(readinessEndpoint, readyHandler(p.Probes()))
	router.HandleFunc(metadataEndpoint, p.metadataHandle)

	if p.IdentityProvider != nil {
		for _, route := range p.IdentityProvider.GetRoutes() {
			router.Handle(route.Endpoint, intercept(route.HandleFunc))
		}
	}
	return router
}

var allowAllOrigins = func(_ string) bool {
	return true
}

func buildInterceptor(interceptors ...HttpInterceptor) func(http.HandlerFunc) http.Handler {
	return func(handlerFunc http.HandlerFunc) http.Handler {
		handler := handlerFuncToHandler(handlerFunc)
		for i := len(interceptors) - 1; i >= 0; i-- {
			handler = interceptors[i](handler)
		}
		return handler
	}
}

func handlerFuncToHandler(handlerFunc http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, r)
	})
}
