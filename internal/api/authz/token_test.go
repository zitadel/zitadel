package authz

import (
	"context"
	"sync"
	"testing"

	"github.com/caos/zitadel/internal/errors"
)

func Test_VerifyAccessToken(t *testing.T) {

	type args struct {
		ctx      context.Context
		token    string
		verifier *TokenVerifier
		method   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no auth header set",
			args: args{
				ctx:   context.Background(),
				token: "",
			},
			wantErr: true,
		},
		{
			name: "wrong auth header set",
			args: args{
				ctx:   context.Background(),
				token: "Basic sds",
			},
			wantErr: true,
		},
		{
			name: "auth header set",
			args: args{
				ctx:   context.Background(),
				token: "Bearer AUTH",
				verifier: &TokenVerifier{
					authZRepo: &testVerifier{memberships: []*Membership{}},
					clients: func() sync.Map {
						m := sync.Map{}
						m.Store("service", &client{name: "name"})
						return m
					}(),
					authMethods: MethodMapping{"/service/method": Option{Permission: "authenticated"}},
				},
				method: "/service/method",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, _, err := verifyAccessToken(tt.args.ctx, tt.args.token, tt.args.verifier, tt.args.method)
			if tt.wantErr && err == nil {
				t.Errorf("got wrong result, should get err: actual: %v ", err)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("got wrong result, should not get err: actual: %v ", err)
			}

			if tt.wantErr && !errors.IsUnauthenticated(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
