package authz

import (
	"context"
	"testing"

	"github.com/zitadel/zitadel/internal/errors"
)

func Test_extractBearerToken(t *testing.T) {

	type args struct {
		ctx      context.Context
		token    string
		verifier AccessTokenVerifier
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
				verifier: AccessTokenVerifierFunc(func(context.Context, string) (string, string, string, string, string, error) {
					return "", "", "", "", "", nil
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := extractBearerToken(tt.args.token)
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
