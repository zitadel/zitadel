package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"testing"

	gschema "github.com/gorilla/schema"
	"github.com/stretchr/testify/require"
)

type SampleSchema struct {
	Value    *SampleSchemaValue `schema:"value"`
	IntValue int                `schema:"intvalue"`
}

type SampleSchemaValue struct{}

func (s *SampleSchemaValue) UnmarshalText(text []byte) error {
	if string(text) == "foo" {
		return nil
	}

	return errors.New("this is a test error")
}

func TestParser_UnwrapParserError(t *testing.T) {

	tests := []struct {
		name                 string
		query                string
		wantErr              bool
		assertUnwrappedError func(err error, unwrappedErr error)
	}{
		{
			name:    "unwrap ok",
			query:   "value=test",
			wantErr: true,
			assertUnwrappedError: func(_, err error) {
				require.Equal(t, "this is a test error", err.Error())
			},
		},
		{
			name:    "multiple errors",
			query:   "value=test&intvalue=foo",
			wantErr: true,
			assertUnwrappedError: func(err error, unwrappedErr error) {
				require.Equal(t, err, unwrappedErr)
			},
		},
		{
			name:  "no error",
			query: "value=foo&intvalue=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			encodedFormData := url.Values{}.Encode()
			r, err := http.NewRequest(http.MethodPost, "http://exmaple.com?"+tt.query, bytes.NewBufferString(encodedFormData))
			require.NoError(t, err)

			data := new(SampleSchema)
			err = p.Parse(r, data)
			if !tt.wantErr {
				require.NoError(t, err)
				require.Nil(t, p.UnwrapParserError(err))
				return
			}

			require.Error(t, err)
			require.IsType(t, gschema.MultiError{}, err)

			unwrappedErr := p.UnwrapParserError(err)
			require.Error(t, unwrappedErr)
			tt.assertUnwrappedError(err, unwrappedErr)
		})
	}
}
