package middleware

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/i18n"
)

var (
	SupportedLanguages = []language.Tag{language.English, language.German}
)

func TestMain(m *testing.M) {
	i18n.SupportLanguages(SupportedLanguages...)
	os.Exit(m.Run())
}

func Test_statusWriter(t *testing.T) {
	tests := []struct {
		name        string
		writeStatus int
		wantStatus  int
	}{
		{
			name:       "default status is 200",
			wantStatus: 200,
		},
		{
			name:        "writes status code and body correctly",
			writeStatus: 201,
			wantStatus:  201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const body = "Hello, World!"
			rec := httptest.NewRecorder()
			sw := newStatusWriter(rec)
			if tt.writeStatus != 0 {
				sw.WriteHeader(tt.writeStatus)
			}
			sw.Write([]byte(body))
			assert.Equal(t, tt.wantStatus, sw.Status())

			resp := rec.Result()
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
			assert.Equal(t, body, rec.Body.String())
		})
	}
}
