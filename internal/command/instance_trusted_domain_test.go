package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_AddTrustedDomain(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		trustedDomain string
	}
	type want struct {
		details *domain.ObjectDetails
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "empty domain, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: "",
			},
			want: want{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-Stk21", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid domain, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: "&.com",
			},
			want: want{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-S3v3w", "Errors.Domain.InvalidCharacter"),
			},
		},
		{
			name: "domain already exists, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewTrustedDomainAddedEvent(context.Background(),
								&instance.NewAggregate("instanceID").Aggregate, "domain.com"),
						),
					),
				),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: "domain.com",
			},
			want: want{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMA-hg42a", "Errors.Instance.Domain.AlreadyExists"),
			},
		},
		{
			name: "domain add ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewTrustedDomainAddedEvent(context.Background(),
							&instance.NewAggregate("instanceID").Aggregate, "domain.com"),
					),
				),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: "domain.com",
			},
			want: want{
				details: &domain.ObjectDetails{
					ResourceOwner: "instanceID",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := c.AddTrustedDomain(tt.args.ctx, tt.args.trustedDomain)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.details, got)
		})
	}
}

func TestCommands_RemoveTrustedDomain(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		trustedDomain string
	}
	type want struct {
		details *domain.ObjectDetails
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "domain does not exists, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: "domain.com",
			},
			want: want{
				err: zerrors.ThrowNotFound(nil, "COMMA-de3z9", "Errors.Instance.Domain.NotFound"),
			},
		},
		{
			name: "domain remove ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewTrustedDomainAddedEvent(context.Background(),
								&instance.NewAggregate("instanceID").Aggregate, "domain.com"),
						),
					),
					expectPush(
						instance.NewTrustedDomainRemovedEvent(context.Background(),
							&instance.NewAggregate("instanceID").Aggregate, "domain.com"),
					),
				),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				trustedDomain: "domain.com",
			},
			want: want{
				details: &domain.ObjectDetails{
					ResourceOwner: "instanceID",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := c.RemoveTrustedDomain(tt.args.ctx, tt.args.trustedDomain)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.details, got)
		})
	}
}

//
//func TestCommands_RemoveTrustedDomain(t *testing.T) {
//	type fields struct {
//		httpClient                      *http.Client
//		jobs                            sync.WaitGroup
//		checkPermission                 domain.PermissionCheck
//		newEncryptedCode                encrypedCodeFunc
//		newEncryptedCodeWithDefault     encryptedCodeWithDefaultFunc
//		newHashedSecret                 hashedSecretFunc
//		eventstore                      *eventstore.Eventstore
//		static                          static.Storage
//		idGenerator                     id.Generator
//		zitadelRoles                    []authz.RoleMapping
//		externalDomain                  string
//		externalSecure                  bool
//		externalPort                    uint16
//		idpConfigEncryption             crypto.EncryptionAlgorithm
//		smtpEncryption                  crypto.EncryptionAlgorithm
//		smsEncryption                   crypto.EncryptionAlgorithm
//		userEncryption                  crypto.EncryptionAlgorithm
//		userPasswordHasher              *crypto.Hasher
//		secretHasher                    *crypto.Hasher
//		machineKeySize                  int
//		applicationKeySize              int
//		domainVerificationAlg           crypto.EncryptionAlgorithm
//		domainVerificationGenerator     crypto.Generator
//		domainVerificationValidator     func(domain, token, verifier string, checkType api_http.CheckType) error
//		sessionTokenCreator             func(sessionID string) (id string, token string, err error)
//		sessionTokenVerifier            func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)
//		defaultAccessTokenLifetime      time.Duration
//		defaultRefreshTokenLifetime     time.Duration
//		defaultRefreshTokenIdleLifetime time.Duration
//		multifactors                    domain.MultifactorConfigs
//		webauthnConfig                  *webauthn_helper.Config
//		keySize                         int
//		keyAlgorithm                    crypto.EncryptionAlgorithm
//		certificateAlgorithm            crypto.EncryptionAlgorithm
//		certKeySize                     int
//		privateKeyLifetime              time.Duration
//		publicKeyLifetime               time.Duration
//		certificateLifetime             time.Duration
//		defaultSecretGenerators         *SecretGenerators
//		samlCertificateAndKeyGenerator  func(id string) ([]byte, []byte, error)
//		GrpcMethodExisting              func(method string) bool
//		GrpcServiceExisting             func(method string) bool
//		ActionFunctionExisting          func(function string) bool
//		EventExisting                   func(event string) bool
//		EventGroupExisting              func(group string) bool
//		GenerateDomain                  func(instanceName, domain string) (string, error)
//	}
//	type args struct {
//		ctx           context.Context
//		trustedDomain string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *domain.ObjectDetails
//		wantErr assert.ErrorAssertionFunc
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Commands{
//				httpClient:                      tt.fields.httpClient,
//				jobs:                            tt.fields.jobs,
//				checkPermission:                 tt.fields.checkPermission,
//				newEncryptedCode:                tt.fields.newEncryptedCode,
//				newEncryptedCodeWithDefault:     tt.fields.newEncryptedCodeWithDefault,
//				newHashedSecret:                 tt.fields.newHashedSecret,
//				eventstore:                      tt.fields.eventstore,
//				static:                          tt.fields.static,
//				idGenerator:                     tt.fields.idGenerator,
//				zitadelRoles:                    tt.fields.zitadelRoles,
//				externalDomain:                  tt.fields.externalDomain,
//				externalSecure:                  tt.fields.externalSecure,
//				externalPort:                    tt.fields.externalPort,
//				idpConfigEncryption:             tt.fields.idpConfigEncryption,
//				smtpEncryption:                  tt.fields.smtpEncryption,
//				smsEncryption:                   tt.fields.smsEncryption,
//				userEncryption:                  tt.fields.userEncryption,
//				userPasswordHasher:              tt.fields.userPasswordHasher,
//				secretHasher:                    tt.fields.secretHasher,
//				machineKeySize:                  tt.fields.machineKeySize,
//				applicationKeySize:              tt.fields.applicationKeySize,
//				domainVerificationAlg:           tt.fields.domainVerificationAlg,
//				domainVerificationGenerator:     tt.fields.domainVerificationGenerator,
//				domainVerificationValidator:     tt.fields.domainVerificationValidator,
//				sessionTokenCreator:             tt.fields.sessionTokenCreator,
//				sessionTokenVerifier:            tt.fields.sessionTokenVerifier,
//				defaultAccessTokenLifetime:      tt.fields.defaultAccessTokenLifetime,
//				defaultRefreshTokenLifetime:     tt.fields.defaultRefreshTokenLifetime,
//				defaultRefreshTokenIdleLifetime: tt.fields.defaultRefreshTokenIdleLifetime,
//				multifactors:                    tt.fields.multifactors,
//				webauthnConfig:                  tt.fields.webauthnConfig,
//				keySize:                         tt.fields.keySize,
//				keyAlgorithm:                    tt.fields.keyAlgorithm,
//				certificateAlgorithm:            tt.fields.certificateAlgorithm,
//				certKeySize:                     tt.fields.certKeySize,
//				privateKeyLifetime:              tt.fields.privateKeyLifetime,
//				publicKeyLifetime:               tt.fields.publicKeyLifetime,
//				certificateLifetime:             tt.fields.certificateLifetime,
//				defaultSecretGenerators:         tt.fields.defaultSecretGenerators,
//				samlCertificateAndKeyGenerator:  tt.fields.samlCertificateAndKeyGenerator,
//				GrpcMethodExisting:              tt.fields.GrpcMethodExisting,
//				GrpcServiceExisting:             tt.fields.GrpcServiceExisting,
//				ActionFunctionExisting:          tt.fields.ActionFunctionExisting,
//				EventExisting:                   tt.fields.EventExisting,
//				EventGroupExisting:              tt.fields.EventGroupExisting,
//				GenerateDomain:                  tt.fields.GenerateDomain,
//			}
//			got, err := c.RemoveTrustedDomain(tt.args.ctx, tt.args.trustedDomain)
//			if !tt.wantErr(t, err, fmt.Sprintf("RemoveTrustedDomain(%v, %v)", tt.args.ctx, tt.args.trustedDomain)) {
//				return
//			}
//			assert.Equalf(t, tt.want, got, "RemoveTrustedDomain(%v, %v)", tt.args.ctx, tt.args.trustedDomain)
//		})
//	}
//}
