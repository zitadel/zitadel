package key

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checkSingleFlag(t *testing.T) {
	type args struct {
		masterKeyFile    string
		masterKeyFromArg string
		masterKeyFromEnv bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"no values, error",
			args{
				masterKeyFile:    "",
				masterKeyFromArg: "",
				masterKeyFromEnv: false,
			},
			assert.Error,
		},
		{
			"multiple values, error",
			args{
				masterKeyFile:    "file",
				masterKeyFromArg: "masterkey",
				masterKeyFromEnv: true,
			},
			assert.Error,
		},
		{
			"only file, ok",
			args{
				masterKeyFile:    "file",
				masterKeyFromArg: "",
				masterKeyFromEnv: false,
			},
			assert.NoError,
		},
		{
			"only argument, ok",
			args{
				masterKeyFile:    "",
				masterKeyFromArg: "masterkey",
				masterKeyFromEnv: false,
			},
			assert.NoError,
		},
		{
			"only env, ok",
			args{
				masterKeyFile:    "",
				masterKeyFromArg: "",
				masterKeyFromEnv: true,
			},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, checkSingleFlag(tt.args.masterKeyFile, tt.args.masterKeyFromArg, tt.args.masterKeyFromEnv), fmt.Sprintf("checkSingleFlag(%v, %v)", tt.args.masterKeyFile, tt.args.masterKeyFromArg))
		})
	}
}
