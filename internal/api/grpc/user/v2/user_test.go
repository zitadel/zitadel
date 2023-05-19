package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func Test_hashedPasswordToCommand(t *testing.T) {
	type args struct {
		hashed *user.HashedPassword
	}
	type res struct {
		want string
		err  func(error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"not hashed",
			args{
				hashed: nil,
			},
			res{
				"",
				nil,
			},
		},
		{
			"hashed, not bcrypt",
			args{
				hashed: &user.HashedPassword{
					Hash:      "hash",
					Algorithm: "custom",
				},
			},
			res{
				"",
				func(err error) bool {
					return errors.Is(err, caos_errs.ThrowInvalidArgument(nil, "USER-JDk4t", "Errors.InvalidArgument"))
				},
			},
		},
		{
			"hashed, bcrypt",
			args{
				hashed: &user.HashedPassword{
					Hash:      "hash",
					Algorithm: "bcrypt",
				},
			},
			res{
				"hash",
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hashedPasswordToCommand(tt.args.hashed)
			if tt.res.err == nil {
				require.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}
