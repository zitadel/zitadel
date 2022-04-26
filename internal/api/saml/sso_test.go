package saml

import (
	"github.com/caos/oidc/pkg/op"
	"github.com/caos/zitadel/internal/api/saml/mock"
	"github.com/caos/zitadel/internal/api/saml/serviceprovider"
	"github.com/caos/zitadel/internal/api/saml/xml/md"
	"github.com/caos/zitadel/internal/api/saml/xml/xml_dsig"
	"github.com/golang/mock/gomock"
	dsig "github.com/russellhaering/goxmldsig"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestSSO_getAcsUrlAndBindingForResponse(t *testing.T) {
	type res struct {
		acs     string
		binding string
	}
	type args struct {
		sp             *serviceprovider.ServiceProvider
		requestBinding string
	}
	tests := []struct {
		name string
		args args
		res  res
	}{{
		"sp with post and redirect, default used",
		args{
			&serviceprovider.ServiceProvider{
				Metadata: &md.EntityDescriptorType{
					SPSSODescriptor: &md.SPSSODescriptorType{
						AssertionConsumerService: []md.IndexedEndpointType{
							{Index: "1", IsDefault: "true", Binding: RedirectBinding, Location: "redirect"},
							{Index: "2", Binding: PostBinding, Location: "post"},
						},
					},
				},
			},
			RedirectBinding,
		},
		res{
			acs:     "redirect",
			binding: RedirectBinding,
		},
	},
		{
			"sp with post and redirect, first index used",
			args{
				&serviceprovider.ServiceProvider{
					Metadata: &md.EntityDescriptorType{
						SPSSODescriptor: &md.SPSSODescriptorType{
							AssertionConsumerService: []md.IndexedEndpointType{
								{Index: "1", Binding: RedirectBinding, Location: "redirect"},
								{Index: "2", Binding: PostBinding, Location: "post"},
							},
						},
					},
				},
				RedirectBinding,
			},
			res{
				acs:     "redirect",
				binding: RedirectBinding,
			},
		},
		{
			"sp with post and redirect, redirect used",
			args{
				&serviceprovider.ServiceProvider{
					Metadata: &md.EntityDescriptorType{
						SPSSODescriptor: &md.SPSSODescriptorType{
							AssertionConsumerService: []md.IndexedEndpointType{
								{Binding: RedirectBinding, Location: "redirect"},
								{Binding: PostBinding, Location: "post"},
							},
						},
					},
				},
				RedirectBinding,
			},
			res{
				acs:     "redirect",
				binding: RedirectBinding,
			},
		},
		{
			"sp with post and redirect, post used",
			args{
				&serviceprovider.ServiceProvider{
					Metadata: &md.EntityDescriptorType{
						SPSSODescriptor: &md.SPSSODescriptorType{
							AssertionConsumerService: []md.IndexedEndpointType{
								{Binding: RedirectBinding, Location: "redirect"},
								{Binding: PostBinding, Location: "post"},
							},
						},
					},
				},
				PostBinding,
			},
			res{
				acs:     "post",
				binding: PostBinding,
			},
		},
		{
			"sp with redirect, post used",
			args{
				&serviceprovider.ServiceProvider{
					Metadata: &md.EntityDescriptorType{
						SPSSODescriptor: &md.SPSSODescriptorType{
							AssertionConsumerService: []md.IndexedEndpointType{
								{Binding: RedirectBinding, Location: "redirect"},
							},
						},
					},
				},
				PostBinding,
			},
			res{
				acs:     "redirect",
				binding: RedirectBinding,
			},
		},
		{
			"sp with post, redirect used",
			args{
				&serviceprovider.ServiceProvider{
					Metadata: &md.EntityDescriptorType{
						SPSSODescriptor: &md.SPSSODescriptorType{
							AssertionConsumerService: []md.IndexedEndpointType{
								{Binding: PostBinding, Location: "post"},
							},
						},
					},
				},
				RedirectBinding,
			},
			res{
				acs:     "post",
				binding: PostBinding,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acs, binding := getAcsUrlAndBindingForResponse(tt.args.sp, tt.args.requestBinding)
			if acs != tt.res.acs && binding != tt.res.binding {
				t.Errorf("getAcsUrlAndBindingForResponse() got = %v/%v, want %v/%v", acs, binding, tt.res.acs, tt.res.binding)
				return
			}
		})
	}
}

func TestSSO_getAuthRequestFromRequest(t *testing.T) {
	type res struct {
		want *AuthRequestForm
		err  bool
	}
	tests := []struct {
		name string
		arg  *http.Request
		res  res
	}{
		{
			"parsing form error",
			&http.Request{URL: &url.URL{RawQuery: "invalid=%%param"}},
			res{
				nil,
				true,
			},
		},
		{
			"signed redirect binding",
			&http.Request{URL: &url.URL{RawQuery: "SAMLRequest=request&SAMLEncoding=encoding&RelayState=state&SigAlg=alg&Signature=sig"}},
			res{
				&AuthRequestForm{
					AuthRequest: "request",
					Encoding:    "encoding",
					RelayState:  "state",
					SigAlg:      "alg",
					Sig:         "sig",
					Binding:     RedirectBinding,
				},
				false,
			},
		},
		{
			"unsigned redirect binding",
			&http.Request{URL: &url.URL{RawQuery: "SAMLRequest=request&SAMLEncoding=encoding&RelayState=state"}},
			res{
				&AuthRequestForm{
					AuthRequest: "request",
					Encoding:    "encoding",
					RelayState:  "state",
					SigAlg:      "",
					Sig:         "",
					Binding:     RedirectBinding,
				},
				false,
			},
		},
		{
			"post binding",
			&http.Request{
				Form: map[string][]string{
					"SAMLRequest": {"request"},
					"RelayState":  {"state"},
				},
				URL: &url.URL{RawQuery: ""}},
			res{
				&AuthRequestForm{
					AuthRequest: "request",
					Encoding:    "",
					RelayState:  "state",
					SigAlg:      "",
					Sig:         "",
					Binding:     PostBinding,
				},
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAuthRequestFromRequest(tt.arg)
			if (err != nil) != tt.res.err {
				t.Errorf("getAuthRequestFromRequest() error = %v, wantErr %v", err, tt.res.err)
			}
			if !reflect.DeepEqual(got, tt.res.want) {
				t.Errorf("getAuthRequestFromRequest() got = %v, want %v", got, tt.res.want)
			}
		})
	}
}

func TestSSO_certificateCheckNecessary(t *testing.T) {
	type args struct {
		sig      *xml_dsig.SignatureType
		metadata *md.EntityDescriptorType
	}
	tests := []struct {
		name string
		args args
		res  bool
	}{
		{
			"sig nil",
			args{
				sig:      nil,
				metadata: &md.EntityDescriptorType{},
			},
			false,
		},
		{
			"keyinfo nil",
			args{
				sig:      &xml_dsig.SignatureType{KeyInfo: nil},
				metadata: &md.EntityDescriptorType{},
			},
			false,
		},
		{
			"keydescriptor nil",
			args{
				sig: &xml_dsig.SignatureType{KeyInfo: &xml_dsig.KeyInfoType{}},
				metadata: &md.EntityDescriptorType{
					SPSSODescriptor: &md.SPSSODescriptorType{
						KeyDescriptor: nil,
					},
				},
			},
			false,
		},
		{
			"keydescriptor length == 0",
			args{
				sig: &xml_dsig.SignatureType{KeyInfo: &xml_dsig.KeyInfoType{}},
				metadata: &md.EntityDescriptorType{
					SPSSODescriptor: &md.SPSSODescriptorType{
						KeyDescriptor: []md.KeyDescriptorType{},
					},
				},
			},
			false,
		},
		{
			"check necessary",
			args{
				sig: &xml_dsig.SignatureType{KeyInfo: &xml_dsig.KeyInfoType{}},
				metadata: &md.EntityDescriptorType{
					SPSSODescriptor: &md.SPSSODescriptorType{
						KeyDescriptor: []md.KeyDescriptorType{{Use: "test"}},
					},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authRequestF := func() *xml_dsig.SignatureType {
				return tt.args.sig
			}
			metadataF := func() *md.EntityDescriptorType {
				return tt.args.metadata
			}

			gotF := certificateCheckNecessary(authRequestF, metadataF)
			got := gotF()
			if got != tt.res {
				t.Errorf("certificateCheckNecessary() got = %v, want %v", got, tt.res)
			}
		})
	}
}

func TestSSO_checkCertificate(t *testing.T) {
	type args struct {
		sig      *xml_dsig.SignatureType
		metadata *md.EntityDescriptorType
	}
	tests := []struct {
		name string
		args args
		err  bool
	}{
		{
			"keydescriptor length == 0",
			args{
				sig: &xml_dsig.SignatureType{KeyInfo: &xml_dsig.KeyInfoType{}},
				metadata: &md.EntityDescriptorType{
					SPSSODescriptor: &md.SPSSODescriptorType{
						KeyDescriptor: []md.KeyDescriptorType{},
					},
				},
			},
			true,
		},
		{
			"x509data length == 0",
			args{
				sig: &xml_dsig.SignatureType{KeyInfo: &xml_dsig.KeyInfoType{X509Data: []xml_dsig.X509DataType{}}},
				metadata: &md.EntityDescriptorType{
					SPSSODescriptor: &md.SPSSODescriptorType{
						KeyDescriptor: []md.KeyDescriptorType{{Use: "test"}},
					},
				},
			},
			true,
		},
		{
			"certificates equal",
			args{
				sig: &xml_dsig.SignatureType{KeyInfo: &xml_dsig.KeyInfoType{X509Data: []xml_dsig.X509DataType{{X509Certificate: "test"}}}},
				metadata: &md.EntityDescriptorType{
					SPSSODescriptor: &md.SPSSODescriptorType{
						KeyDescriptor: []md.KeyDescriptorType{{Use: "test", KeyInfo: xml_dsig.KeyInfoType{X509Data: []xml_dsig.X509DataType{{X509Certificate: "test"}}}}},
					},
				},
			},
			false,
		},
		{
			"certificates not equal",
			args{
				sig: &xml_dsig.SignatureType{KeyInfo: &xml_dsig.KeyInfoType{X509Data: []xml_dsig.X509DataType{{X509Certificate: "test1"}}}},
				metadata: &md.EntityDescriptorType{
					SPSSODescriptor: &md.SPSSODescriptorType{
						KeyDescriptor: []md.KeyDescriptorType{{Use: "test", KeyInfo: xml_dsig.KeyInfoType{X509Data: []xml_dsig.X509DataType{{X509Certificate: "test2"}}}}},
					},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authRequestF := func() *xml_dsig.SignatureType {
				return tt.args.sig
			}
			metadataF := func() *md.EntityDescriptorType {
				return tt.args.metadata
			}

			gotF := checkCertificate(authRequestF, metadataF)
			got := gotF()
			if (got != nil) != tt.err {
				t.Errorf("checkCertificate() got = %v, want %v", got, tt.err)
			}
		})
	}
}

func TestSSO_ssoHandleFunc(t *testing.T) {
	type request struct {
		ID           string
		SAMLRequest  string
		Signature    string
		SigAlg       string
		SAMLEncoding string
		RelayState   string
		Binding      string
	}
	type res struct {
		code int
		err  bool
	}
	type sp struct {
		entityID string
		metadata string
	}
	type args struct {
		metadataEndpoint string
		config           *IdentityProviderConfig
		certificate      string
		key              string
		request          request
		sp               sp
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"signed redirect request",
			args{
				metadataEndpoint: "/saml/metadata",
				config: &IdentityProviderConfig{
					SignatureAlgorithm: dsig.RSASHA256SignatureMethod,
					Metadata:           &MetadataIDP{},
					Endpoints: &EndpointConfig{
						SingleSignOn: Endpoint{URL: "http://localhost:50002/saml/SSO", Path: "/saml/SSO"},
					},
				},
				certificate: "-----BEGIN CERTIFICATE-----\nMIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=\n-----END CERTIFICATE-----\n",
				key:         "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC7XKdCRxUZXjdq\nVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWc\nWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2\ngIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+a\nEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7\nLfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v\n4cxTNPn/AgMBAAECggEAF+rV9yH30Ysza8GwrXCR9qDN1Dp3QmmsavnXkonEvPoq\nEr2T3o0//6mBp6CLDboMQGQBjblJwl+3Y6PgZolvHAMOsMdHfYNPEo7FSzUBzEw+\nqRrs5HkMyvoPgfV6X8F97W3tiD4Q/AmHkMILl+MxbnfPXM54gWqPuwIqxY1uaCk5\nREwyb7WBon3rd58ceOI1SLRjod6SbqWBMMSN3cJ+5VEPObFjw/RlhNQ5rBI8G5Kt\nso2zBU5C4BB2CvqlWy98WDKJkTvWHbiTjZCy8BQ+gQ6UJM2vaNELFOVpuMGQnMIi\noWiX10Jg2e1gP9j3TdrohlGF8M3+TXjSFKNmeX0DUQKBgQDx7UazUWS5RtkgnjH9\nw2xH2xkstJVD7nAS8VTxNwcrgjVXPvTJha9El904obUjyRX7ppb02tuH5ML/bZh6\n9lL4bP5+SHcJ10e4q8CK/KAGHD6BYAbaGXRq0CoSk5a3vv5XPdob4T5qKCIHFpnu\nMfbvdbEoameLOyRYOGu/yVZIiwKBgQDGQs7FRTisHV0xooiRmlvYF0dcd19qpLed\nqhgJNqBPOTEvvGvJNRoi39haEY3cuTqsxZ5FAlFlVFMUUozz+d0xBLLInoVY/Y4h\nhSdGmdw/A6oHodLqyEp3N5RZNdLlh8/nDS3xXzMotAl75bW5kc2ttcRhRdtyNJ9Z\nup0PgppO3QKBgEC45upAQz8iCiKkz+EA8C4FGqYQJcLHvmoC8GOcAioMqrKNoDVt\ns2cZbdChynEpcd0iQ058YrDnbZeiPWHgFnBp0Gf+gQI7+u8X2+oTDci0s7Au/YZJ\nuxB8YlUX8QF1clvqqzg8OVNzKy9UR5gm+9YyWVPjq5HfH6kOZx0nAxNjAoGAERt8\nqgsCC9/wxbKnpCC0oh3IG5N1WUdjTKh7sHfVN2DQ/LR+fHsniTDVg1gWbKBTDsty\nj7PWgC7ZiFxjKz45NtyX7LW4/efLFttdezsVhR500nnFMFseCdFy7Iu3afThHKfH\nehdj27RFSTqWBrAtFjsj+dzERcOCqIRwvwDe/cUCgYEA5+1mzVXDVjKsWylKJPk+\nZZA4LUfvmTj3VLNDZrlSAI/xEikCFio0QWEA2TQYTAwbXTrKwQSeHQRhv7OTc1h+\nMhpAgvs189ze5J4jiNmULEkkrO+Cxxnw8tyV+UFRZtzW9gUoVBwXiZ/Wbl9sfnlO\nwLJHc0j6OltPcPJmxHP8gQI=\n-----END PRIVATE KEY-----\n",
				request: request{
					ID:          "test",
					Binding:     RedirectBinding,
					SAMLRequest: "nJJBj9MwEIX/ijX3NG6a7DbWJlLZClFpYatN4cBt6k6oJccungmw/x61XaQioRy42vP5ved5D4yDP5nVKMfwQt9HYlG/Bh/YnC8aGFMwEdmxCTgQG7GmW318MsVMG2SmJC4GuEFO08wpRYk2elCbdQPukFlNd/c9LQpczPve6r3taVHWdbWoal3bfr7c03JJc1BfKLGLoYFipkFtmEfaBBYM0kChiyLTZVbc7XRtyntTVrOyrr6CWhOLCygX8ihyMnnuo0V/jCym0loX+dl33nXPoFZ/Ij3GwONAqaP0w1n6/PL0D3qptb7CaBnU9i3bOxcOLnyb/oj9dYjNh91um22fux20l2WYS7Kk3sc0oEw/cj5xh6y/jBoK4uQV2gmfAwkeUPAhv5Fq30rwCQfarLfRO/v6H/KSMLCjIKBW3sefj4lQqAFJI0HeXiX/rlr7OwAA//8=",
					RelayState:  "K6LS7mdqUO4SGedbfa8nBIyX-7K8gGbrHMqIMwVn6zCKLLoADHjEHUAm",
					Signature:   "PWZ6JPNpAGE7mYLKD3dCUG9AZcThrMRQGtvdv31ewx3hms5Oglc677iAUEcbIBrvKtMrCPVwXPNxT6wQ0rg4qIgyKgoyS53ZTaxaFHPrB7wkkzqtK7GvWgdEqceT8iooK5SCLHFMJ3m30LqEbX7zFw62yE34+e7ypfZSM5Lrf0QFwPzX+LNCuYA+Ob9D5SKc132tn21J2vBRmNJ1zCY0ksRzQfyfErjAzcGVx8qK9jpaeyvsVBZSkH/I6+1hb8lQWE48xala9NbqfbMATGBCQj1UvpVMMfp6PE7KPk5Y1YDeSqPeRIEKH+Gnip6Hve5Ji1aiRp5bytVf1VHwTHSq8w==",
					SigAlg:      "http://www.w3.org/2000/09/xmldsig#rsa-sha1",
				},
				sp: sp{
					entityID: "http://localhost:8000/saml/metadata",
					metadata: "PEVudGl0eURlc2NyaXB0b3IgeG1sbnM9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDptZXRhZGF0YSIgdmFsaWRVbnRpbD0iMjAyMi0wNC0yOFQxMTozMjowNC43OTdaIiBlbnRpdHlJRD0iaHR0cDovL2xvY2FsaG9zdDo4MDAwL3NhbWwvbWV0YWRhdGEiPgogIDxTUFNTT0Rlc2NyaXB0b3IgeG1sbnM9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDptZXRhZGF0YSIgdmFsaWRVbnRpbD0iMjAyMi0wNC0yOFQxMTozMjowNC43OTY5MjNaIiBwcm90b2NvbFN1cHBvcnRFbnVtZXJhdGlvbj0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOnByb3RvY29sIiBBdXRoblJlcXVlc3RzU2lnbmVkPSJ0cnVlIiBXYW50QXNzZXJ0aW9uc1NpZ25lZD0idHJ1ZSI+CiAgICA8S2V5RGVzY3JpcHRvciB1c2U9ImVuY3J5cHRpb24iPgogICAgICA8S2V5SW5mbyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC8wOS94bWxkc2lnIyI+CiAgICAgICAgPFg1MDlEYXRhIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwLzA5L3htbGRzaWcjIj4KICAgICAgICAgIDxYNTA5Q2VydGlmaWNhdGUgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyMiPk1JSUN2RENDQWFRQ0NRRDZFOFpHc1EydXNqQU5CZ2txaGtpRzl3MEJBUXNGQURBZ01SNHdIQVlEVlFRRERCVnRlWE5sY25acFkyVXVaWGhoYlhCc1pTNWpiMjB3SGhjTk1qSXdNakUzTVRRd05qTTVXaGNOTWpNd01qRTNNVFF3TmpNNVdqQWdNUjR3SEFZRFZRUUREQlZ0ZVhObGNuWnBZMlV1WlhoaGJYQnNaUzVqYjIwd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFS0FvSUJBUUM3WEtkQ1J4VVpYamRxVnF3d3dPSnFjMUNoMG5PU21rK1VlcmtVcWx2aVdIZGVMUitGb2xIS2pxTHpDQmxvQXo0eFZjMERGZlI3NmdXY1dBSEpsb3FaN0dCUzdOcERoelY4RytjWFErYlRVMEx1MmU3M3pDUWIzMFhVZEtoV2lHZkRLYVUrMXhnOUNELzJnSWZzWVBzM1RUcTFzcTdvQ3M1cUxkVUhhVkw1a2NSYUhLZG5UaTdjczVpOXh6czNUc1VuWGNySlB3eWRqcCthRWt5UmgwN29NcFhCRW9iR2lzZkYycDFNQTZwVlcyZ2pteXdmN0Q1aVlFRkVMUWhNN3BvcVBOMy9rZkJ2VTFuN0xmZ3E3b3htdi84TEZpNFpvcHI1bnlxc3oyNlhQdFV5MVdxVHpnem5BbVArbk4wb0JURVJGVmJYWGRSYTNrMnY0Y3hUTlBuL0FnTUJBQUV3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUpZeFJPV1NPWmJPelh6YWZkR2pRS3NNZ045NDhHL2hId1Z1Wm5leUFjVm9MTUZUczFXZXlhOVorc25NcDF1MEFkREdtUVRTOXpHbkQ3c3lEWUdPbWdpZ09MY012TE1vV2Y1dENRQmJFdWtXOE83RFBqUlIwWHlwQ2hHU3NIc3FMR08wQjBIYVRlbDBIZFA5U2k4MjdPQ2tjOVErV2JzRkcvOC80VG9HV0wrdWxhMVd1TGF3b3pvajh1bVBpOUQ4aVhDb1czNXkyU1RVK1dGUUc3VytLZmR1KzJDWXovMHRHZHdWcU5HNFdzZmF3V2NoclMwMHZHRktqbS9mSmM4NzZnQWZ4aU1IMUk5Zlp2WVNBeEFaM3NWSS8vTWwyc1VkZ2YwNjd5d1E3NW9hTFNTMk5JbW16NWFvczN2dVdtT1hoSUxkN2lUVStCRDhVdjZ2V2JJN0kxTT08L1g1MDlDZXJ0aWZpY2F0ZT4KICAgICAgICA8L1g1MDlEYXRhPgogICAgICA8L0tleUluZm8+CiAgICAgIDxFbmNyeXB0aW9uTWV0aG9kIEFsZ29yaXRobT0iaHR0cDovL3d3dy53My5vcmcvMjAwMS8wNC94bWxlbmMjYWVzMTI4LWNiYyI+PC9FbmNyeXB0aW9uTWV0aG9kPgogICAgICA8RW5jcnlwdGlvbk1ldGhvZCBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDEvMDQveG1sZW5jI2FlczE5Mi1jYmMiPjwvRW5jcnlwdGlvbk1ldGhvZD4KICAgICAgPEVuY3J5cHRpb25NZXRob2QgQWxnb3JpdGhtPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxLzA0L3htbGVuYyNhZXMyNTYtY2JjIj48L0VuY3J5cHRpb25NZXRob2Q+CiAgICAgIDxFbmNyeXB0aW9uTWV0aG9kIEFsZ29yaXRobT0iaHR0cDovL3d3dy53My5vcmcvMjAwMS8wNC94bWxlbmMjcnNhLW9hZXAtbWdmMXAiPjwvRW5jcnlwdGlvbk1ldGhvZD4KICAgIDwvS2V5RGVzY3JpcHRvcj4KICAgIDxLZXlEZXNjcmlwdG9yIHVzZT0ic2lnbmluZyI+CiAgICAgIDxLZXlJbmZvIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwLzA5L3htbGRzaWcjIj4KICAgICAgICA8WDUwOURhdGEgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyMiPgogICAgICAgICAgPFg1MDlDZXJ0aWZpY2F0ZSB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC8wOS94bWxkc2lnIyI+TUlJQ3ZEQ0NBYVFDQ1FENkU4WkdzUTJ1c2pBTkJna3Foa2lHOXcwQkFRc0ZBREFnTVI0d0hBWURWUVFEREJWdGVYTmxjblpwWTJVdVpYaGhiWEJzWlM1amIyMHdIaGNOTWpJd01qRTNNVFF3TmpNNVdoY05Nak13TWpFM01UUXdOak01V2pBZ01SNHdIQVlEVlFRRERCVnRlWE5sY25acFkyVXVaWGhoYlhCc1pTNWpiMjB3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRQzdYS2RDUnhVWlhqZHFWcXd3d09KcWMxQ2gwbk9TbWsrVWVya1VxbHZpV0hkZUxSK0ZvbEhLanFMekNCbG9BejR4VmMwREZmUjc2Z1djV0FISmxvcVo3R0JTN05wRGh6VjhHK2NYUStiVFUwTHUyZTczekNRYjMwWFVkS2hXaUdmREthVSsxeGc5Q0QvMmdJZnNZUHMzVFRxMXNxN29DczVxTGRVSGFWTDVrY1JhSEtkblRpN2NzNWk5eHpzM1RzVW5YY3JKUHd5ZGpwK2FFa3lSaDA3b01wWEJFb2JHaXNmRjJwMU1BNnBWVzJnam15d2Y3RDVpWUVGRUxRaE03cG9xUE4zL2tmQnZVMW43TGZncTdveG12LzhMRmk0Wm9wcjVueXFzejI2WFB0VXkxV3FUemd6bkFtUCtuTjBvQlRFUkZWYlhYZFJhM2sydjRjeFROUG4vQWdNQkFBRXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSll4Uk9XU09aYk96WHphZmRHalFLc01nTjk0OEcvaEh3VnVabmV5QWNWb0xNRlRzMVdleWE5Witzbk1wMXUwQWRER21RVFM5ekduRDdzeURZR09tZ2lnT0xjTXZMTW9XZjV0Q1FCYkV1a1c4TzdEUGpSUjBYeXBDaEdTc0hzcUxHTzBCMEhhVGVsMEhkUDlTaTgyN09Da2M5UStXYnNGRy84LzRUb0dXTCt1bGExV3VMYXdvem9qOHVtUGk5RDhpWENvVzM1eTJTVFUrV0ZRRzdXK0tmZHUrMkNZei8wdEdkd1ZxTkc0V3NmYXdXY2hyUzAwdkdGS2ptL2ZKYzg3NmdBZnhpTUgxSTlmWnZZU0F4QVozc1ZJLy9NbDJzVWRnZjA2N3l3UTc1b2FMU1MyTkltbXo1YW9zM3Z1V21PWGhJTGQ3aVRVK0JEOFV2NnZXYkk3STFNPTwvWDUwOUNlcnRpZmljYXRlPgogICAgICAgIDwvWDUwOURhdGE+CiAgICAgIDwvS2V5SW5mbz4KICAgIDwvS2V5RGVzY3JpcHRvcj4KICAgIDxTaW5nbGVMb2dvdXRTZXJ2aWNlIEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLVBPU1QiIExvY2F0aW9uPSJodHRwOi8vbG9jYWxob3N0OjgwMDAvc2FtbC9zbG8iIFJlc3BvbnNlTG9jYXRpb249Imh0dHA6Ly9sb2NhbGhvc3Q6ODAwMC9zYW1sL3NsbyI+PC9TaW5nbGVMb2dvdXRTZXJ2aWNlPgogICAgPEFzc2VydGlvbkNvbnN1bWVyU2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1QT1NUIiBMb2NhdGlvbj0iaHR0cDovL2xvY2FsaG9zdDo4MDAwL3NhbWwvYWNzIiBpbmRleD0iMSI+PC9Bc3NlcnRpb25Db25zdW1lclNlcnZpY2U+CiAgICA8QXNzZXJ0aW9uQ29uc3VtZXJTZXJ2aWNlIEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLUFydGlmYWN0IiBMb2NhdGlvbj0iaHR0cDovL2xvY2FsaG9zdDo4MDAwL3NhbWwvYWNzIiBpbmRleD0iMiI+PC9Bc3NlcnRpb25Db25zdW1lclNlcnZpY2U+CiAgPC9TUFNTT0Rlc2NyaXB0b3I+CjwvRW50aXR5RGVzY3JpcHRvcj4=",
				},
			},
			res{
				code: 303,
				err:  false,
			}},
		{
			"signed post request",
			args{
				metadataEndpoint: "/saml/metadata",
				config: &IdentityProviderConfig{
					SignatureAlgorithm: dsig.RSASHA256SignatureMethod,
					Metadata:           &MetadataIDP{},
					Endpoints: &EndpointConfig{
						SingleSignOn: Endpoint{URL: "http://localhost:50002/saml/SSO", Path: "/saml/SSO"},
					},
				},
				certificate: "-----BEGIN CERTIFICATE-----\nMIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=\n-----END CERTIFICATE-----\n",
				key:         "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC7XKdCRxUZXjdq\nVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWc\nWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2\ngIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+a\nEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7\nLfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v\n4cxTNPn/AgMBAAECggEAF+rV9yH30Ysza8GwrXCR9qDN1Dp3QmmsavnXkonEvPoq\nEr2T3o0//6mBp6CLDboMQGQBjblJwl+3Y6PgZolvHAMOsMdHfYNPEo7FSzUBzEw+\nqRrs5HkMyvoPgfV6X8F97W3tiD4Q/AmHkMILl+MxbnfPXM54gWqPuwIqxY1uaCk5\nREwyb7WBon3rd58ceOI1SLRjod6SbqWBMMSN3cJ+5VEPObFjw/RlhNQ5rBI8G5Kt\nso2zBU5C4BB2CvqlWy98WDKJkTvWHbiTjZCy8BQ+gQ6UJM2vaNELFOVpuMGQnMIi\noWiX10Jg2e1gP9j3TdrohlGF8M3+TXjSFKNmeX0DUQKBgQDx7UazUWS5RtkgnjH9\nw2xH2xkstJVD7nAS8VTxNwcrgjVXPvTJha9El904obUjyRX7ppb02tuH5ML/bZh6\n9lL4bP5+SHcJ10e4q8CK/KAGHD6BYAbaGXRq0CoSk5a3vv5XPdob4T5qKCIHFpnu\nMfbvdbEoameLOyRYOGu/yVZIiwKBgQDGQs7FRTisHV0xooiRmlvYF0dcd19qpLed\nqhgJNqBPOTEvvGvJNRoi39haEY3cuTqsxZ5FAlFlVFMUUozz+d0xBLLInoVY/Y4h\nhSdGmdw/A6oHodLqyEp3N5RZNdLlh8/nDS3xXzMotAl75bW5kc2ttcRhRdtyNJ9Z\nup0PgppO3QKBgEC45upAQz8iCiKkz+EA8C4FGqYQJcLHvmoC8GOcAioMqrKNoDVt\ns2cZbdChynEpcd0iQ058YrDnbZeiPWHgFnBp0Gf+gQI7+u8X2+oTDci0s7Au/YZJ\nuxB8YlUX8QF1clvqqzg8OVNzKy9UR5gm+9YyWVPjq5HfH6kOZx0nAxNjAoGAERt8\nqgsCC9/wxbKnpCC0oh3IG5N1WUdjTKh7sHfVN2DQ/LR+fHsniTDVg1gWbKBTDsty\nj7PWgC7ZiFxjKz45NtyX7LW4/efLFttdezsVhR500nnFMFseCdFy7Iu3afThHKfH\nehdj27RFSTqWBrAtFjsj+dzERcOCqIRwvwDe/cUCgYEA5+1mzVXDVjKsWylKJPk+\nZZA4LUfvmTj3VLNDZrlSAI/xEikCFio0QWEA2TQYTAwbXTrKwQSeHQRhv7OTc1h+\nMhpAgvs189ze5J4jiNmULEkkrO+Cxxnw8tyV+UFRZtzW9gUoVBwXiZ/Wbl9sfnlO\nwLJHc0j6OltPcPJmxHP8gQI=\n-----END PRIVATE KEY-----\n",
				request: request{
					ID:          "test",
					Binding:     PostBinding,
					SAMLRequest: "PHNhbWxwOkF1dGhuUmVxdWVzdCB4bWxuczpzYW1sPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YXNzZXJ0aW9uIiB4bWxuczpzYW1scD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOnByb3RvY29sIiBJRD0iaWQtOGZjNTc5MTI4YWNmYjYzNTAzODYwYjI0MTkzOGRmOTBlNWE5YzdhMiIgVmVyc2lvbj0iMi4wIiBJc3N1ZUluc3RhbnQ9IjIwMjItMDQtMjZUMTI6NTY6NDYuNDkzWiIgRGVzdGluYXRpb249Imh0dHA6Ly9sb2NhbGhvc3Q6NTAwMDIvc2FtbC9TU08iIEFzc2VydGlvbkNvbnN1bWVyU2VydmljZVVSTD0iaHR0cDovL2xvY2FsaG9zdDo4MDAwL3NhbWwvYWNzIiBQcm90b2NvbEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLVBPU1QiPjxzYW1sOklzc3VlciBGb3JtYXQ9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpuYW1laWQtZm9ybWF0OmVudGl0eSI+aHR0cDovL2xvY2FsaG9zdDo4MDAwL3NhbWwvbWV0YWRhdGE8L3NhbWw6SXNzdWVyPjxkczpTaWduYXR1cmUgeG1sbnM6ZHM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyMiPjxkczpTaWduZWRJbmZvPjxkczpDYW5vbmljYWxpemF0aW9uTWV0aG9kIEFsZ29yaXRobT0iaHR0cDovL3d3dy53My5vcmcvMjAwMS8xMC94bWwtZXhjLWMxNG4jIi8+PGRzOlNpZ25hdHVyZU1ldGhvZCBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyNyc2Etc2hhMSIvPjxkczpSZWZlcmVuY2UgVVJJPSIjaWQtOGZjNTc5MTI4YWNmYjYzNTAzODYwYjI0MTkzOGRmOTBlNWE5YzdhMiI+PGRzOlRyYW5zZm9ybXM+PGRzOlRyYW5zZm9ybSBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyNlbnZlbG9wZWQtc2lnbmF0dXJlIi8+PGRzOlRyYW5zZm9ybSBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDEvMTAveG1sLWV4Yy1jMTRuIyIvPjwvZHM6VHJhbnNmb3Jtcz48ZHM6RGlnZXN0TWV0aG9kIEFsZ29yaXRobT0iaHR0cDovL3d3dy53My5vcmcvMjAwMC8wOS94bWxkc2lnI3NoYTEiLz48ZHM6RGlnZXN0VmFsdWU+NXVJVnJiRVg2d1NGbEJpUUx5VFJGeW5naTYwPTwvZHM6RGlnZXN0VmFsdWU+PC9kczpSZWZlcmVuY2U+PC9kczpTaWduZWRJbmZvPjxkczpTaWduYXR1cmVWYWx1ZT5NMXB0S3laVTJCYzhlSUVUSUlSRHUxK256R0ZXRUI0RFlPZmpRMVFXYkZSVEUxK2RJYm5TMzlBRTgvZEM4aXhYcmdWMGVRckdCMytDWW1VeFpCa3RNbFVzd09scVIzNXJyY1NBZ0lVUDVZODQwb3FLbWJBaDJTNC8rQlVvRHgzaytLdGFNZEpWNXlpSWZZMHN2aUhwYllsMFhNMGNJTVdKTEtBc1ZReGJhQTJENTh2VXJyZDZWczZqRGsreUh3WFU2ODNhd2ZPZGpsQktwZmg4a21hb1poVk51Y1ZpT2NUeU5XWDg4Y3ltbkwyZ2UzR0dpZVU4bHBzd3lMK2UwclROL3FaUkVKL3RrNlJkZDJiRzNGSVo1TW1Fd3dsOHRwd0VqUjFzcDZhSEFVejVmT1Bmc0x1QUlZd1M4T2RBNklVcGFGWDVLeENYSkp2NzFhWFMrbVBOWFE9PTwvZHM6U2lnbmF0dXJlVmFsdWU+PGRzOktleUluZm8+PGRzOlg1MDlEYXRhPjxkczpYNTA5Q2VydGlmaWNhdGU+TUlJQ3ZEQ0NBYVFDQ1FENkU4WkdzUTJ1c2pBTkJna3Foa2lHOXcwQkFRc0ZBREFnTVI0d0hBWURWUVFEREJWdGVYTmxjblpwWTJVdVpYaGhiWEJzWlM1amIyMHdIaGNOTWpJd01qRTNNVFF3TmpNNVdoY05Nak13TWpFM01UUXdOak01V2pBZ01SNHdIQVlEVlFRRERCVnRlWE5sY25acFkyVXVaWGhoYlhCc1pTNWpiMjB3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRQzdYS2RDUnhVWlhqZHFWcXd3d09KcWMxQ2gwbk9TbWsrVWVya1VxbHZpV0hkZUxSK0ZvbEhLanFMekNCbG9BejR4VmMwREZmUjc2Z1djV0FISmxvcVo3R0JTN05wRGh6VjhHK2NYUStiVFUwTHUyZTczekNRYjMwWFVkS2hXaUdmREthVSsxeGc5Q0QvMmdJZnNZUHMzVFRxMXNxN29DczVxTGRVSGFWTDVrY1JhSEtkblRpN2NzNWk5eHpzM1RzVW5YY3JKUHd5ZGpwK2FFa3lSaDA3b01wWEJFb2JHaXNmRjJwMU1BNnBWVzJnam15d2Y3RDVpWUVGRUxRaE03cG9xUE4zL2tmQnZVMW43TGZncTdveG12LzhMRmk0Wm9wcjVueXFzejI2WFB0VXkxV3FUemd6bkFtUCtuTjBvQlRFUkZWYlhYZFJhM2sydjRjeFROUG4vQWdNQkFBRXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSll4Uk9XU09aYk96WHphZmRHalFLc01nTjk0OEcvaEh3VnVabmV5QWNWb0xNRlRzMVdleWE5Witzbk1wMXUwQWRER21RVFM5ekduRDdzeURZR09tZ2lnT0xjTXZMTW9XZjV0Q1FCYkV1a1c4TzdEUGpSUjBYeXBDaEdTc0hzcUxHTzBCMEhhVGVsMEhkUDlTaTgyN09Da2M5UStXYnNGRy84LzRUb0dXTCt1bGExV3VMYXdvem9qOHVtUGk5RDhpWENvVzM1eTJTVFUrV0ZRRzdXK0tmZHUrMkNZei8wdEdkd1ZxTkc0V3NmYXdXY2hyUzAwdkdGS2ptL2ZKYzg3NmdBZnhpTUgxSTlmWnZZU0F4QVozc1ZJLy9NbDJzVWRnZjA2N3l3UTc1b2FMU1MyTkltbXo1YW9zM3Z1V21PWGhJTGQ3aVRVK0JEOFV2NnZXYkk3STFNPTwvZHM6WDUwOUNlcnRpZmljYXRlPjwvZHM6WDUwOURhdGE+PC9kczpLZXlJbmZvPjwvZHM6U2lnbmF0dXJlPjxzYW1scDpOYW1lSURQb2xpY3kgRm9ybWF0PSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6bmFtZWlkLWZvcm1hdDp0cmFuc2llbnQiIEFsbG93Q3JlYXRlPSJ0cnVlIi8+PC9zYW1scDpBdXRoblJlcXVlc3Q+",
					RelayState:  "gcH0iLEWgKzJM0rlyJZAhtn8xcY1i85ONmb-FfHwPWO4yIMZSZmUT6aV",
					Signature:   "",
					SigAlg:      "",
				},
				sp: sp{
					entityID: "http://localhost:8000/saml/metadata",
					metadata: "PEVudGl0eURlc2NyaXB0b3IgeG1sbnM9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDptZXRhZGF0YSIgdmFsaWRVbnRpbD0iMjAyMi0wNC0yOFQxMTozMjowNC43OTdaIiBlbnRpdHlJRD0iaHR0cDovL2xvY2FsaG9zdDo4MDAwL3NhbWwvbWV0YWRhdGEiPgogIDxTUFNTT0Rlc2NyaXB0b3IgeG1sbnM9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDptZXRhZGF0YSIgdmFsaWRVbnRpbD0iMjAyMi0wNC0yOFQxMTozMjowNC43OTY5MjNaIiBwcm90b2NvbFN1cHBvcnRFbnVtZXJhdGlvbj0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOnByb3RvY29sIiBBdXRoblJlcXVlc3RzU2lnbmVkPSJ0cnVlIiBXYW50QXNzZXJ0aW9uc1NpZ25lZD0idHJ1ZSI+CiAgICA8S2V5RGVzY3JpcHRvciB1c2U9ImVuY3J5cHRpb24iPgogICAgICA8S2V5SW5mbyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC8wOS94bWxkc2lnIyI+CiAgICAgICAgPFg1MDlEYXRhIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwLzA5L3htbGRzaWcjIj4KICAgICAgICAgIDxYNTA5Q2VydGlmaWNhdGUgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyMiPk1JSUN2RENDQWFRQ0NRRDZFOFpHc1EydXNqQU5CZ2txaGtpRzl3MEJBUXNGQURBZ01SNHdIQVlEVlFRRERCVnRlWE5sY25acFkyVXVaWGhoYlhCc1pTNWpiMjB3SGhjTk1qSXdNakUzTVRRd05qTTVXaGNOTWpNd01qRTNNVFF3TmpNNVdqQWdNUjR3SEFZRFZRUUREQlZ0ZVhObGNuWnBZMlV1WlhoaGJYQnNaUzVqYjIwd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFS0FvSUJBUUM3WEtkQ1J4VVpYamRxVnF3d3dPSnFjMUNoMG5PU21rK1VlcmtVcWx2aVdIZGVMUitGb2xIS2pxTHpDQmxvQXo0eFZjMERGZlI3NmdXY1dBSEpsb3FaN0dCUzdOcERoelY4RytjWFErYlRVMEx1MmU3M3pDUWIzMFhVZEtoV2lHZkRLYVUrMXhnOUNELzJnSWZzWVBzM1RUcTFzcTdvQ3M1cUxkVUhhVkw1a2NSYUhLZG5UaTdjczVpOXh6czNUc1VuWGNySlB3eWRqcCthRWt5UmgwN29NcFhCRW9iR2lzZkYycDFNQTZwVlcyZ2pteXdmN0Q1aVlFRkVMUWhNN3BvcVBOMy9rZkJ2VTFuN0xmZ3E3b3htdi84TEZpNFpvcHI1bnlxc3oyNlhQdFV5MVdxVHpnem5BbVArbk4wb0JURVJGVmJYWGRSYTNrMnY0Y3hUTlBuL0FnTUJBQUV3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUpZeFJPV1NPWmJPelh6YWZkR2pRS3NNZ045NDhHL2hId1Z1Wm5leUFjVm9MTUZUczFXZXlhOVorc25NcDF1MEFkREdtUVRTOXpHbkQ3c3lEWUdPbWdpZ09MY012TE1vV2Y1dENRQmJFdWtXOE83RFBqUlIwWHlwQ2hHU3NIc3FMR08wQjBIYVRlbDBIZFA5U2k4MjdPQ2tjOVErV2JzRkcvOC80VG9HV0wrdWxhMVd1TGF3b3pvajh1bVBpOUQ4aVhDb1czNXkyU1RVK1dGUUc3VytLZmR1KzJDWXovMHRHZHdWcU5HNFdzZmF3V2NoclMwMHZHRktqbS9mSmM4NzZnQWZ4aU1IMUk5Zlp2WVNBeEFaM3NWSS8vTWwyc1VkZ2YwNjd5d1E3NW9hTFNTMk5JbW16NWFvczN2dVdtT1hoSUxkN2lUVStCRDhVdjZ2V2JJN0kxTT08L1g1MDlDZXJ0aWZpY2F0ZT4KICAgICAgICA8L1g1MDlEYXRhPgogICAgICA8L0tleUluZm8+CiAgICAgIDxFbmNyeXB0aW9uTWV0aG9kIEFsZ29yaXRobT0iaHR0cDovL3d3dy53My5vcmcvMjAwMS8wNC94bWxlbmMjYWVzMTI4LWNiYyI+PC9FbmNyeXB0aW9uTWV0aG9kPgogICAgICA8RW5jcnlwdGlvbk1ldGhvZCBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDEvMDQveG1sZW5jI2FlczE5Mi1jYmMiPjwvRW5jcnlwdGlvbk1ldGhvZD4KICAgICAgPEVuY3J5cHRpb25NZXRob2QgQWxnb3JpdGhtPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxLzA0L3htbGVuYyNhZXMyNTYtY2JjIj48L0VuY3J5cHRpb25NZXRob2Q+CiAgICAgIDxFbmNyeXB0aW9uTWV0aG9kIEFsZ29yaXRobT0iaHR0cDovL3d3dy53My5vcmcvMjAwMS8wNC94bWxlbmMjcnNhLW9hZXAtbWdmMXAiPjwvRW5jcnlwdGlvbk1ldGhvZD4KICAgIDwvS2V5RGVzY3JpcHRvcj4KICAgIDxLZXlEZXNjcmlwdG9yIHVzZT0ic2lnbmluZyI+CiAgICAgIDxLZXlJbmZvIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwLzA5L3htbGRzaWcjIj4KICAgICAgICA8WDUwOURhdGEgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyMiPgogICAgICAgICAgPFg1MDlDZXJ0aWZpY2F0ZSB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC8wOS94bWxkc2lnIyI+TUlJQ3ZEQ0NBYVFDQ1FENkU4WkdzUTJ1c2pBTkJna3Foa2lHOXcwQkFRc0ZBREFnTVI0d0hBWURWUVFEREJWdGVYTmxjblpwWTJVdVpYaGhiWEJzWlM1amIyMHdIaGNOTWpJd01qRTNNVFF3TmpNNVdoY05Nak13TWpFM01UUXdOak01V2pBZ01SNHdIQVlEVlFRRERCVnRlWE5sY25acFkyVXVaWGhoYlhCc1pTNWpiMjB3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRQzdYS2RDUnhVWlhqZHFWcXd3d09KcWMxQ2gwbk9TbWsrVWVya1VxbHZpV0hkZUxSK0ZvbEhLanFMekNCbG9BejR4VmMwREZmUjc2Z1djV0FISmxvcVo3R0JTN05wRGh6VjhHK2NYUStiVFUwTHUyZTczekNRYjMwWFVkS2hXaUdmREthVSsxeGc5Q0QvMmdJZnNZUHMzVFRxMXNxN29DczVxTGRVSGFWTDVrY1JhSEtkblRpN2NzNWk5eHpzM1RzVW5YY3JKUHd5ZGpwK2FFa3lSaDA3b01wWEJFb2JHaXNmRjJwMU1BNnBWVzJnam15d2Y3RDVpWUVGRUxRaE03cG9xUE4zL2tmQnZVMW43TGZncTdveG12LzhMRmk0Wm9wcjVueXFzejI2WFB0VXkxV3FUemd6bkFtUCtuTjBvQlRFUkZWYlhYZFJhM2sydjRjeFROUG4vQWdNQkFBRXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSll4Uk9XU09aYk96WHphZmRHalFLc01nTjk0OEcvaEh3VnVabmV5QWNWb0xNRlRzMVdleWE5Witzbk1wMXUwQWRER21RVFM5ekduRDdzeURZR09tZ2lnT0xjTXZMTW9XZjV0Q1FCYkV1a1c4TzdEUGpSUjBYeXBDaEdTc0hzcUxHTzBCMEhhVGVsMEhkUDlTaTgyN09Da2M5UStXYnNGRy84LzRUb0dXTCt1bGExV3VMYXdvem9qOHVtUGk5RDhpWENvVzM1eTJTVFUrV0ZRRzdXK0tmZHUrMkNZei8wdEdkd1ZxTkc0V3NmYXdXY2hyUzAwdkdGS2ptL2ZKYzg3NmdBZnhpTUgxSTlmWnZZU0F4QVozc1ZJLy9NbDJzVWRnZjA2N3l3UTc1b2FMU1MyTkltbXo1YW9zM3Z1V21PWGhJTGQ3aVRVK0JEOFV2NnZXYkk3STFNPTwvWDUwOUNlcnRpZmljYXRlPgogICAgICAgIDwvWDUwOURhdGE+CiAgICAgIDwvS2V5SW5mbz4KICAgIDwvS2V5RGVzY3JpcHRvcj4KICAgIDxTaW5nbGVMb2dvdXRTZXJ2aWNlIEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLVBPU1QiIExvY2F0aW9uPSJodHRwOi8vbG9jYWxob3N0OjgwMDAvc2FtbC9zbG8iIFJlc3BvbnNlTG9jYXRpb249Imh0dHA6Ly9sb2NhbGhvc3Q6ODAwMC9zYW1sL3NsbyI+PC9TaW5nbGVMb2dvdXRTZXJ2aWNlPgogICAgPEFzc2VydGlvbkNvbnN1bWVyU2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1QT1NUIiBMb2NhdGlvbj0iaHR0cDovL2xvY2FsaG9zdDo4MDAwL3NhbWwvYWNzIiBpbmRleD0iMSI+PC9Bc3NlcnRpb25Db25zdW1lclNlcnZpY2U+CiAgICA8QXNzZXJ0aW9uQ29uc3VtZXJTZXJ2aWNlIEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLUFydGlmYWN0IiBMb2NhdGlvbj0iaHR0cDovL2xvY2FsaG9zdDo4MDAwL3NhbWwvYWNzIiBpbmRleD0iMiI+PC9Bc3NlcnRpb25Db25zdW1lclNlcnZpY2U+CiAgPC9TUFNTT0Rlc2NyaXB0b3I+CjwvRW50aXR5RGVzY3JpcHRvcj4=",
				},
			},
			res{
				code: 303,
				err:  false,
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint := op.NewEndpoint(tt.args.metadataEndpoint)
			spInst, err := serviceprovider.NewServiceProvider(tt.args.sp.entityID, &serviceprovider.ServiceProviderConfig{Metadata: tt.args.sp.metadata}, "")
			if err != nil {
				t.Errorf("error while creating service provider")
				return
			}

			mockStorage := idpStorageWithResponseCertAndSP(
				t,
				[]byte(tt.args.certificate),
				[]byte(tt.args.key),
				"http://localhost:8000/saml/metadata",
				spInst,
				tt.args.request.ID,
			)
			if mockStorage == nil {
				return
			}

			idp, err := NewIdentityProvider(&endpoint, tt.args.config, mockStorage)
			if (err != nil) != tt.res.err {
				t.Errorf("NewIdentityProvider() error = %v", err.Error())
				return
			}
			if idp == nil {
				return
			}

			callURL := idp.SingleSignOnEndpoint.Relative()
			form := url.Values{}
			var req *http.Request
			if tt.args.request.Binding == RedirectBinding {
				callURL += "?SAMLRequest=" + url.QueryEscape(tt.args.request.SAMLRequest)
				if tt.args.request.RelayState != "" {
					callURL += "&RelayState=" + url.QueryEscape(tt.args.request.RelayState)
				}
				if tt.args.request.SigAlg != "" {
					callURL += "&SigAlg=" + url.QueryEscape(tt.args.request.SigAlg)
				}
				if tt.args.request.Signature != "" {
					callURL += "&Signature=" + url.QueryEscape(tt.args.request.Signature)
				}
				if tt.args.request.SAMLEncoding != "" {
					callURL += "&SAMLEncoding=" + url.QueryEscape(tt.args.request.SAMLEncoding)
				}
				req = httptest.NewRequest(http.MethodGet, callURL, nil)
			} else if tt.args.request.Binding == PostBinding {
				req = httptest.NewRequest(http.MethodPost, callURL, nil)
				form.Add("SAMLRequest", tt.args.request.SAMLRequest)
				form.Add("RelayState", tt.args.request.RelayState)
				req.Form = form
			}

			w := httptest.NewRecorder()

			idp.ssoHandleFunc(w, req)

			res := w.Result()
			defer func() {
				_ = res.Body.Close()
			}()
			_, err = ioutil.ReadAll(res.Body)
			if res.StatusCode != tt.res.code {
				t.Errorf("ssoHandleFunc() code got = %v, want %v", res.StatusCode, tt.res)
			}
		})
	}
}

func idpStorageWithResponseCertAndSP(
	t *testing.T,
	cert []byte,
	pKey []byte,
	entityID string,
	sp *serviceprovider.ServiceProvider,
	authRequestID string,
) *mock.MockIDPStorage {
	mockStorage := idpStorageWithResponseCert(t, cert, pKey)
	mockStorage.EXPECT().GetEntityByID(gomock.Any(), entityID).Return(sp, nil)

	request := mock.NewMockAuthRequestInt(gomock.NewController(t))
	request.EXPECT().GetID().Return(authRequestID)
	mockStorage.EXPECT().CreateAuthRequest(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(request, nil)

	return mockStorage
}
