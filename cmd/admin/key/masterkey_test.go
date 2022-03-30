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
			},
			assert.Error,
		},
		{
			"both values, error",
			args{
				masterKeyFile:    "file",
				masterKeyFromArg: "masterkey",
			},
			assert.Error,
		},
		{
			"only file, ok",
			args{
				masterKeyFile:    "file",
				masterKeyFromArg: "",
			},
			assert.NoError,
		},
		{
			"only argument, ok",
			args{
				masterKeyFile:    "",
				masterKeyFromArg: "masterkey",
			},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, checkSingleFlag(tt.args.masterKeyFile, tt.args.masterKeyFromArg), fmt.Sprintf("checkSingleFlag(%v, %v)", tt.args.masterKeyFile, tt.args.masterKeyFromArg))
		})
	}
}
