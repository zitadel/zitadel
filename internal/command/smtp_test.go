package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id.Generator
		alg         crypto.EncryptionAlgorithm
	}
	type args struct {
		instanceID string
		smtp       *AddSMTPConfig
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
				smtp: &AddSMTPConfig{},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-PQN0wsqSyi", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "smtp config, custom domain not existing",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								true,
							),
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "configid"),
				alg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				smtp: &AddSMTPConfig{
					ResourceOwner: "INSTANCE",
					From:          "from@domain.ch",
					Host:          "host:587",
					User:          "user",
					Password:      "password",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "INST-xtWIiR2ZbR", "Errors.SMTPConfig.SenderAdressNotCustomDomain"))
				},
			},
		},
		{
			name: "add smtp config, from empty",
			fields: fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "configid"),
			},
			args: args{
				smtp: &AddSMTPConfig{
					ResourceOwner: "INSTANCE",
					From:          "   ",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-SAAFpV8VKV", "Errors.Invalid.Argument"))
				},
			},
		},
		{
			name: "add smtp config, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"domain.ch",
								false,
							),
						),
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true, true, false,
							),
						),
					),
					expectPush(
						instance.NewSMTPConfigAddedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"configid",
							"test",
							true,
							"from@domain.ch",
							"name",
							"",
							"host:587",
							"user",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("password"),
							},
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "configid"),
				alg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				smtp: &AddSMTPConfig{
					ResourceOwner: "INSTANCE",
					Description:   "test",
					Tls:           true,
					From:          "from@domain.ch",
					FromName:      "name",
					Host:          "host:587",
					User:          "user",
					Password:      "password",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "add smtp config with reply to address, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"domain.ch",
								false,
							),
						),
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true, true, false,
							),
						),
					),
					expectPush(
						instance.NewSMTPConfigAddedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"configid",
							"test",
							true,
							"from@domain.ch",
							"name",
							"replyto@domain.ch",
							"host:587",
							"user",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("password"),
							},
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "configid"),
				alg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				smtp: &AddSMTPConfig{
					ResourceOwner:  "INSTANCE",
					Description:    "test",
					Tls:            true,
					From:           "from@domain.ch",
					FromName:       "name",
					ReplyToAddress: "replyto@domain.ch",
					Host:           "host:587",
					User:           "user",
					Password:       "password",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "smtp config, port is missing",
			fields: fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "configid"),
				alg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				smtp: &AddSMTPConfig{
					ResourceOwner: "INSTANCE",
					Description:   "test",
					Tls:           true,
					From:          "from@domain.ch",
					FromName:      "name",
					Host:          "host",
					User:          "user",
					Password:      "password",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-EvAtufIinh", "Errors.Invalid.Argument"))
				},
			},
		},
		{
			name: "smtp config, host is empty",
			fields: fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "configid"),
				alg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				smtp: &AddSMTPConfig{
					ResourceOwner: "INSTANCE",
					Description:   "test",
					Tls:           true,
					From:          "from@domain.ch",
					FromName:      "name",
					Host:          "   ",
					User:          "user",
					Password:      "password",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-EvAtufIinh", "Errors.Invalid.Argument"))
				},
			},
		},
		{
			name: "add smtp config, ipv6 works",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"domain.ch",
								false,
							),
						),
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true, true, false,
							),
						),
					),
					expectPush(
						instance.NewSMTPConfigAddedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"configid",
							"test",
							true,
							"from@domain.ch",
							"name",
							"",
							"[2001:db8::1]:2525",
							"user",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("password"),
							},
						),
					),
				),
				alg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "configid"),
			},
			args: args{
				smtp: &AddSMTPConfig{
					ResourceOwner: "INSTANCE",
					Description:   "test",
					Tls:           true,
					From:          "from@domain.ch",
					FromName:      "name",
					Host:          "[2001:db8::1]:2525",
					User:          "user",
					Password:      "password",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
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
			err := r.AddSMTPConfig(context.Background(), tt.args.smtp)
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

func TestCommandSide_ChangeSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		smtp *ChangeSMTPConfig
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
			name: "resourceowner empty, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				smtp: &ChangeSMTPConfig{},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-jwA8gxldy3", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "id empty, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				smtp: &ChangeSMTPConfig{
					ResourceOwner: "INSTANCE",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-2JPlSRzuHy", "Errors.IDMissing"))
				},
			},
		},
		{
			name: "empty config, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				smtp: &ChangeSMTPConfig{
					ResourceOwner: "INSTANCE",
					ID:            "configID",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-gyPUXOTA4N", "Errors.Invalid.Argument"))
				},
			},
		},
		{
			name: "smtp not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				smtp: &ChangeSMTPConfig{
					ResourceOwner: "INSTANCE",
					ID:            "ID",
					Description:   "test",
					Tls:           true,
					From:          "from@domain.ch",
					FromName:      "name",
					Host:          "host:587",
					User:          "user",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-j5IDFtt3T1", "Errors.SMTPConfig.NotFound"))
				},
			},
		},
		{
			name: "smtp domain not matched",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"domain.ch",
								false,
							),
						),
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true, true, true,
							),
						),
						eventFromEventPusher(
							instance.NewSMTPConfigAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
								"test",
								true,
								"from@domain.ch",
								"name",
								"",
								"host:587",
								"user",
								&crypto.CryptoValue{},
							),
						),
					),
				),
			},
			args: args{
				smtp: &ChangeSMTPConfig{
					ResourceOwner: "INSTANCE",
					ID:            "ID",
					Description:   "test",
					Tls:           true,
					From:          "from@wrongdomain.ch",
					FromName:      "name",
					Host:          "host:587",
					User:          "user",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"domain.ch",
								false,
							),
						),
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true, true, true,
							),
						),
						eventFromEventPusher(
							instance.NewSMTPConfigAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
								"test",
								true,
								"from@domain.ch",
								"name",
								"",
								"host:587",
								"user",
								&crypto.CryptoValue{},
							),
						),
					),
				),
			},
			args: args{
				smtp: &ChangeSMTPConfig{
					ResourceOwner: "INSTANCE",
					ID:            "ID",
					Description:   "test",
					Tls:           true,
					From:          "from@domain.ch",
					FromName:      "name",
					Host:          "host:587",
					User:          "user",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "smtp config change, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"domain.ch",
								false,
							),
						),
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true, true, true,
							),
						),
						eventFromEventPusher(
							instance.NewSMTPConfigAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
								"",
								true,
								"from@domain.ch",
								"name",
								"",
								"host:587",
								"user",
								&crypto.CryptoValue{},
							),
						),
					),
					expectPush(
						newSMTPConfigChangedEvent(
							context.Background(),
							"ID",
							"test",
							false,
							"from2@domain.ch",
							"name2",
							"replyto@domain.ch",
							"host2:587",
							"user2",
						),
					),
				),
			},
			args: args{
				smtp: &ChangeSMTPConfig{
					ResourceOwner:  "INSTANCE",
					ID:             "ID",
					Description:    "test",
					Tls:            false,
					From:           "from2@domain.ch",
					FromName:       "name2",
					ReplyToAddress: "replyto@domain.ch",
					Host:           "host2:587",
					User:           "user2",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "smtp config, port is missing",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				smtp: &ChangeSMTPConfig{
					ResourceOwner: "INSTANCE",
					ID:            "ID",
					Description:   "test",
					Tls:           true,
					From:          "from@domain.ch",
					FromName:      "name",
					Host:          "host",
					User:          "user",
					Password:      "password",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-kZNVkuL32L", "Errors.Invalid.Argument"))
				},
			},
		},
		{
			name: "smtp config, host is empty",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				smtp: &ChangeSMTPConfig{
					ResourceOwner: "INSTANCE",
					ID:            "ID",
					Description:   "test",
					Tls:           true,
					From:          "from@domain.ch",
					FromName:      "name",
					Host:          "   ",
					User:          "user",
					Password:      "password",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-kZNVkuL32L", "Errors.Invalid.Argument"))
				},
			},
		},
		{
			name: "smtp config change, ipv6 works",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"domain.ch",
								false,
							),
						),
						eventFromEventPusher(
							instance.NewDomainPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true, true, true,
							),
						),
						eventFromEventPusher(
							instance.NewSMTPConfigAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
								"",
								true,
								"from@domain.ch",
								"name",
								"",
								"host:587",
								"user",
								&crypto.CryptoValue{},
							),
						),
					),
					expectPush(
						newSMTPConfigChangedEvent(
							context.Background(),
							"ID",
							"test",
							false,
							"from2@domain.ch",
							"name2",
							"replyto@domain.ch",
							"[2001:db8::1]:2525",
							"user2",
						),
					),
				),
			},
			args: args{
				smtp: &ChangeSMTPConfig{
					ResourceOwner:  "INSTANCE",
					ID:             "ID",
					Description:    "test",
					Tls:            false,
					From:           "from2@domain.ch",
					FromName:       "name2",
					ReplyToAddress: "replyto@domain.ch",
					Host:           "[2001:db8::1]:2525",
					User:           "user2",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			err := r.ChangeSMTPConfig(context.Background(), tt.args.smtp)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, tt.args.smtp.Details)
			}
		})
	}
}

func TestCommandSide_ChangeSMTPConfigPassword(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
		alg        crypto.EncryptionAlgorithm
	}
	type args struct {
		instanceID string
		id         string
		password   string
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
			name: "smtp config, error resourceOwner empty",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-gHAyvUXCAF", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "smtp config, error id empty",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				instanceID: "INSTANCE",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-BCkAf7LcJA", "Errors.IDMissing"))
				},
			},
		},
		{
			name: "smtp config, error not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				instanceID: "INSTANCE",
				password:   "",
				id:         "ID",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-rDHzqjGuKQ", "Errors.SMTPConfig.NotFound"))
				},
			},
		},
		{
			name: "change smtp config password, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMTPConfigAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
								"test",
								true,
								"from",
								"name",
								"",
								"host:587",
								"user",
								&crypto.CryptoValue{},
							),
						),
						eventFromEventPusher(
							instance.NewSMTPConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
							),
						),
					),
					expectPush(
						instance.NewSMTPConfigPasswordChangedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"ID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("password"),
							},
						),
					),
				),
				alg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				password:   "password",
				id:         "ID",
				instanceID: "INSTANCE",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:     tt.fields.eventstore(t),
				smtpEncryption: tt.fields.alg,
			}
			got, err := r.ChangeSMTPConfigPassword(context.Background(), tt.args.instanceID, tt.args.id, tt.args.password)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddSMTPConfigHTTP(t *testing.T) {
	type fields struct {
		eventstore                  func(t *testing.T) *eventstore.Eventstore
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
		defaultSecretGenerators     *SecretGenerators
		idGenerator                 id.Generator
	}
	type args struct {
		http *AddSMTPConfigHTTP
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
			name: "add smtp config, resourceowner empty",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				http: &AddSMTPConfigHTTP{},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-FTNDXc8ACS", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "add smtp config, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						instance.NewSMTPConfigHTTPAddedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
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
				http: &AddSMTPConfigHTTP{
					ResourceOwner: "INSTANCE",
					Description:   "test",
					Endpoint:      "endpoint",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
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
			err := r.AddSMTPConfigHTTP(context.Background(), tt.args.http)
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

func TestCommandSide_ChangeSMTPConfigHTTP(t *testing.T) {
	type fields struct {
		eventstore                  func(t *testing.T) *eventstore.Eventstore
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
		defaultSecretGenerators     *SecretGenerators
	}
	type args struct {
		http *ChangeSMTPConfigHTTP
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
	}{{
		name: "resourceowner empty, precondition error",
		fields: fields{
			eventstore: expectEventstore(),
		},
		args: args{
			http: &ChangeSMTPConfigHTTP{},
		},
		res: res{
			err: func(err error) bool {
				return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-k7QCGOWyJA", "Errors.ResourceOwnerMissing"))
			},
		},
	},
		{
			name: "id empty, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				http: &ChangeSMTPConfigHTTP{
					ResourceOwner: "INSTANCE",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-2MHkV8ObWo", "Errors.IDMissing"))
				},
			},
		},
		{
			name: "smtp not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				http: &ChangeSMTPConfigHTTP{
					ResourceOwner: "INSTANCE",
					ID:            "ID",
					Description:   "test",
					Endpoint:      "endpoint",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-xIrdledqv4", "Errors.SMTPConfig.NotFound"))
				},
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMTPConfigHTTPAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
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
				),
			},
			args: args{
				http: &ChangeSMTPConfigHTTP{
					ResourceOwner: "INSTANCE",
					ID:            "ID",
					Description:   "test",
					Endpoint:      "endpoint",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "smtp config change, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMTPConfigHTTPAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
								"",
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
					expectPush(
						newSMTPConfigHTTPChangedEvent(
							context.Background(),
							"ID",
							"test",
							"endpoint2",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("87654321"),
							},
						),
					),
				),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("87654321", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args: args{
				http: &ChangeSMTPConfigHTTP{
					ResourceOwner:        "INSTANCE",
					ID:                   "ID",
					Description:          "test",
					Endpoint:             "endpoint2",
					ExpirationSigningKey: true,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:                  tt.fields.eventstore(t),
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
			}
			err := r.ChangeSMTPConfigHTTP(context.Background(), tt.args.http)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, tt.args.http.Details)
			}
		})
	}
}

func TestCommandSide_ActivateSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
		alg        crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx        context.Context
		instanceID string
		id         string
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
	}{{
		name: "resourceowner empty, precondition error",
		fields: fields{
			eventstore: expectEventstore(),
		},
		args: args{
			instanceID: "",
		},
		res: res{
			err: func(err error) bool {
				return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-h5htMCebv3", "Errors.ResourceOwnerMissing"))
			},
		},
	},
		{
			name: "id empty, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				instanceID: "INSTANCE",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-1hPl6oVMJa", "Errors.IDMissing"))
				},
			},
		},
		{
			name: "smtp not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				id:         "id",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-E9K20hxOS9", "Errors.SMTPConfig.NotFound"))
				},
			},
		},
		{
			name: "activate smtp config, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMTPConfigAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
								"test",
								true,
								"from",
								"name",
								"",
								"host:587",
								"user",
								&crypto.CryptoValue{},
							),
						),
					),
					expectPush(
						instance.NewSMTPConfigActivatedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"ID",
						),
					),
				),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				id:         "ID",
				instanceID: "INSTANCE",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "activate smtp config, already active",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMTPConfigAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
								"test",
								true,
								"from",
								"name",
								"",
								"host:587",
								"user",
								&crypto.CryptoValue{},
							),
						),
						eventFromEventPusher(
							instance.NewSMTPConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
							),
						),
					),
				),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				id:         "ID",
				instanceID: "INSTANCE",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-vUHBSmBzaw", "Errors.SMTPConfig.AlreadyActive"))
				},
			},
		},
		{
			name: "activate smtp config http, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMTPConfigHTTPAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
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
					expectPush(
						instance.NewSMTPConfigActivatedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"ID",
						),
					),
				),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				id:         "ID",
				instanceID: "INSTANCE",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:     tt.fields.eventstore(t),
				smtpEncryption: tt.fields.alg,
			}
			got, err := r.ActivateSMTPConfig(tt.args.ctx, tt.args.instanceID, tt.args.id)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_DeactivateSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
		alg        crypto.EncryptionAlgorithm
	}
	type args struct {
		instanceID string
		id         string
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
	}{{
		name: "resourceOwner empty, precondition error",
		fields: fields{
			eventstore: expectEventstore(),
		},
		args: args{},
		res: res{
			err: func(err error) bool {
				return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-pvNHou89Tw", "Errors.ResourceOwnerMissing"))
			},
		},
	},
		{
			name: "id empty, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				instanceID: "INSTANCE",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-jLTIMrtApO", "Errors.IDMissing"))
				},
			},
		},
		{
			name: "smtp not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				instanceID: "INSTANCE",
				id:         "id",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-k39PJ", "Errors.SMTPConfig.NotFound"))
				},
			},
		},
		{
			name: "deactivate smtp config, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMTPConfigAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
								"test",
								true,
								"from",
								"name",
								"",
								"host:587",
								"user",
								&crypto.CryptoValue{},
							),
						),
						eventFromEventPusher(
							instance.NewSMTPConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
							),
						),
					),
					expectPush(
						instance.NewSMTPConfigDeactivatedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"ID",
						),
					),
				),
			},
			args: args{
				id:         "ID",
				instanceID: "INSTANCE",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "deactivate smtp config, already deactivated",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMTPConfigAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
								"test",
								true,
								"from",
								"name",
								"",
								"host:587",
								"user",
								&crypto.CryptoValue{},
							),
						),
						eventFromEventPusher(
							instance.NewSMTPConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
							),
						),
						eventFromEventPusher(
							instance.NewSMTPConfigDeactivatedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
							),
						),
					),
				),
			},
			args: args{
				id:         "ID",
				instanceID: "INSTANCE",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-km8g3", "Errors.SMTPConfig.AlreadyDeactivated"))
				},
			},
		},
		{
			name: "deactivate smtp config http, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMTPConfigHTTPAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
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
						eventFromEventPusher(
							instance.NewSMTPConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
							),
						),
					),
					expectPush(
						instance.NewSMTPConfigDeactivatedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"ID",
						),
					),
				),
			},
			args: args{
				id:         "ID",
				instanceID: "INSTANCE",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:     tt.fields.eventstore(t),
				smtpEncryption: tt.fields.alg,
			}
			got, err := r.DeactivateSMTPConfig(context.Background(), tt.args.instanceID, tt.args.id)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_RemoveSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
		alg        crypto.EncryptionAlgorithm
	}
	type args struct {
		instanceID string
		id         string
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
			name: "resourceowner empty, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				id: "ID",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-t2WsPRgGaK", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "id empty, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				instanceID: "INSTANCE",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-0ZV5whuUfu", "Errors.IDMissing"))
				},
			},
		},
		{
			name: "smtp config, error not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				instanceID: "INSTANCE",
				id:         "ID",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-09CXlTDL6w", "Errors.SMTPConfig.NotFound"))
				},
			},
		},
		{
			name: "remove smtp config, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMTPConfigAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
								"test",
								true,
								"from",
								"name",
								"",
								"host:587",
								"user",
								&crypto.CryptoValue{},
							),
						),
					),
					expectPush(
						instance.NewSMTPConfigRemovedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"ID",
						),
					),
				),
			},
			args: args{
				id:         "ID",
				instanceID: "INSTANCE",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "remove smtp config http, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMTPConfigHTTPAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
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
					expectPush(
						instance.NewSMTPConfigRemovedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"ID",
						),
					),
				),
			},
			args: args{
				id:         "ID",
				instanceID: "INSTANCE",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:     tt.fields.eventstore(t),
				smtpEncryption: tt.fields.alg,
			}
			got, err := r.RemoveSMTPConfig(context.Background(), tt.args.instanceID, tt.args.id)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_TestSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
		alg        crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx        context.Context
		instanceID string
		id         string
		email      string
		config     smtp.Config
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "id empty, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "email empty, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				instanceID: "INSTANCE",
				id:         "id",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "if password is empty, smtp id must not",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				id:         "",
				email:      "email",
				config: smtp.Config{
					From:     "test@example,com",
					FromName: "Test",
					SMTP: smtp.SMTP{
						User:     "user",
						Password: "",
						Host:     "example.com:2525",
					},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "password empty and smtp config not found, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				id:         "id",
				email:      "email",
				config: smtp.Config{
					From:     "test@example,com",
					FromName: "Test",
					SMTP: smtp.SMTP{
						User:     "user",
						Password: "",
						Host:     "example.com:2525",
					},
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "valid new smtp config, wrong auth, ok",
			fields: fields{
				eventstore: expectEventstore(),
				alg:        crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				email:      "email",
				config: smtp.Config{
					From:     "test@example.com",
					FromName: "Test",
					SMTP: smtp.SMTP{
						User:     "user",
						Password: "password",
						Host:     "mail.smtp2go.com:2525",
					},
				},
			},
			res: res{
				err: zerrors.IsInternal,
			},
		},
		{
			name: "valid smtp config using stored password, wrong auth, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							instance.NewSMTPConfigAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
								"test",
								true,
								"from",
								"name",
								"",
								"mail.smtp2go.com:2525",
								"user",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("password"),
								},
							),
						),
					),
				),
				alg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				email:      "email",
				id:         "ID",
				config: smtp.Config{
					From:     "test@example.com",
					FromName: "Test",
					SMTP: smtp.SMTP{
						User:     "user",
						Password: "",
						Host:     "mail.smtp2go.com:2525",
					},
				},
			},
			res: res{
				err: zerrors.IsInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:     tt.fields.eventstore(t),
				smtpEncryption: tt.fields.alg,
			}
			err := r.TestSMTPConfig(tt.args.ctx, tt.args.instanceID, tt.args.id, tt.args.email, &tt.args.config)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCommandSide_TestSMTPConfigById(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		alg        crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx        context.Context
		instanceID string
		id         string
		email      string
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "id empty, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "email empty, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				instanceID: "INSTANCE",
				id:         "id",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "smtp config not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				id:         "id",
				email:      "email",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "valid smtp config, wrong auth, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewSMTPConfigAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"ID",
								"test",
								true,
								"from",
								"name",
								"",
								"mail.smtp2go.com:2525",
								"user",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("password"),
								},
							),
						),
					),
				),
				alg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				id:         "ID",
				instanceID: "INSTANCE",
				email:      "test@example.com",
			},
			res: res{
				err: zerrors.IsInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:     tt.fields.eventstore,
				smtpEncryption: tt.fields.alg,
			}
			err := r.TestSMTPConfigById(tt.args.ctx, tt.args.instanceID, tt.args.id, tt.args.email)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func newSMTPConfigChangedEvent(ctx context.Context, id, description string, tls bool, fromAddress, fromName, replyTo, host, user string) *instance.SMTPConfigChangedEvent {
	changes := []instance.SMTPConfigChanges{
		instance.ChangeSMTPConfigDescription(description),
		instance.ChangeSMTPConfigTLS(tls),
		instance.ChangeSMTPConfigFromAddress(fromAddress),
		instance.ChangeSMTPConfigFromName(fromName),
		instance.ChangeSMTPConfigReplyToAddress(replyTo),
		instance.ChangeSMTPConfigSMTPHost(host),
		instance.ChangeSMTPConfigSMTPUser(user),
	}
	event, _ := instance.NewSMTPConfigChangeEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		id,
		changes,
	)
	return event
}

func newSMTPConfigHTTPChangedEvent(ctx context.Context, id, description, endpoint string, signingKey *crypto.CryptoValue) *instance.SMTPConfigHTTPChangedEvent {
	changes := []instance.SMTPConfigHTTPChanges{
		instance.ChangeSMTPConfigHTTPDescription(description),
		instance.ChangeSMTPConfigHTTPEndpoint(endpoint),
		instance.ChangeSMTPConfigHTTPSigningKey(signingKey),
	}
	event, _ := instance.NewSMTPConfigHTTPChangeEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		id,
		changes,
	)
	return event
}
