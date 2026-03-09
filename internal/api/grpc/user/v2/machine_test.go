package user

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/muhlemmer/gu"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func Test_patchMachineUserToCommand(t *testing.T) {
	type args struct {
		userId   string
		userName *string
		machine  *user.UpdateUserRequest_Machine
		metadata []*user.Metadata
	}
	tests := []struct {
		name string
		args args
		want *command.ChangeMachine
	}{{
		name: "single property",
		args: args{
			userId: "userId",
			machine: &user.UpdateUserRequest_Machine{
				Name: gu.Ptr("name"),
			},
		},
		want: &command.ChangeMachine{
			ID:   "userId",
			Name: gu.Ptr("name"),
		},
	}, {
		name: "all properties",
		args: args{
			userId:   "userId",
			userName: gu.Ptr("userName"),
			machine: &user.UpdateUserRequest_Machine{
				Name:            gu.Ptr("name"),
				Description:     gu.Ptr("description"),
				AccessTokenType: gu.Ptr(user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT),
			},
			metadata: []*user.Metadata{
				{
					Key:   "key1",
					Value: []byte("value1"),
				},
				{
					Key:   "key2",
					Value: []byte("value2"),
				},
			},
		},
		want: &command.ChangeMachine{
			ID:              "userId",
			Username:        gu.Ptr("userName"),
			Name:            gu.Ptr("name"),
			Description:     gu.Ptr("description"),
			AccessTokenType: gu.Ptr(domain.OIDCTokenTypeJWT),
			Metadata: []*domain.Metadata{
				{
					Key:   "key1",
					Value: []byte("value1"),
				},
				{
					Key:   "key2",
					Value: []byte("value2"),
				},
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateMachineUserToCommand(tt.args.userId, tt.args.userName, tt.args.machine, tt.args.metadata)
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateComparable(language.Tag{})); diff != "" {
				t.Errorf("patchMachineUserToCommand() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
