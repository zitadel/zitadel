package user

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func Test_patchMachineUserToCommand(t *testing.T) {
	type args struct {
		userId   string
		userName *string
		machine  *user.UpdateUserRequest_Machine
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
				Name: ptr("name"),
			},
		},
		want: &command.ChangeMachine{
			ID:   "userId",
			Name: ptr("name"),
		},
	}, {
		name: "all properties",
		args: args{
			userId:   "userId",
			userName: ptr("userName"),
			machine: &user.UpdateUserRequest_Machine{
				Name:        ptr("name"),
				Description: ptr("description"),
			},
		},
		want: &command.ChangeMachine{
			ID:          "userId",
			Username:    ptr("userName"),
			Name:        ptr("name"),
			Description: ptr("description"),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := patchMachineUserToCommand(tt.args.userId, tt.args.userName, tt.args.machine)
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateComparable(language.Tag{})); diff != "" {
				t.Errorf("patchMachineUserToCommand() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
