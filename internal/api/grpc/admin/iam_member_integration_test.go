//go:build integration

package admin_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/integration"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/object"
)

var iamRoles = []string{
	"IAM_OWNER",
	"IAM_OWNER_VIEWER",
	"IAM_ORG_MANAGER",
	"IAM_USER_MANAGER",
	"ADMIN_IMPERSONATOR",
	"END_USER_IMPERSONATOR",
}

func TestServer_AddIAMMember(t *testing.T) {
	user := Tester.CreateHumanUserVerified(AdminCTX, Tester.Organisation.ID, gofakeit.Email())

	type args struct {
		ctx context.Context
		req *admin_pb.AddIAMMemberRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *admin_pb.AddIAMMemberResponse
		wantErr bool
	}{
		{
			name: "permission error",
			args: args{
				ctx: Tester.WithAuthorization(CTX, integration.OrgOwner),
				req: &admin_pb.AddIAMMemberRequest{
					UserId: user.GetUserId(),
					Roles:  iamRoles,
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: AdminCTX,
				req: &admin_pb.AddIAMMemberRequest{
					UserId: user.GetUserId(),
					Roles:  iamRoles,
				},
			},
			want: &admin_pb.AddIAMMemberResponse{
				Details: &object.ObjectDetails{
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tester.Client.Admin.AddIAMMember(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}
