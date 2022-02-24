package key

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/cmd/helper"
	"github.com/caos/zitadel/internal/crypto"
)

func Test_keysFromArgs(t *testing.T) {
	type args struct {
		args []string
	}
	type res struct {
		keys []crypto.Key
		err  error
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
				keys: []crypto.Key{},
			},
		},
		{
			"invalid arg",
			args{
				args: []string{"keyID", "value"},
			},
			res{
				err: helper.NewUserError("argument is not in the valid format [keyID=key]"),
			},
		},
		{
			"single arg",
			args{
				args: []string{"keyID=value"},
			},
			res{
				keys: []crypto.Key{
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
				keys: []crypto.Key{
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
			if !errors.Is(err, tt.res.err) {
				t.Errorf("keysFromArgs() error = %v, err %v", err, tt.res.err)
				return
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
		keys []crypto.Key
		err  error
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
				err: helper.NewUserError("unable to extract keys from file"),
			},
		},
		{
			"single key",
			args{
				file: bytes.NewReader([]byte("keyID: value")),
			},
			res{
				keys: []crypto.Key{
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
				keys: []crypto.Key{
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
			if !errors.Is(err, tt.res.err) {
				t.Errorf("keysFromArgs() error = %v, err %v", err, tt.res.err)
				return
			}
			assert.EqualValues(t, got, tt.res.keys)
		})
	}
}
