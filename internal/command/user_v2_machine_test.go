package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_ChangeUserMachine(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx     context.Context
		orgID   string
		machine *ChangeMachine
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}

	userAgg := user.NewAggregate("user1", "org1")
	userAddedEvent := user.NewMachineAddedEvent(context.Background(),
		&userAgg.Aggregate,
		"username",
		"name",
		"description",
		true,
		domain.OIDCTokenTypeBearer,
	)

	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "change machine username, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(userAddedEvent),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				machine: &ChangeMachine{
					Username: gu.Ptr("changed"),
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "change machine username, not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				machine: &ChangeMachine{
					Username: gu.Ptr("changed"),
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-ugjs0upun6", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "change machine username, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(userAddedEvent),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						user.NewUsernameChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"username",
							"changed",
							true,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				machine: &ChangeMachine{
					Username: gu.Ptr("changed"),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change machine username, no change",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(userAddedEvent),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				machine: &ChangeMachine{
					Username: gu.Ptr("username"),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change machine description, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(userAddedEvent),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				machine: &ChangeMachine{
					Description: gu.Ptr("changed"),
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "change machine description, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(userAddedEvent),
					),
					expectPush(
						user.NewMachineChangedEvent(context.Background(),
							&userAgg.Aggregate,
							[]user.MachineChanges{
								user.ChangeDescription("changed"),
							},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				machine: &ChangeMachine{
					Description: gu.Ptr("changed"),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change machine description, no change",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(userAddedEvent),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				machine: &ChangeMachine{
					Description: gu.Ptr("description"),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			err := r.ChangeUserMachine(tt.args.ctx, tt.args.machine)
			if tt.res.err == nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, tt.args.machine.Details)
			}
		})
	}
}
