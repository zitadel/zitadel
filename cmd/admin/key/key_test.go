package key

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errors "github.com/caos/zitadel/internal/errors"

	"github.com/caos/zitadel/internal/crypto"
)

func Test_keysFromArgs(t *testing.T) {
	type args struct {
		args []string
	}
	type res struct {
		keys []*crypto.Key
		err  func(error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"no args",
			args{},
			res{
				keys: []*crypto.Key{},
			},
		},
		{
			"invalid arg",
			args{
				args: []string{"keyID", "value"},
			},
			res{
				err: caos_errors.IsInternal,
			},
		},
		{
			"single arg",
			args{
				args: []string{"keyID=value"},
			},
			res{
				keys: []*crypto.Key{
					{
						ID:    "keyID",
						Value: "value",
					},
				},
			},
		},
		{
			"multiple args",
			args{
				args: []string{"keyID=value", "keyID2=value2"},
			},
			res{
				keys: []*crypto.Key{
					{
						ID:    "keyID",
						Value: "value",
					},
					{
						ID:    "keyID2",
						Value: "value2",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := keysFromArgs(tt.args.args)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if !reflect.DeepEqual(got, tt.res.keys) {
				t.Errorf("keysFromArgs() got = %v, want %v", got, tt.res.keys)
			}
		})
	}
}

func Test_keysFromYAML(t *testing.T) {
	type args struct {
		file io.Reader
	}
	type res struct {
		keys []*crypto.Key
		err  func(error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"invalid yaml",
			args{
				file: bytes.NewReader([]byte("keyID=ds")),
			},
			res{
				err: caos_errors.IsInternal,
			},
		},
		{
			"single key",
			args{
				file: bytes.NewReader([]byte("keyID: value")),
			},
			res{
				keys: []*crypto.Key{
					{
						ID:    "keyID",
						Value: "value",
					},
				},
			},
		},
		{
			"multiple keys",
			args{
				file: bytes.NewReader([]byte("keyID: value\nkeyID2: value2")),
			},
			res{
				keys: []*crypto.Key{
					{
						ID:    "keyID",
						Value: "value",
					},
					{
						ID:    "keyID2",
						Value: "value2",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := keysFromYAML(tt.args.file)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			assert.ElementsMatch(t, got, tt.res.keys)
		})
	}
}
