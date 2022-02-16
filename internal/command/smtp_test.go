package command

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/notification/channels/smtp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/iam"
)

func TestCommandSide_AddSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		alg        crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx  context.Context
		smtp *smtp.EmailConfig
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
			name: "smtp config, error already exists",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewSMTPConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								true,
								"from",
								"name",
								"host",
								"user",
								&crypto.CryptoValue{},
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				smtp: &smtp.EmailConfig{
					Tls: true,
				},
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "add smtp config, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(iam.NewSMTPConfigAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								true,
								"from",
								"name",
								"host",
								"user",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("password"),
								},
							),
							),
						},
					),
				),
				alg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				smtp: &smtp.EmailConfig{
					Tls:      true,
					From:     "from",
					FromName: "name",
					SMTP: smtp.SMTP{
						Host:     "host",
						User:     "user",
						Password: "password",
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:         tt.fields.eventstore,
				smtpPasswordCrypto: tt.fields.alg,
			}
			got, err := r.AddSMTPConfig(tt.args.ctx, tt.args.smtp)
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
		ctx  context.Context
		smtp *smtp.EmailConfig
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
			name: "smtp not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:  context.Background(),
				smtp: &smtp.EmailConfig{},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},

		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewSMTPConfigAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								true,
								"from",
								"name",
								"host",
								"user",
								&crypto.CryptoValue{},
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				smtp: &smtp.EmailConfig{
					Tls:      true,
					From:     "from",
					FromName: "name",
					SMTP: smtp.SMTP{
						Host: "host",
						User: "user",
					},
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "smtp config change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewSMTPConfigAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								true,
								"from",
								"name",
								"host",
								"user",
								&crypto.CryptoValue{},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newSMTPConfigChangedEvent(
									context.Background(),
									false,
									"from2",
									"name2",
									"host2",
									"user2",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				smtp: &smtp.EmailConfig{
					Tls:      false,
					From:     "from2",
					FromName: "name2",
					SMTP: smtp.SMTP{
						Host: "host2",
						User: "user2",
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeSMTPConfig(tt.args.ctx, tt.args.smtp)
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
		ctx      context.Context
		password string
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
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "change smtp config password, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewSMTPConfigAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								true,
								"from",
								"name",
								"host",
								"user",
								&crypto.CryptoValue{},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(iam.NewSMTPConfigPasswordChangedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("password"),
								},
							),
							),
						},
					),
				),
				alg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:      context.Background(),
				password: "password",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:         tt.fields.eventstore,
				smtpPasswordCrypto: tt.fields.alg,
			}
			got, err := r.ChangeSMTPConfigPassword(tt.args.ctx, tt.args.password)
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

func newSMTPConfigChangedEvent(ctx context.Context, tls bool, fromAddress, fromName, host, user string) *iam.SMTPConfigChangedEvent {
	changes := []iam.SMTPConfigChanges{
		iam.ChangeSMTPConfigTLS(tls),
		iam.ChangeSMTPConfigFromAddress(fromAddress),
		iam.ChangeSMTPConfigFromName(fromName),
		iam.ChangeSMTPConfigSMTPHost(host),
		iam.ChangeSMTPConfigSMTPUser(user),
	}
	event, _ := iam.NewSMTPConfigChangeEvent(ctx,
		&iam.NewAggregate().Aggregate,
		changes,
	)
	return event
}
