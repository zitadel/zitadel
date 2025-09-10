package serrors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestErrorHandler(t *testing.T) {
	i18n.MustLoadSupportedLanguagesFromDir()

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "scim error",
			err:        ThrowInvalidSyntax(zerrors.ThrowInvalidArgument(nil, "FOO", "Invalid syntax")),
			wantStatus: http.StatusBadRequest,
			wantBody: `{
				"schemas":[
					"urn:ietf:params:scim:api:messages:2.0:Error",
					"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail"
				],
				"scimType":"invalidSyntax",
				"detail":"Invalid syntax",
				"status":"400",
				"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail": {
					"id":"FOO",
					"message":"Invalid syntax"
				}
			}`,
		},
		{
			name:       "zitadel error",
			err:        zerrors.ThrowInvalidArgument(nil, "FOO", "Invalid syntax"),
			wantStatus: http.StatusBadRequest,
			wantBody: `{
				"schemas":[
					"urn:ietf:params:scim:api:messages:2.0:Error",
					"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail"
				],
				"scimType":"invalidValue",
				"detail":"Invalid syntax",
				"status":"400",
				"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail": {
					"id":"FOO",
					"message":"Invalid syntax"
				}
			}`,
		},
		{
			name:       "zitadel internal error",
			err:        zerrors.ThrowInternal(nil, "FOO", "Internal error"),
			wantStatus: http.StatusInternalServerError,
			wantBody: `{
				"schemas":[
					"urn:ietf:params:scim:api:messages:2.0:Error",
					"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail"
				],
				"detail":"Internal error",
				"status":"500",
				"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail": {
					"id":"FOO",
					"message":"Internal error"
				}
			}`,
		},
		{
			name:       "unknown error",
			err:        errors.New("FOO"),
			wantStatus: http.StatusInternalServerError,
			wantBody: `{
				"schemas":[
					"urn:ietf:params:scim:api:messages:2.0:Error"
				],
				"detail":"Unknown internal server error",
				"status":"500"
			}`,
		},
		{
			name:       "no error",
			err:        nil,
			wantStatus: http.StatusOK,
			wantBody:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			recorder := httptest.NewRecorder()
			ErrorHandler(i18n.NewZitadelTranslator(language.English))(
				func(http.ResponseWriter, *http.Request) error {
					return tt.err
				}).ServeHTTP(recorder, req)
			assert.Equal(t, tt.wantStatus, recorder.Code)

			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, recorder.Body.String())
			}
		})
	}
}
