package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func TestCommandSide_AddSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		alg        crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx  context.Context
		smtp *smtp.Config
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
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Tls:      true,
					From:     "from@domain.ch",
					FromName: "name",
					SMTP: smtp.SMTP{
						Host:     "host:587",
						User:     "user",
						Password: "password",
					},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "smtp config, error already exists",
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
						eventFromEventPusher(
							instance.NewSMTPConfigAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								"from@domain.ch",
								"name",
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
					Tls:      true,
					From:     "from@domain.ch",
					FromName: "name",
					SMTP: smtp.SMTP{
						Host:     "host:587",
						User:     "user",
						Password: "password",
					},
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"INSTANCE",
								instance.NewSMTPConfigAddedEvent(
									context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									true,
									"from@domain.ch",
									"name",
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
						},
					),
				),
				alg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Tls:      true,
					From:     "from@domain.ch",
					FromName: "name",
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
				eventstore: eventstoreExpect(t),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Tls:      true,
					From:     "from@domain.ch",
					FromName: "name",
					SMTP: smtp.SMTP{
						Host:     "host",
						User:     "user",
						Password: "password",
					},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
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
					Tls:      true,
					From:     "from@domain.ch",
					FromName: "name",
					SMTP: smtp.SMTP{
						Host:     "   ",
						User:     "user",
						Password: "password",
					},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"INSTANCE",
								instance.NewSMTPConfigAddedEvent(
									context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									true,
									"from@domain.ch",
									"name",
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
						},
					),
				),
				alg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Tls:      true,
					From:     "from@domain.ch",
					FromName: "name",
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
		smtp *smtp.Config
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
			name: "empty config, invalid argument error",
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
				err: caos_errs.IsErrorInvalidArgument,
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
					Tls:      true,
					From:     "from@domain.ch",
					FromName: "name",
					SMTP: smtp.SMTP{
						Host: "host:587",
						User: "user",
					},
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
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
								true,
								"from@domain.ch",
								"name",
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
					Tls:      true,
					From:     "from@wrongdomain.ch",
					FromName: "name",
					SMTP: smtp.SMTP{
						Host: "host:587",
						User: "user",
					},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
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
								true,
								"from@domain.ch",
								"name",
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
					Tls:      true,
					From:     "from@domain.ch",
					FromName: "name",
					SMTP: smtp.SMTP{
						Host: "host:587",
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
								true,
								"from@domain.ch",
								"name",
								"host:587",
								"user",
								&crypto.CryptoValue{},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"INSTANCE",
								newSMTPConfigChangedEvent(
									context.Background(),
									false,
									"from2@domain.ch",
									"name2",
									"host2:587",
									"user2",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Tls:      false,
					From:     "from2@domain.ch",
					FromName: "name2",
					SMTP: smtp.SMTP{
						Host: "host2:587",
						User: "user2",
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
				eventstore: eventstoreExpect(t),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Tls:      true,
					From:     "from@domain.ch",
					FromName: "name",
					SMTP: smtp.SMTP{
						Host:     "host",
						User:     "user",
						Password: "password",
					},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
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
					Tls:      true,
					From:     "from@domain.ch",
					FromName: "name",
					SMTP: smtp.SMTP{
						Host:     "   ",
						User:     "user",
						Password: "password",
					},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
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
								true,
								"from@domain.ch",
								"name",
								"host:587",
								"user",
								&crypto.CryptoValue{},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"INSTANCE",
								newSMTPConfigChangedEvent(
									context.Background(),
									false,
									"from2@domain.ch",
									"name2",
									"[2001:db8::1]:2525",
									"user2",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				smtp: &smtp.Config{
					Tls:      false,
					From:     "from2@domain.ch",
					FromName: "name2",
					SMTP: smtp.SMTP{
						Host: "[2001:db8::1]:2525",
						User: "user2",
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
							instance.NewSMTPConfigAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								"from",
								"name",
								"host:587",
								"user",
								&crypto.CryptoValue{},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"INSTANCE",
								instance.NewSMTPConfigPasswordChangedEvent(
									context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
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
				ctx:      authz.WithInstanceID(context.Background(), "INSTANCE"),
				password: "password",
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

func TestCommandSide_RemoveSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		alg        crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx context.Context
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
			},
			res: res{
				err: caos_errs.IsNotFound,
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
								true,
								"from",
								"name",
								"host:587",
								"user",
								&crypto.CryptoValue{},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"INSTANCE",
								instance.NewSMTPConfigRemovedEvent(
									context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
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
			got, err := r.RemoveSMTPConfig(tt.args.ctx)
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

func newSMTPConfigChangedEvent(ctx context.Context, tls bool, fromAddress, fromName, host, user string) *instance.SMTPConfigChangedEvent {
	changes := []instance.SMTPConfigChanges{
		instance.ChangeSMTPConfigTLS(tls),
		instance.ChangeSMTPConfigFromAddress(fromAddress),
		instance.ChangeSMTPConfigFromName(fromName),
		instance.ChangeSMTPConfigSMTPHost(host),
		instance.ChangeSMTPConfigSMTPUser(user),
	}
	event, _ := instance.NewSMTPConfigChangeEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		changes,
	)
	return event
}
