package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_ChangeOrgSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
		alg        crypto.EncryptionAlgorithm
	}
	type args struct {
		smtp *ChangeOrgSMTPConfig
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
				smtp: &ChangeOrgSMTPConfig{},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-wJk3s0", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "id empty, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				smtp: &ChangeOrgSMTPConfig{
					ResourceOwner: "ORG1",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-2MHqR1", "Errors.IDMissing"))
				},
			},
		},
		{
			name: "from empty, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				smtp: &ChangeOrgSMTPConfig{
					ResourceOwner: "ORG1",
					ID:            "configID",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-g9PXN4", "Errors.Invalid.Argument"))
				},
			},
		},
		{
			name: "host missing port, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				smtp: &ChangeOrgSMTPConfig{
					ResourceOwner: "ORG1",
					ID:            "configID",
					From:          "from@example.com",
					Host:          "host",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-kZ3Vk2", "Errors.Invalid.Argument"))
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
				smtp: &ChangeOrgSMTPConfig{
					ResourceOwner: "ORG1",
					ID:            "ID",
					From:          "from@example.com",
					Host:          "host:587",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "ORG-j5IDn3", "Errors.SMTPConfig.NotFound"))
				},
			},
		},
		{
			name: "no changes, returns details",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOrgSMTPConfigAddedEvent(
								context.Background(),
								&org.NewAggregate("ORG1").Aggregate,
								"ID",
								"test",
								true,
								"from@example.com",
								"name",
								"",
								"host:587",
								"user",
								&instance.PlainAuth{
									Password: &crypto.CryptoValue{},
								},
								nil,
							),
						),
					),
				),
			},
			args: args{
				smtp: &ChangeOrgSMTPConfig{
					ResourceOwner: "ORG1",
					ID:            "ID",
					Description:   "test",
					Tls:           true,
					From:          "from@example.com",
					FromName:      "name",
					Host:          "host:587",
					User:          "user",
					PlainAuth:     &PlainAuth{},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "ORG1",
				},
			},
		},
		{
			name: "smtp config change, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOrgSMTPConfigAddedEvent(
								context.Background(),
								&org.NewAggregate("ORG1").Aggregate,
								"ID",
								"",
								true,
								"from@example.com",
								"name",
								"",
								"host:587",
								"user",
								&instance.PlainAuth{
									Password: &crypto.CryptoValue{},
								},
								nil,
							),
						),
					),
					expectPush(
						newOrgSMTPConfigChangedEvent(
							context.Background(),
							"ORG1",
							"ID",
							"test",
							false,
							"from2@example.com",
							"name2",
							"replyto@example.com",
							"host2:587",
							"user2",
							&instance.PlainAuth{},
							nil,
						),
					),
				),
			},
			args: args{
				smtp: &ChangeOrgSMTPConfig{
					ResourceOwner:  "ORG1",
					ID:             "ID",
					Description:    "test",
					Tls:            false,
					From:           "from2@example.com",
					FromName:       "name2",
					ReplyToAddress: "replyto@example.com",
					Host:           "host2:587",
					User:           "user2",
					PlainAuth:      &PlainAuth{},
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
				smtpEncryption: tt.fields.alg,
			}
			err := r.ChangeOrgSMTPConfig(context.Background(), tt.args.smtp)
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

func TestCommandSide_ChangeOrgSMTPConfigPassword(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
		alg        crypto.EncryptionAlgorithm
	}
	type args struct {
		orgID    string
		id       string
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
			name: "resourceowner empty, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-gHAyv1", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "id empty, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				orgID: "ORG1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-BCkAf1", "Errors.IDMissing"))
				},
			},
		},
		{
			name: "smtp config not active, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				orgID: "ORG1",
				id:    "ID",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "ORG-rDHzq1", "Errors.SMTPConfig.NotFound"))
				},
			},
		},
		{
			name: "change smtp config password, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOrgSMTPConfigAddedEvent(
								context.Background(),
								&org.NewAggregate("ORG1").Aggregate,
								"ID",
								"test",
								true,
								"from@example.com",
								"name",
								"",
								"host:587",
								"user",
								&instance.PlainAuth{
									Password: &crypto.CryptoValue{},
								},
								nil,
							),
						),
						eventFromEventPusher(
							org.NewOrgSMTPConfigActivatedEvent(
								context.Background(),
								&org.NewAggregate("ORG1").Aggregate,
								"ID",
							),
						),
					),
					expectPush(
						org.NewOrgSMTPConfigPasswordChangedEvent(
							context.Background(),
							&org.NewAggregate("ORG1").Aggregate,
							"ID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("newpassword"),
							},
						),
					),
				),
				alg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				orgID:    "ORG1",
				id:       "ID",
				password: "newpassword",
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
				smtpEncryption: tt.fields.alg,
			}
			got, err := r.ChangeOrgSMTPConfigPassword(context.Background(), tt.args.orgID, tt.args.id, tt.args.password)
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
