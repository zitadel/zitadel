package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_RemoveOrgSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		orgID string
		id    string
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-t2WsP1", "Errors.ResourceOwnerMissing"))
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-0ZV5w1", "Errors.IDMissing"))
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
				orgID: "ORG1",
				id:    "ID",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "ORG-09CXl1", "Errors.SMTPConfig.NotFound"))
				},
			},
		},
		{
			name: "remove smtp config, ok",
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
					expectPush(
						org.NewOrgSMTPConfigRemovedEvent(
							context.Background(),
							&org.NewAggregate("ORG1").Aggregate,
							"ID",
						),
					),
				),
			},
			args: args{
				orgID: "ORG1",
				id:    "ID",
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.RemoveOrgSMTPConfig(context.Background(), tt.args.orgID, tt.args.id)
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

func TestCommandSide_ActivateOrgSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		orgID string
		id    string
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
				orgID: "",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-h5htM1", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "id empty, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				orgID: "ORG1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-1hPl61", "Errors.IDMissing"))
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
				orgID: "ORG1",
				id:    "ID",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "ORG-E9K201", "Errors.SMTPConfig.NotFound"))
				},
			},
		},
		{
			name: "activate smtp config, ok",
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
					expectPush(
						org.NewOrgSMTPConfigActivatedEvent(
							context.Background(),
							&org.NewAggregate("ORG1").Aggregate,
							"ID",
						),
					),
				),
			},
			args: args{
				orgID: "ORG1",
				id:    "ID",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "ORG1",
				},
			},
		},
		{
			name: "activate smtp config, already active",
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
				),
			},
			args: args{
				orgID: "ORG1",
				id:    "ID",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "ORG-vUHBS1", "Errors.SMTPConfig.AlreadyActive"))
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.ActivateOrgSMTPConfig(context.Background(), tt.args.orgID, tt.args.id)
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

func TestCommandSide_DeactivateOrgSMTPConfig(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		orgID string
		id    string
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
			args: args{},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-pvNHo1", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "id empty, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				orgID: "ORG1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-jLTIM1", "Errors.IDMissing"))
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
				orgID: "ORG1",
				id:    "ID",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "ORG-k39PJ1", "Errors.SMTPConfig.NotFound"))
				},
			},
		},
		{
			name: "deactivate smtp config, ok",
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
						org.NewOrgSMTPConfigDeactivatedEvent(
							context.Background(),
							&org.NewAggregate("ORG1").Aggregate,
							"ID",
						),
					),
				),
			},
			args: args{
				orgID: "ORG1",
				id:    "ID",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "ORG1",
				},
			},
		},
		{
			name: "deactivate smtp config, already deactivated",
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
						eventFromEventPusher(
							org.NewOrgSMTPConfigDeactivatedEvent(
								context.Background(),
								&org.NewAggregate("ORG1").Aggregate,
								"ID",
							),
						),
					),
				),
			},
			args: args{
				orgID: "ORG1",
				id:    "ID",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "ORG-km8g31", "Errors.SMTPConfig.AlreadyDeactivated"))
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.DeactivateOrgSMTPConfig(context.Background(), tt.args.orgID, tt.args.id)
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
