package saml

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"github.com/caos/oidc/pkg/op"
	"github.com/caos/zitadel/internal/api/saml/key"
	"github.com/caos/zitadel/internal/api/saml/mock"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/golang/mock/gomock"
	dsig "github.com/russellhaering/goxmldsig"
	"gopkg.in/square/go-jose.v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIDP_certificateHandleFunc(t *testing.T) {
	type args struct {
		metadataEndpoint string
		config           *IdentityProviderConfig
		certificate      string
		key              string
	}
	type res struct {
		code int
	}

	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"certificate 1",
			args{
				metadataEndpoint: "/saml/metadata",
				config: &IdentityProviderConfig{
					SignatureAlgorithm: dsig.RSASHA256SignatureMethod,
					Metadata:           &MetadataIDP{},
					Endpoints:          &EndpointConfig{},
				},
				certificate: "-----BEGIN CERTIFICATE-----\nMIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=\n-----END CERTIFICATE-----\n",
				key:         "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC7XKdCRxUZXjdq\nVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWc\nWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2\ngIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+a\nEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7\nLfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v\n4cxTNPn/AgMBAAECggEAF+rV9yH30Ysza8GwrXCR9qDN1Dp3QmmsavnXkonEvPoq\nEr2T3o0//6mBp6CLDboMQGQBjblJwl+3Y6PgZolvHAMOsMdHfYNPEo7FSzUBzEw+\nqRrs5HkMyvoPgfV6X8F97W3tiD4Q/AmHkMILl+MxbnfPXM54gWqPuwIqxY1uaCk5\nREwyb7WBon3rd58ceOI1SLRjod6SbqWBMMSN3cJ+5VEPObFjw/RlhNQ5rBI8G5Kt\nso2zBU5C4BB2CvqlWy98WDKJkTvWHbiTjZCy8BQ+gQ6UJM2vaNELFOVpuMGQnMIi\noWiX10Jg2e1gP9j3TdrohlGF8M3+TXjSFKNmeX0DUQKBgQDx7UazUWS5RtkgnjH9\nw2xH2xkstJVD7nAS8VTxNwcrgjVXPvTJha9El904obUjyRX7ppb02tuH5ML/bZh6\n9lL4bP5+SHcJ10e4q8CK/KAGHD6BYAbaGXRq0CoSk5a3vv5XPdob4T5qKCIHFpnu\nMfbvdbEoameLOyRYOGu/yVZIiwKBgQDGQs7FRTisHV0xooiRmlvYF0dcd19qpLed\nqhgJNqBPOTEvvGvJNRoi39haEY3cuTqsxZ5FAlFlVFMUUozz+d0xBLLInoVY/Y4h\nhSdGmdw/A6oHodLqyEp3N5RZNdLlh8/nDS3xXzMotAl75bW5kc2ttcRhRdtyNJ9Z\nup0PgppO3QKBgEC45upAQz8iCiKkz+EA8C4FGqYQJcLHvmoC8GOcAioMqrKNoDVt\ns2cZbdChynEpcd0iQ058YrDnbZeiPWHgFnBp0Gf+gQI7+u8X2+oTDci0s7Au/YZJ\nuxB8YlUX8QF1clvqqzg8OVNzKy9UR5gm+9YyWVPjq5HfH6kOZx0nAxNjAoGAERt8\nqgsCC9/wxbKnpCC0oh3IG5N1WUdjTKh7sHfVN2DQ/LR+fHsniTDVg1gWbKBTDsty\nj7PWgC7ZiFxjKz45NtyX7LW4/efLFttdezsVhR500nnFMFseCdFy7Iu3afThHKfH\nehdj27RFSTqWBrAtFjsj+dzERcOCqIRwvwDe/cUCgYEA5+1mzVXDVjKsWylKJPk+\nZZA4LUfvmTj3VLNDZrlSAI/xEikCFio0QWEA2TQYTAwbXTrKwQSeHQRhv7OTc1h+\nMhpAgvs189ze5J4jiNmULEkkrO+Cxxnw8tyV+UFRZtzW9gUoVBwXiZ/Wbl9sfnlO\nwLJHc0j6OltPcPJmxHP8gQI=\n-----END PRIVATE KEY-----\n",
			},
			res{code: 200},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mock.NewMockIDPStorage(gomock.NewController(t))
			cert, err := crypto.BytesToCertificate([]byte(tt.args.certificate))
			if err != nil {
				t.Errorf("failed to parse certificate")
			}

			block, _ := pem.Decode([]byte(tt.args.key))
			b := block.Bytes
			priv, err := x509.ParsePKCS8PrivateKey(b)
			if err != nil {
				t.Errorf("failed to parse private key")
			}

			mockStorage.EXPECT().GetResponseSigningKey(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, certAndKeyCh chan<- key.CertificateAndKey) {
				certAndKeyCh <- key.CertificateAndKey{
					Certificate: &jose.SigningKey{
						Key: jose.JSONWebKey{
							Key: cert,
						},
					},
					Key: &jose.SigningKey{
						Key: jose.JSONWebKey{
							Key: priv,
						},
					},
				}
			})

			endpoint := op.NewEndpoint(tt.args.metadataEndpoint)

			idp, err := NewIdentityProvider(&endpoint, tt.args.config, mockStorage)
			if err != nil {
				t.Errorf("NewIdentityProvider() error = %v", err.Error())
			}

			req := httptest.NewRequest(http.MethodGet, idp.CertificateEndpoint.Relative(), nil)
			w := httptest.NewRecorder()

			idp.certificateHandleFunc(w, req)

			res := w.Result()
			defer func() {
				_ = res.Body.Close()
			}()
			data, err := ioutil.ReadAll(res.Body)
			if res.StatusCode != tt.res.code {
				t.Errorf("ParseCertificates() code got = %v, want %v", res.StatusCode, tt.res)
			}
			if string(data) != tt.args.certificate {
				t.Errorf("ParseCertificates() count got = %v, want %v", string(data), tt.args.certificate)
			}
		})
	}
}
