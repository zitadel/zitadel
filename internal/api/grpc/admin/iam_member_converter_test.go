package admin

import (
	"testing"

	"github.com/zitadel/zitadel/internal/test"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

func TestAddIAMMemberToDomain(t *testing.T) {
	type args struct {
		req *admin.AddIAMMemberRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "all fields filled",
			args: args{
				req: &admin.AddIAMMemberRequest{
					UserId: "1232452",
					Roles:  []string{"admin"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddIAMMemberToCommand(tt.args.req, "INSTANCE")
			test.AssertFieldsMapped(t, got, "ObjectRoot")
		})
	}
}

func TestUpdateIAMMemberToDomain(t *testing.T) {
	type args struct {
		req *admin.UpdateIAMMemberRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "all fields filled",
			args: args{
				req: &admin.UpdateIAMMemberRequest{
					UserId: "1232452",
					Roles:  []string{"admin"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UpdateIAMMemberToCommand(tt.args.req, "INSTANCE")
			test.AssertFieldsMapped(t, got, "ObjectRoot")
		})
	}
}
