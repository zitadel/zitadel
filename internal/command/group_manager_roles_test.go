package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	repo "github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_SetGroupManagerRoles(t *testing.T) {
	t.Parallel()

	zitadelRoles := []authz.RoleMapping{
		{Role: "ORG_OWNER"},
		{Role: "ORG_OWNER_VIEWER"},
		// IAM_OWNER is intentionally a known, valid system role here. The
		// command must still reject it because group manager roles are
		// scoped to the org level by domain.OrgRolePrefix.
		{Role: domain.RoleIAMOwner},
	}

	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx     context.Context
		groupID string
		roles   []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr func(error) bool
	}{
		{
			name: "missing group id, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:   context.Background(),
				roles: []string{"ORG_OWNER_VIEWER"},
			},
			wantErr: zerrors.IsErrorInvalidArgument,
		},
		{
			name: "invalid role, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:     context.Background(),
				groupID: "group1",
				roles:   []string{"NOT_A_ROLE"},
			},
			wantErr: zerrors.IsErrorInvalidArgument,
		},
		{
			// IAM_OWNER is in the zitadelRoles registry but lacks the ORG_
			// prefix. The command must reject it before reaching the
			// eventstore — without the prefix guard the role would escalate
			// every group member to an instance-wide IAM_OWNER.
			name: "IAM_OWNER role rejected, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:     context.Background(),
				groupID: "group1",
				roles:   []string{domain.RoleIAMOwner},
			},
			wantErr: zerrors.IsErrorInvalidArgument,
		},
		{
			name: "group not found, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:     context.Background(),
				groupID: "group1",
				roles:   []string{"ORG_OWNER_VIEWER"},
			},
			wantErr: zerrors.IsPreconditionFailed,
		},
		{
			name: "missing permission, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(addNewGroupEvent("group1", "org1")),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:     context.Background(),
				groupID: "group1",
				roles:   []string{"ORG_OWNER_VIEWER"},
			},
			wantErr: zerrors.IsPermissionDenied,
		},
		{
			name: "roles unchanged, no events pushed, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(addNewGroupEvent("group1", "org1")),
					),
					expectFilter( // existing manager roles
						eventFromEventPusher(
							repo.NewGroupManagerRolesSetEvent(context.Background(),
								&repo.NewAggregate("group1", "org1").Aggregate,
								[]string{"ORG_OWNER_VIEWER"},
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:     context.Background(),
				groupID: "group1",
				roles:   []string{"ORG_OWNER_VIEWER"},
			},
			want: &domain.ObjectDetails{
				ID:            "group1",
				ResourceOwner: "org1",
			},
		},
		{
			name: "manager roles set, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(addNewGroupEvent("group1", "org1")),
					),
					expectFilter(), // no existing manager roles
					expectPush(
						repo.NewGroupManagerRolesSetEvent(context.Background(),
							&repo.NewAggregate("group1", "org1").Aggregate,
							[]string{"ORG_OWNER_VIEWER"},
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:     context.Background(),
				groupID: "group1",
				roles:   []string{"ORG_OWNER_VIEWER"},
			},
			want: &domain.ObjectDetails{
				ID:            "group1",
				ResourceOwner: "org1",
			},
		},
		{
			name: "manager roles removed with empty list, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(addNewGroupEvent("group1", "org1")),
					),
					expectFilter( // existing manager roles
						eventFromEventPusher(
							repo.NewGroupManagerRolesSetEvent(context.Background(),
								&repo.NewAggregate("group1", "org1").Aggregate,
								[]string{"ORG_OWNER_VIEWER"},
							),
						),
					),
					expectPush(
						repo.NewGroupManagerRolesSetEvent(context.Background(),
							&repo.NewAggregate("group1", "org1").Aggregate,
							nil,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:     context.Background(),
				groupID: "group1",
			},
			want: &domain.ObjectDetails{
				ID:            "group1",
				ResourceOwner: "org1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
				zitadelRoles:    zitadelRoles,
			}
			got, err := c.SetGroupManagerRoles(tt.args.ctx, tt.args.groupID, tt.args.roles)
			if tt.wantErr != nil {
				require.True(t, tt.wantErr(err))
				return
			}
			require.NoError(t, err)
			assertObjectDetails(t, tt.want, got)
		})
	}
}
