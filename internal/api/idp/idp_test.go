package idp

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/command"
	z_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/form"
)

func Test_redirectToSuccessURL(t *testing.T) {
	type args struct {
		id         string
		userID     string
		token      string
		failureURL string
		successURL string
	}
	type res struct {
		want string
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"redirect",
			args{
				id:         "id",
				token:      "token",
				failureURL: "https://example.com/failure",
				successURL: "https://example.com/success",
			},
			res{
				"https://example.com/success?id=id&token=token",
			},
		},
		{
			"redirect with userID",
			args{
				id:         "id",
				userID:     "user",
				token:      "token",
				failureURL: "https://example.com/failure",
				successURL: "https://example.com/success",
			},
			res{
				"https://example.com/success?id=id&token=token&user=user",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com", nil)
			resp := httptest.NewRecorder()

			wm := command.NewIDPIntentWriteModel(tt.args.id, tt.args.id)
			wm.FailureURL, _ = url.Parse(tt.args.failureURL)
			wm.SuccessURL, _ = url.Parse(tt.args.successURL)

			redirectToSuccessURL(resp, req, wm, tt.args.token, tt.args.userID)
			assert.Equal(t, tt.res.want, resp.Header().Get("Location"))
		})
	}
}

func Test_redirectToFailureURL(t *testing.T) {
	type args struct {
		id         string
		failureURL string
		successURL string
		err        string
		desc       string
	}
	type res struct {
		want string
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"redirect",
			args{
				id:         "id",
				failureURL: "https://example.com/failure",
				successURL: "https://example.com/success",
			},
			res{
				"https://example.com/failure?error=&error_description=&id=id",
			},
		},
		{
			"redirect with error",
			args{
				id:         "id",
				failureURL: "https://example.com/failure",
				successURL: "https://example.com/success",
				err:        "test",
				desc:       "testdesc",
			},
			res{
				"https://example.com/failure?error=test&error_description=testdesc&id=id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com", nil)
			resp := httptest.NewRecorder()

			wm := command.NewIDPIntentWriteModel(tt.args.id, tt.args.id)
			wm.FailureURL, _ = url.Parse(tt.args.failureURL)
			wm.SuccessURL, _ = url.Parse(tt.args.successURL)

			redirectToFailureURL(resp, req, wm, tt.args.err, tt.args.desc)
			assert.Equal(t, tt.res.want, resp.Header().Get("Location"))
		})
	}
}

func Test_redirectToFailureURLErr(t *testing.T) {
	type args struct {
		id         string
		failureURL string
		successURL string
		err        error
	}
	type res struct {
		want string
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"redirect with error",
			args{
				id:         "id",
				failureURL: "https://example.com/failure",
				successURL: "https://example.com/success",
				err:        z_errors.ThrowError(nil, "test", "testdesc"),
			},
			res{
				"https://example.com/failure?error=test&error_description=testdesc&id=id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com", nil)
			resp := httptest.NewRecorder()

			wm := command.NewIDPIntentWriteModel(tt.args.id, tt.args.id)
			wm.FailureURL, _ = url.Parse(tt.args.failureURL)
			wm.SuccessURL, _ = url.Parse(tt.args.successURL)

			redirectToFailureURLErr(resp, req, wm, tt.args.err)
			assert.Equal(t, tt.res.want, resp.Header().Get("Location"))
		})
	}
}

func Test_parseCallbackRequest(t *testing.T) {
	type args struct {
		url string
	}
	type res struct {
		want *externalIDPCallbackData
		err  bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"no state",
			args{
				url: "https://example.com?state=&code=code&error=error&error_description=desc",
			},
			res{
				err: true,
			},
		},
		{
			"parse",
			args{
				url: "https://example.com?state=state&code=code&error=error&error_description=desc",
			},
			res{
				want: &externalIDPCallbackData{
					State:            "state",
					Code:             "code",
					Error:            "error",
					ErrorDescription: "desc",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.args.url, nil)
			handler := Handler{parser: form.NewParser()}

			data, err := handler.parseCallbackRequest(req)
			if tt.res.err {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.res.want, data)
		})
	}
}
