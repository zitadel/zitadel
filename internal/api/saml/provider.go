package saml

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/amdonov/xmlsig"
	"github.com/caos/logging"
	"github.com/caos/oidc/pkg/op"
	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/api/saml/key"
	"github.com/caos/zitadel/internal/api/saml/xml/md"
	"github.com/caos/zitadel/internal/auth/repository"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/telemetry/metrics"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/square/go-jose.v2"
	"net/http"
)

const (
	locksTable = "projections.locks"
	signingKey = "signing_key"
)

type ProviderConfig struct {
	CertificateAlgorithm string

	Metadata      *Metadata
	Organisation  *Organisation
	ContactPerson *ContactPerson

	IDP           *IdentityProviderConfig `yaml:"IDP"`
	StorageConfig *StorageConfig          `yaml:"StorageConfig"`

	UserAgentCookieConfig *middleware.UserAgentCookieConfig
	defaultLoginURL       string
}

type Metadata struct {
	Path                string
	URL                 string
	SignatureAlgorithm  string
	DigestAlgorithm     string
	EncryptionAlgorithm string
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
)

type Provider struct {
	storage      Storage
	httpHandler  http.Handler
	interceptors []HttpInterceptor
	caCert       string
	caKey        string

	MetadataEndpoint *op.Endpoint
	Metadata         *md.EntityDescriptorType
	Signer           xmlsig.Signer

	IdentityProvider *IdentityProvider
}

func NewProvider(
	conf ProviderConfig,
	command *command.Commands,
	query *query.Queries,
	repo repository.Repository,
	keyConfig systemdefaults.KeyConfig,
	es *eventstore.Eventstore,
	projections types.SQL,
	certChan <-chan interface{},
	localDevMode bool,
) (*Provider, error) {
	metricTypes := []metrics.MetricType{metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode, metrics.MetricTypeTotalCount}
	cookieHandler, err := middleware.NewUserAgentHandler(conf.UserAgentCookieConfig, id.SonyFlakeGenerator, localDevMode)
	if err != nil {
		return nil, err
	}

	encAlg, err := crypto.NewAESCrypto(keyConfig.EncryptionConfig)
	if err != nil {
		return nil, err
	}
	sqlClient, err := projections.Start()
	if err != nil {
		return nil, err
	}

	storage := &ProviderStorage{
		certChan:                  certChan,
		certificateAlgorithm:      conf.CertificateAlgorithm,
		encAlg:                    encAlg,
		locker:                    crdb.NewLocker(sqlClient, locksTable, signingKey),
		eventstore:                es,
		certificateRotationCheck:  keyConfig.CertificateRotationCheck.Duration,
		certificateGracefulPeriod: keyConfig.CertificateGracefulPeriod.Duration,
		repo:                      repo,
		command:                   command,
		query:                     query,
		defaultLoginURL:           conf.StorageConfig.DefaultLoginURL,
	}

	getCACert(storage)
	cert, key := getMetadataCert(storage)

	certPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert,
		},
	)

	keyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	tlsCert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return nil, err
	}

	signer, err := xmlsig.NewSignerWithOptions(tlsCert, xmlsig.SignerOptions{
		SignatureAlgorithm: conf.Metadata.SignatureAlgorithm,
		DigestAlgorithm:    conf.Metadata.DigestAlgorithm,
	})
	if err != nil {
		return nil, err
	}

	metadata := op.NewEndpointWithURL(conf.Metadata.Path, conf.Metadata.URL)
	idp, err := NewIdentityProvider(
		&metadata,
		conf.IDP,
		storage,
	)
	if err != nil {
		return nil, err
	}

	prov := &Provider{
		MetadataEndpoint: &metadata,
		Metadata:         conf.getMetadata(idp),
		Signer:           signer,
		storage:          storage,
		IdentityProvider: idp,
		interceptors: []HttpInterceptor{
			middleware.MetricsHandler(metricTypes),
			middleware.TelemetryHandler(),
			middleware.NoCacheInterceptor,
			cookieHandler,
			http_utils.CopyHeadersToContext,
		},
	}

	prov.httpHandler = CreateRouter(prov, prov.interceptors...)

	return prov, nil
}

func (p *Provider) HttpHandler() http.Handler {
	return p.httpHandler
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
func getCACert(storage Storage) ([]byte, *rsa.PrivateKey) {
	ctx := context.Background()
	certAndKeyCh := make(chan key.CertificateAndKey)
	go storage.GetCA(ctx, certAndKeyCh)
	for {
		select {
		case <-ctx.Done():
			//TODO
		case certAndKey := <-certAndKeyCh:
			if certAndKey.Key.Key == nil || certAndKey.Certificate.Key == nil {
				logging.Log("OP-DAvt4").Warn("signer has no key")
				continue
			}
			certWebKey := certAndKey.Certificate.Key.(jose.JSONWebKey)
			keyWebKey := certAndKey.Key.Key.(jose.JSONWebKey)

			return certWebKey.Key.([]byte), keyWebKey.Key.(*rsa.PrivateKey)
		}
	}
}

func getMetadataCert(storage Storage) ([]byte, *rsa.PrivateKey) {
	ctx := context.Background()
	certAndKeyCh := make(chan key.CertificateAndKey)
	go storage.GetMetadataSigningKey(ctx, certAndKeyCh)

	for {
		select {
		case <-ctx.Done():
			//TODO
		case certAndKey := <-certAndKeyCh:
			if certAndKey.Key.Key == nil || certAndKey.Certificate.Key == nil {
				logging.Log("OP-DAvt4").Warn("signer has no key")
				continue
			}
			certWebKey := certAndKey.Certificate.Key.(jose.JSONWebKey)
			keyWebKey := certAndKey.Key.Key.(jose.JSONWebKey)

			return certWebKey.Key.([]byte), keyWebKey.Key.(*rsa.PrivateKey)
		}
	}
}

type HttpInterceptor func(http.Handler) http.Handler

func CreateRouter(p *Provider, interceptors ...HttpInterceptor) *mux.Router {
	intercept := buildInterceptor(interceptors...)
	router := mux.NewRouter()
	router.Use(handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"authorization", "content-type"}),
		handlers.AllowedOriginValidator(allowAllOrigins),
	))
	router.HandleFunc(healthEndpoint, healthHandler)
	router.HandleFunc(readinessEndpoint, readyHandler(p.Probes()))
	router.HandleFunc(p.MetadataEndpoint.Relative(), p.metadataHandle)

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
