package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddOrgSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id.Generator
		alg         crypto.EncryptionAlgorithm
	}
	type args struct {
		smtp *AddOrgSMTPConfig
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "resourceowner empty, invalid argument",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				smtp: &AddOrgSMTPConfig{},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-sn93Ss", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "add org smtp config, from empty",
			fields: fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "configid"),
			},
			args: args{
				smtp: &AddOrgSMTPConfig{
					ResourceOwner: "ORG1",
					From:          "   ",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-ASF3g2", "Errors.Invalid.Argument"))
				},
			},
		},
		{
			name: "add org smtp config, host missing port",
			fields: fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "configid"),
			},
			args: args{
				smtp: &AddOrgSMTPConfig{
					ResourceOwner: "ORG1",
					From:          "from@example.com",
					Host:          "host",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-gK9RE2", "Errors.Invalid.Argument"))
				},
			},
		},
		{
			name: "add org smtp config, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewOrgSMTPConfigAddedEvent(
							context.Background(),
							&org.NewAggregate("ORG1").Aggregate,
							"configid",
							"test",
							true,
							"from@example.com",
							"name",
							"",
							"host:587",
							"user",
							&instance.PlainAuth{
								Password: &crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("password"),
								},
							},
							nil,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "configid"),
				alg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				smtp: &AddOrgSMTPConfig{
					ResourceOwner: "ORG1",
					Description:   "test",
					Tls:           true,
					From:          "from@example.com",
					FromName:      "name",
					Host:          "host:587",
					User:          "user",
					PlainAuth: &PlainAuth{
						Password: "password",
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "ORG1",
				},
			},
		},
		{
			name: "add org smtp config with reply to, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewOrgSMTPConfigAddedEvent(
							context.Background(),
							&org.NewAggregate("ORG1").Aggregate,
							"configid",
							"test",
							true,
							"from@example.com",
							"name",
							"replyto@example.com",
							"host:587",
							"user",
							&instance.PlainAuth{
								Password: &crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("password"),
								},
							},
							nil,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "configid"),
				alg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				smtp: &AddOrgSMTPConfig{
					ResourceOwner:  "ORG1",
					Description:    "test",
					Tls:            true,
					From:           "from@example.com",
					FromName:       "name",
					ReplyToAddress: "replyto@example.com",
					Host:           "host:587",
					User:           "user",
					PlainAuth: &PlainAuth{
						Password: "password",
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "ORG1",
				},
			},
		},
		{
			name: "add org smtp config without auth, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewOrgSMTPConfigAddedEvent(
							context.Background(),
							&org.NewAggregate("ORG1").Aggregate,
							"configid",
							"test",
							true,
							"from@example.com",
							"name",
							"",
							"host:587",
							"",
							nil,
							nil,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "configid"),
				alg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				smtp: &AddOrgSMTPConfig{
					ResourceOwner: "ORG1",
					Description:   "test",
					Tls:           true,
					From:          "from@example.com",
					FromName:      "name",
					Host:          "host:587",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "ORG1",
				},
			},
		},
		{
			name: "add org smtp config with xoauth2, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewOrgSMTPConfigAddedEvent(
							context.Background(),
							&org.NewAggregate("ORG1").Aggregate,
							"configid",
							"test",
							true,
							"from@example.com",
							"name",
							"",
							"host:587",
							"user",
							nil,
							&instance.XOAuth2Auth{
								TokenEndpoint: "auth.example.com/token",
								Scopes:        []string{"scope"},
								ClientCredentials: &instance.XOAuth2ClientCredentials{
									ClientId: "client-id",
									ClientSecret: &crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("client-secret"),
									},
								},
							},
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "configid"),
				alg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				smtp: &AddOrgSMTPConfig{
					ResourceOwner: "ORG1",
					Description:   "test",
					Tls:           true,
					From:          "from@example.com",
					FromName:      "name",
					Host:          "host:587",
					User:          "user",
					XOAuth2Auth: &XOAuth2Auth{
						TokenEndpoint: "auth.example.com/token",
						Scopes:        []string{"scope"},
						ClientCredentialsAuth: &OAuth2ClientCredentials{
							ClientId:     "client-id",
							ClientSecret: "client-secret",
						},
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "ORG1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:     tt.fields.eventstore(t),
				idGenerator:    tt.fields.idGenerator,
				smtpEncryption: tt.fields.alg,
			}
			err := r.AddOrgSMTPConfig(context.Background(), tt.args.smtp)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, tt.args.smtp.Details)
				assert.NotEmpty(t, tt.args.smtp.ID)
			}
		})
	}
}

func TestCommandSide_AddOrgSMTPConfigHTTP(t *testing.T) {
	type fields struct {
		eventstore                  func(t *testing.T) *eventstore.Eventstore
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
		defaultSecretGenerators     *SecretGenerators
		idGenerator                 id.Generator
	}
	type args struct {
		http *AddOrgSMTPConfigHTTP
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "resourceowner empty, invalid argument",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				http: &AddOrgSMTPConfigHTTP{},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-FTNDXc1", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "add org smtp config http, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewOrgSMTPConfigHTTPAddedEvent(
							context.Background(),
							&org.NewAggregate("ORG1").Aggregate,
							"configid",
							"test",
							"endpoint",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("12345678"),
							},
						),
					),
				),
				idGenerator:                 id_mock.NewIDGeneratorExpectIDs(t, "configid"),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("12345678", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args: args{
				http: &AddOrgSMTPConfigHTTP{
					ResourceOwner: "ORG1",
					Description:   "test",
					Endpoint:      "endpoint",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "ORG1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:                  tt.fields.eventstore(t),
				idGenerator:                 tt.fields.idGenerator,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
			}
			err := r.AddOrgSMTPConfigHTTP(context.Background(), tt.args.http)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, tt.args.http.Details)
				assert.NotEmpty(t, tt.args.http.ID)
			}
		})
	}
}
