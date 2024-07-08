package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	id_mock "github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id_generator.Generator
		alg         crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx        context.Context
		instanceID string
		smtp       *smtp.Config
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
			name: "smtp config, custom domain not existing",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					From: "from@domain.ch",
					SMTP: smtp.SMTP{
						Host:     "host:587",
						User:     "user",
						Password: "password",
					},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "add smtp config, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Description: "test",
					Tls:         true,
					From:        "from@domain.ch",
					FromName:    "name",
					SMTP: smtp.SMTP{
						Host:     "host:587",
						User:     "user",
						Password: "password",
					},
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
				eventstore: eventstoreExpect(
					t,
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
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Description:    "test",
					Tls:            true,
					From:           "from@domain.ch",
					FromName:       "name",
					ReplyToAddress: "replyto@domain.ch",
					SMTP: smtp.SMTP{
						Host:     "host:587",
						User:     "user",
						Password: "password",
					},
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
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "configid"),
				alg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Description: "test",
					Tls:         true,
					From:        "from@domain.ch",
					FromName:    "name",
					SMTP: smtp.SMTP{
						Host:     "host",
						User:     "user",
						Password: "password",
					},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "smtp config, host is empty",
			fields: fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "configid"),
				alg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Description: "test",
					Tls:         true,
					From:        "from@domain.ch",
					FromName:    "name",
					SMTP: smtp.SMTP{
						Host:     "   ",
						User:     "user",
						Password: "password",
					},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "add smtp config, ipv6 works",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Description: "test",
					Tls:         true,
					From:        "from@domain.ch",
					FromName:    "name",
					SMTP: smtp.SMTP{
						Host:     "[2001:db8::1]:2525",
						User:     "user",
						Password: "password",
					},
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
				eventstore:     tt.fields.eventstore,
				smtpEncryption: tt.fields.alg,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			_, got, err := r.AddSMTPConfig(tt.args.ctx, tt.args.instanceID, tt.args.smtp)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ChangeSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		instanceID string
		id         string
		smtp       *smtp.Config
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
			name: "id empty, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "empty config, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{},
				id:   "configID",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "smtp not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Description: "test",
					Tls:         true,
					From:        "from@domain.ch",
					FromName:    "name",
					SMTP: smtp.SMTP{
						Host: "host:587",
						User: "user",
					},
				},
				instanceID: "INSTANCE",
				id:         "ID",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "smtp domain not matched",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				instanceID: "INSTANCE",
				id:         "ID",
				smtp: &smtp.Config{
					Description: "test",
					Tls:         true,
					From:        "from@wrongdomain.ch",
					FromName:    "name",
					SMTP: smtp.SMTP{
						Host: "host:587",
						User: "user",
					},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Description: "test",
					Tls:         true,
					From:        "from@domain.ch",
					FromName:    "name",
					SMTP: smtp.SMTP{
						Host: "host:587",
						User: "user",
					},
				},
				instanceID: "INSTANCE",
				id:         "ID",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "smtp config change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Description:    "test",
					Tls:            false,
					From:           "from2@domain.ch",
					FromName:       "name2",
					ReplyToAddress: "replyto@domain.ch",
					SMTP: smtp.SMTP{
						Host: "host2:587",
						User: "user2",
					},
				},
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
			name: "smtp config, port is missing",
			fields: fields{
				eventstore: eventstoreExpect(t),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Description: "test",
					Tls:         true,
					From:        "from@domain.ch",
					FromName:    "name",
					SMTP: smtp.SMTP{
						Host:     "host",
						User:     "user",
						Password: "password",
					},
				},
				instanceID: "INSTANCE",
				id:         "ID",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "smtp config, host is empty",
			fields: fields{
				eventstore: eventstoreExpect(t),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Description: "test",
					Tls:         true,
					From:        "from@domain.ch",
					FromName:    "name",
					SMTP: smtp.SMTP{
						Host:     "   ",
						User:     "user",
						Password: "password",
					},
				},
				instanceID: "INSTANCE",
				id:         "ID",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "smtp config change, ipv6 works",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Description:    "test",
					Tls:            false,
					From:           "from2@domain.ch",
					FromName:       "name2",
					ReplyToAddress: "replyto@domain.ch",
					SMTP: smtp.SMTP{
						Host: "[2001:db8::1]:2525",
						User: "user2",
					},
				},
				instanceID: "INSTANCE",
				id:         "ID",
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
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeSMTPConfig(tt.args.ctx, tt.args.instanceID, tt.args.id, tt.args.smtp)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ChangeSMTPConfigPassword(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		alg        crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx        context.Context
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
			name: "smtp config, error not found",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:      context.Background(),
				password: "",
				id:       "ID",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "change smtp config password, ok",
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
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
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
				eventstore:     tt.fields.eventstore,
				smtpEncryption: tt.fields.alg,
			}
			got, err := r.ChangeSMTPConfigPassword(tt.args.ctx, tt.args.instanceID, tt.args.id, tt.args.password)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ActivateSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		alg        crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx         context.Context
		instanceID  string
		id          string
		activatedId string
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
			name: "smtp not existing, not found error",
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
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "activate smtp config, ok",
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
				ctx:         authz.WithInstanceID(context.Background(), "INSTANCE"),
				id:          "ID",
				instanceID:  "INSTANCE",
				activatedId: "",
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
				eventstore:     tt.fields.eventstore,
				smtpEncryption: tt.fields.alg,
			}
			got, err := r.ActivateSMTPConfig(tt.args.ctx, tt.args.instanceID, tt.args.id, tt.args.activatedId)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_DeactivateSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		alg        crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx         context.Context
		instanceID  string
		id          string
		activatedId string
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
			name: "smtp not existing, not found error",
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
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "deactivate smtp config, ok",
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
				ctx:         authz.WithInstanceID(context.Background(), "INSTANCE"),
				id:          "ID",
				instanceID:  "INSTANCE",
				activatedId: "",
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
				eventstore:     tt.fields.eventstore,
				smtpEncryption: tt.fields.alg,
			}
			got, err := r.DeactivateSMTPConfig(tt.args.ctx, tt.args.instanceID, tt.args.id)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_RemoveSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
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
	}{
		{
			name: "smtp config, error not found",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				id:  "ID",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "remove smtp config, ok",
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
				eventstore:     tt.fields.eventstore,
				smtpEncryption: tt.fields.alg,
			}
			got, err := r.RemoveSMTPConfig(tt.args.ctx, tt.args.instanceID, tt.args.id)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_TestSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
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
			name: "if password is empty, smtp id must not",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
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
				eventstore: eventstoreExpect(
					t,
				),
				alg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
				eventstore:     tt.fields.eventstore,
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
