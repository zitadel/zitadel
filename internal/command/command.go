package command

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	api_http "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/command/preparation"
	sd "github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/static"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	webauthn_helper "github.com/zitadel/zitadel/internal/webauthn"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Commands struct {
	httpClient *http.Client

	jobs sync.WaitGroup

	checkPermission             domain.PermissionCheck
	newEncryptedCode            encrypedCodeFunc
	newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
	newHashedSecret             hashedSecretFunc

	eventstore     *eventstore.Eventstore
	static         static.Storage
	zitadelRoles   []authz.RoleMapping
	externalDomain string
	externalSecure bool
	externalPort   uint16

	idpConfigEncryption             crypto.EncryptionAlgorithm
	smtpEncryption                  crypto.EncryptionAlgorithm
	smsEncryption                   crypto.EncryptionAlgorithm
	userEncryption                  crypto.EncryptionAlgorithm
	userPasswordHasher              *crypto.Hasher
	secretHasher                    *crypto.Hasher
	machineKeySize                  int
	applicationKeySize              int
	domainVerificationAlg           crypto.EncryptionAlgorithm
	domainVerificationGenerator     crypto.Generator
	domainVerificationValidator     func(domain, token, verifier string, checkType api_http.CheckType) error
	sessionTokenCreator             func(sessionID string) (id string, token string, err error)
	sessionTokenVerifier            func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)
	defaultAccessTokenLifetime      time.Duration
	defaultRefreshTokenLifetime     time.Duration
	defaultRefreshTokenIdleLifetime time.Duration

	multifactors            domain.MultifactorConfigs
	webauthnConfig          *webauthn_helper.Config
	keySize                 int
	keyAlgorithm            crypto.EncryptionAlgorithm
	certificateAlgorithm    crypto.EncryptionAlgorithm
	certKeySize             int
	privateKeyLifetime      time.Duration
	publicKeyLifetime       time.Duration
	certificateLifetime     time.Duration
	defaultSecretGenerators *SecretGenerators

	samlCertificateAndKeyGenerator func(id string) ([]byte, []byte, error)

	GrpcMethodExisting     func(method string) bool
	GrpcServiceExisting    func(method string) bool
	ActionFunctionExisting func(function string) bool
	EventExisting          func(event string) bool
	EventGroupExisting     func(group string) bool

	GenerateDomain func(instanceName, domain string) (string, error)
}

func StartCommands(
	es *eventstore.Eventstore,
	defaults sd.SystemDefaults,
	zitadelRoles []authz.RoleMapping,
	staticStore static.Storage,
	webAuthN *webauthn_helper.Config,
	externalDomain string,
	externalSecure bool,
	externalPort uint16,
	idpConfigEncryption, otpEncryption, smtpEncryption, smsEncryption, userEncryption, domainVerificationEncryption, oidcEncryption, samlEncryption crypto.EncryptionAlgorithm,
	httpClient *http.Client,
	permissionCheck domain.PermissionCheck,
	sessionTokenVerifier func(ctx context.Context, sessionToken string, sessionID string, tokenID string) (err error),
	defaultAccessTokenLifetime,
	defaultRefreshTokenLifetime,
	defaultRefreshTokenIdleLifetime time.Duration,
	defaultSecretGenerators *SecretGenerators,
) (repo *Commands, err error) {
	if externalDomain == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Df21s", "no external domain specified")
	}
	// reuse the oidcEncryption to be able to handle both tokens in the interceptor later on
	sessionAlg := oidcEncryption

	secretHasher, err := defaults.SecretHasher.NewHasher()
	if err != nil {
		return nil, fmt.Errorf("secret hasher: %w", err)
	}
	userPasswordHasher, err := defaults.PasswordHasher.NewHasher()
	if err != nil {
		return nil, fmt.Errorf("password hasher: %w", err)
	}
	repo = &Commands{
		eventstore:                      es,
		static:                          staticStore,
		zitadelRoles:                    zitadelRoles,
		externalDomain:                  externalDomain,
		externalSecure:                  externalSecure,
		externalPort:                    externalPort,
		keySize:                         defaults.KeyConfig.Size,
		certKeySize:                     defaults.KeyConfig.CertificateSize,
		privateKeyLifetime:              defaults.KeyConfig.PrivateKeyLifetime,
		publicKeyLifetime:               defaults.KeyConfig.PublicKeyLifetime,
		certificateLifetime:             defaults.KeyConfig.CertificateLifetime,
		idpConfigEncryption:             idpConfigEncryption,
		smtpEncryption:                  smtpEncryption,
		smsEncryption:                   smsEncryption,
		userEncryption:                  userEncryption,
		userPasswordHasher:              userPasswordHasher,
		secretHasher:                    secretHasher,
		machineKeySize:                  int(defaults.SecretGenerators.MachineKeySize),
		applicationKeySize:              int(defaults.SecretGenerators.ApplicationKeySize),
		domainVerificationAlg:           domainVerificationEncryption,
		domainVerificationGenerator:     crypto.NewEncryptionGenerator(defaults.DomainVerification.VerificationGenerator, domainVerificationEncryption),
		domainVerificationValidator:     api_http.ValidateDomain,
		keyAlgorithm:                    oidcEncryption,
		certificateAlgorithm:            samlEncryption,
		webauthnConfig:                  webAuthN,
		httpClient:                      httpClient,
		checkPermission:                 permissionCheck,
		newEncryptedCode:                newEncryptedCode,
		newEncryptedCodeWithDefault:     newEncryptedCodeWithDefaultConfig,
		sessionTokenCreator:             sessionTokenCreator(sessionAlg),
		sessionTokenVerifier:            sessionTokenVerifier,
		defaultAccessTokenLifetime:      defaultAccessTokenLifetime,
		defaultRefreshTokenLifetime:     defaultRefreshTokenLifetime,
		defaultRefreshTokenIdleLifetime: defaultRefreshTokenIdleLifetime,
		defaultSecretGenerators:         defaultSecretGenerators,
		samlCertificateAndKeyGenerator:  samlCertificateAndKeyGenerator(defaults.KeyConfig.CertificateSize, defaults.KeyConfig.CertificateLifetime),
		// always true for now until we can check with an eventlist
		EventExisting: func(event string) bool { return true },
		// always true for now until we can check with an eventlist
		EventGroupExisting:     func(group string) bool { return true },
		GrpcServiceExisting:    func(service string) bool { return false },
		GrpcMethodExisting:     func(method string) bool { return false },
		ActionFunctionExisting: domain.FunctionExists(),
		multifactors: domain.MultifactorConfigs{
			OTP: domain.OTPConfig{
				CryptoMFA: otpEncryption,
				Issuer:    defaults.Multifactors.OTP.Issuer,
			},
		},
		GenerateDomain: domain.NewGeneratedInstanceDomain,
	}

	if defaultSecretGenerators != nil && defaultSecretGenerators.ClientSecret != nil {
		repo.newHashedSecret = newHashedSecretWithDefault(secretHasher, defaultSecretGenerators.ClientSecret)
	}
	return repo, nil
}

type AppendReducer interface {
	AppendEvents(...eventstore.Event)
	// TODO: Why is it allowed to return an error here?
	Reduce() error
}

func (c *Commands) pushAppendAndReduce(ctx context.Context, object AppendReducer, cmds ...eventstore.Command) error {
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return err
	}
	return AppendAndReduce(object, events...)
}

func AppendAndReduce(object AppendReducer, events ...eventstore.Event) error {
	object.AppendEvents(events...)
	return object.Reduce()
}

func queryAndReduce(ctx context.Context, filter preparation.FilterToQueryReducer, wm eventstore.QueryReducer) error {
	events, err := filter(ctx, wm.Query())
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}
	wm.AppendEvents(events...)
	return wm.Reduce()
}

type existsWriteModel interface {
	Exists() bool
	eventstore.QueryReducer
}

func exists(ctx context.Context, filter preparation.FilterToQueryReducer, wm existsWriteModel) (bool, error) {
	err := queryAndReduce(ctx, filter, wm)
	if err != nil {
		return false, err
	}
	return wm.Exists(), nil
}

func samlCertificateAndKeyGenerator(keySize int, lifetime time.Duration) func(id string) ([]byte, []byte, error) {
	return func(id string) ([]byte, []byte, error) {
		priv, pub, err := crypto.GenerateKeyPair(keySize)
		if err != nil {
			return nil, nil, err
		}

		serial, err := id_generator.NumericFromID(id)
		if err != nil {
			return nil, nil, err
		}
		now := time.Now()
		template := x509.Certificate{
			SerialNumber: big.NewInt(serial),
			Subject: pkix.Name{
				Organization: []string{"ZITADEL"},
				SerialNumber: id,
			},
			NotBefore:             now,
			NotAfter:              now.Add(lifetime),
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			BasicConstraintsValid: true,
		}

		derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, pub, priv)
		if err != nil {
			return nil, nil, zerrors.ThrowInternalf(err, "COMMAND-x92u101j", "failed to create certificate")
		}

		keyBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}
		certBlock := &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}
		return pem.EncodeToMemory(keyBlock), pem.EncodeToMemory(certBlock), nil
	}
}

// Close blocks until all async jobs are finished,
// the context expires or after eventstore.PushTimeout.
func (c *Commands) Close(ctx context.Context) error {
	if c.eventstore.PushTimeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.eventstore.PushTimeout)
		defer cancel()
	}

	done := make(chan struct{})
	go func() {
		c.jobs.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// asyncPush attempts to push events to the eventstore in a separate Go routine.
// This can be used to speed up request times when the outcome of the push is
// not important for business logic but have a pure logging function.
// For example this can be used for Secret Check Success and Failed events.
// On push error, a log line describing the error will be emitted.
func (c *Commands) asyncPush(ctx context.Context, cmds ...eventstore.Command) {
	// Create a new context, as the request scoped context might get
	// canceled before we where able to push.
	// The eventstore has its own PushTimeout setting,
	// so we don't need to have a context with timeout here.
	ctx = context.WithoutCancel(ctx)

	c.jobs.Add(1)

	go func() {
		defer c.jobs.Done()
		localCtx, span := tracing.NewSpan(ctx)

		_, err := c.eventstore.Push(localCtx, cmds...)
		if err != nil {
			for _, cmd := range cmds {
				logging.WithError(err).Warnf("could not push event %q", cmd.Type())
			}
		}

		span.EndWithError(err)
	}()
}
