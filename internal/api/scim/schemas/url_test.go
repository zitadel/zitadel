package schemas

import (
	"reflect"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/zitadel/logging"
)

func TestHttpURL_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		u       *HttpURL
		want    []byte
		wantErr bool
	}{
		{
			name:    "http url",
			u:       mustParseURL("http://example.com"),
			want:    []byte(`"http://example.com"`),
			wantErr: false,
		},
		{
			name:    "https url",
			u:       mustParseURL("https://example.com"),
			want:    []byte(`"https://example.com"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, string(got), string(tt.want))
		})
	}
}

func TestHttpURL_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    *HttpURL
		wantErr bool
	}{
		{
			name:    "http url",
			data:    []byte(`"http://example.com"`),
			want:    mustParseURL("http://example.com"),
			wantErr: false,
		},
		{
			name:    "https url",
			data:    []byte(`"https://example.com"`),
			want:    mustParseURL("https://example.com"),
			wantErr: false,
		},
		{
			name:    "ftp url should fail",
			data:    []byte(`"ftp://example.com"`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no url should fail",
			data:    []byte(`"test"`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "number should fail",
			data:    []byte(`120`),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := new(HttpURL)
			err := json.Unmarshal(tt.data, url)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			assert.Equal(t, tt.want.String(), url.String())
		})
	}
}

func TestHttpURL_String(t *testing.T) {
	tests := []struct {
		name string
		u    *HttpURL
		want string
	}{
		{
			name: "http url",
			u:    mustParseURL("http://example.com"),
			want: "http://example.com",
		},
		{
			name: "https url",
			u:    mustParseURL("https://example.com"),
			want: "https://example.com",
		},
		{
			name: "nil",
			u:    nil,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseHTTPURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		want    *HttpURL
		wantErr bool
	}{
		{
			name:    "http url",
			rawURL:  "http://example.com",
			want:    mustParseURL("http://example.com"),
			wantErr: false,
		},
		{
			name:    "https url",
			rawURL:  "https://example.com",
			want:    mustParseURL("https://example.com"),
			wantErr: false,
		},
		{
			name:    "ftp url should fail",
			rawURL:  "ftp://example.com",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no url should fail",
			rawURL:  "test",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHTTPURL(tt.rawURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHTTPURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseHTTPURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func mustParseURL(rawURL string) *HttpURL {
	url, err := ParseHTTPURL(rawURL)
	logging.OnError(err).Fatal("failed to parse URL")
	return url
}
