package actions

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/dop251/goja"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/logstore"
)

func Test_isHostBlocked(t *testing.T) {
	SetLogstoreService(logstore.New(nil, nil, nil))
	var denyList = []AddressChecker{
		mustNewIPChecker(t, "192.168.5.0/24"),
		mustNewIPChecker(t, "127.0.0.1"),
		&DomainChecker{Domain: "test.com"},
	}
	type args struct {
		address *url.URL
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "in range",
			args: args{
				address: mustNewURL(t, "https://192.168.5.4/hodor"),
			},
			want: true,
		},
		{
			name: "exact ip",
			args: args{
				address: mustNewURL(t, "http://127.0.0.1:8080/hodor"),
			},
			want: true,
		},
		{
			name: "address match",
			args: args{
				address: mustNewURL(t, "https://test.com:42/hodor"),
			},
			want: true,
		},
		{
			name: "address not match",
			args: args{
				address: mustNewURL(t, "https://test2.com/hodor"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHostBlocked(denyList, tt.args.address); got != tt.want {
				t.Errorf("isHostBlocked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mustNewIPChecker(t *testing.T, ip string) AddressChecker {
	t.Helper()
	checker, err := NewIPChecker(ip)
	if err != nil {
		t.Errorf("unable to parse cidr of %q because: %v", ip, err)
		t.FailNow()
	}
	return checker
}

func mustNewURL(t *testing.T, raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		t.Errorf("unable to parse address of %q because: %v", raw, err)
		t.FailNow()
	}
	return u
}

func TestHTTP_fetchConfigFromArg(t *testing.T) {
	runtime := goja.New()
	runtime.SetFieldNameMapper(goja.UncapFieldNameMapper())

	type args struct {
		arg *goja.Object
	}
	tests := []struct {
		name       string
		args       args
		wantConfig fetchConfig
		wantErr    func(error) bool
	}{
		{
			name: "no fetch option provided",
			args: args{
				arg: runtime.ToValue(
					struct{}{},
				).ToObject(runtime),
			},
			wantConfig: fetchConfig{},
			wantErr: func(err error) bool {
				return err == nil
			},
		},
		{
			name: "header set as string",
			args: args{
				arg: runtime.ToValue(
					&struct {
						Headers map[string]string
					}{
						Headers: map[string]string{
							"Authorization": "Bearer token",
						},
					},
				).ToObject(runtime),
			},
			wantConfig: fetchConfig{
				Headers: http.Header{
					"Authorization": {"Bearer token"},
				},
			},
			wantErr: func(err error) bool {
				return err == nil
			},
		},
		{
			name: "header set as list",
			args: args{
				arg: runtime.ToValue(
					&struct {
						Headers map[string][]any
					}{
						Headers: map[string][]any{
							"Authorization": {"Bearer token"},
						},
					},
				).ToObject(runtime),
			},
			wantConfig: fetchConfig{
				Headers: http.Header{
					"Authorization": {"Bearer token"},
				},
			},
			wantErr: func(err error) bool {
				return err == nil
			},
		},
		{
			name: "method set",
			args: args{
				arg: runtime.ToValue(
					&struct {
						Method string
					}{
						Method: http.MethodPost,
					},
				).ToObject(runtime),
			},
			wantConfig: fetchConfig{
				Method: http.MethodPost,
			},
			wantErr: func(err error) bool {
				return err == nil
			},
		},
		{
			name: "body set",
			args: args{
				arg: runtime.ToValue(
					&struct {
						Body struct{ Id string }
					}{
						Body: struct{ Id string }{
							Id: "asdf123",
						},
					},
				).ToObject(runtime),
			},
			wantConfig: fetchConfig{
				Body: bytes.NewReader([]byte(`{"id":"asdf123"}`)),
			},
			wantErr: func(err error) bool {
				return err == nil
			},
		},
		{
			name: "invalid header",
			args: args{
				arg: runtime.ToValue(
					&struct {
						NotExists struct{}
					}{
						NotExists: struct{}{},
					},
				).ToObject(runtime),
			},
			wantConfig: fetchConfig{},
			wantErr: func(err error) bool {
				return errors.IsErrorInvalidArgument(err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &HTTP{
				runtime: runtime,
				client:  http.DefaultClient,
			}
			gotConfig := new(fetchConfig)

			err := c.fetchConfigFromArg(tt.args.arg, gotConfig)
			if !tt.wantErr(err) {
				t.Errorf("HTTP.fetchConfigFromArg() error = %v", err)
				return
			}
			if !reflect.DeepEqual(gotConfig.Headers, tt.wantConfig.Headers) {
				t.Errorf("config.Headers got = %#v, want %#v", gotConfig.Headers, tt.wantConfig.Headers)
			}
			if gotConfig.Method != tt.wantConfig.Method {
				t.Errorf("config.Method got = %#v, want %#v", gotConfig.Method, tt.wantConfig.Method)
			}

			if tt.wantConfig.Body == nil {
				if gotConfig.Body != nil {
					t.Errorf("didn't expect a body")
				}
				return
			}

			gotBody, _ := io.ReadAll(gotConfig.Body)
			wantBody, _ := io.ReadAll(tt.wantConfig.Body)

			if !reflect.DeepEqual(gotBody, wantBody) {
				t.Errorf("config.Body got = %s, want %s", gotBody, wantBody)
			}
		})
	}
}

func TestHTTP_buildHTTPRequest(t *testing.T) {
	runtime := goja.New()
	runtime.SetFieldNameMapper(goja.UncapFieldNameMapper())

	type args struct {
		args []goja.Value
	}
	tests := []struct {
		name        string
		args        args
		wantReq     *http.Request
		shouldPanic bool
	}{
		{
			name: "only url",
			args: args{
				args: []goja.Value{
					runtime.ToValue("http://my-url.ch"),
				},
			},
			wantReq: &http.Request{
				Method: http.MethodGet,
				URL:    mustNewURL(t, "http://my-url.ch"),
				Header: defaultFetchConfig.Headers,
				Body:   nil,
			},
		},
		{
			name: "no params",
			args: args{
				args: []goja.Value{
					runtime.ToValue("http://my-url.ch"),
					runtime.ToValue(&struct{}{}),
				},
			},
			wantReq: &http.Request{
				Method: http.MethodGet,
				URL:    mustNewURL(t, "http://my-url.ch"),
				Header: defaultFetchConfig.Headers,
				Body:   nil,
			},
		},
		{
			name: "overwrite headers",
			args: args{
				args: []goja.Value{
					runtime.ToValue("http://my-url.ch"),
					runtime.ToValue(struct {
						Headers map[string][]interface{}
					}{
						Headers: map[string][]interface{}{"Authorization": {"some token"}},
					}),
				},
			},
			wantReq: &http.Request{
				Method: http.MethodGet,
				URL:    mustNewURL(t, "http://my-url.ch"),
				Header: http.Header{
					"Authorization": []string{"some token"},
				},
				Body: nil,
			},
		},
		{
			name: "post with body",
			args: args{
				args: []goja.Value{
					runtime.ToValue("http://my-url.ch"),
					runtime.ToValue(struct {
						Body struct{ MyData string }
					}{
						Body: struct{ MyData string }{MyData: "hello world"},
					}),
				},
			},
			wantReq: &http.Request{
				Method: http.MethodGet,
				URL:    mustNewURL(t, "http://my-url.ch"),
				Header: defaultFetchConfig.Headers,
				Body:   io.NopCloser(bytes.NewReader([]byte(`{"myData":"hello world"}`))),
			},
		},
		{
			name: "too many args",
			args: args{
				args: []goja.Value{
					runtime.ToValue("http://my-url.ch"),
					runtime.ToValue("http://my-url.ch"),
					runtime.ToValue("http://my-url.ch"),
				},
			},
			wantReq:     nil,
			shouldPanic: true,
		},
		{
			name: "no args",
			args: args{
				args: []goja.Value{},
			},
			wantReq:     nil,
			shouldPanic: true,
		},
		{
			name: "invalid config",
			args: args{
				args: []goja.Value{
					runtime.ToValue("http://my-url.ch"),
					runtime.ToValue(struct {
						Invalid bool
					}{
						Invalid: true,
					}),
				},
			},
			wantReq:     nil,
			shouldPanic: true,
		},
		{
			name: "invalid method",
			args: args{
				args: []goja.Value{
					runtime.ToValue("http://my-url.ch"),
					runtime.ToValue(struct {
						Method string
					}{
						Method: " asdf asdf",
					}),
				},
			},
			wantReq:     nil,
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			panicked := false
			if tt.shouldPanic {
				defer func() {
					if panicked != tt.shouldPanic {
						t.Errorf("wanted panic: %v got %v", tt.shouldPanic, panicked)
					}
				}()
				defer func() {
					recover()
					panicked = true
				}()
			}

			c := &HTTP{
				runtime: runtime,
			}

			gotReq := c.buildHTTPRequest(context.Background(), tt.args.args)

			if tt.shouldPanic {
				return
			}

			if gotReq.URL.String() != tt.wantReq.URL.String() {
				t.Errorf("url = %s, want %s", gotReq.URL, tt.wantReq.URL)
			}

			if !reflect.DeepEqual(gotReq.Header, tt.wantReq.Header) {
				t.Errorf("headers = %v, want %v", gotReq.Header, tt.wantReq.Header)
			}

			if gotReq.Method != tt.wantReq.Method {
				t.Errorf("method = %s, want %s", gotReq.Method, tt.wantReq.Method)
			}

			if tt.wantReq.Body == nil {
				if gotReq.Body != nil {
					t.Errorf("didn't expect a body")
				}
				return
			}

			gotBody, _ := io.ReadAll(gotReq.Body)
			wantBody, _ := io.ReadAll(tt.wantReq.Body)

			if !reflect.DeepEqual(gotBody, wantBody) {
				t.Errorf("config.Body got = %s, want %s", gotBody, wantBody)
			}
		})
	}
}
