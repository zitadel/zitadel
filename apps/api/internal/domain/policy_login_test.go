package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDefaultRedirectURI(t *testing.T) {
	type args struct {
		rawURL string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"invalid url, false",
			args{
				rawURL: string('\n'),
			},
			false,
		},
		{
			"empty schema, false",
			args{
				rawURL: "url",
			},
			false,
		},
		{
			"empty http host, false",
			args{
				rawURL: "http://",
			},
			false,
		},
		{
			"empty https host, false",
			args{
				rawURL: "https://",
			},
			false,
		},
		{
			"https, ok",
			args{
				rawURL: "https://test",
			},
			true,
		},
		{
			"custom schema, ok",
			args{
				rawURL: "custom://",
			},
			true,
		},
		{
			"empty url, ok",
			args{
				rawURL: "",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ValidateDefaultRedirectURI(tt.args.rawURL), "ValidateDefaultRedirectURI(%v)", tt.args.rawURL)
		})
	}
}
