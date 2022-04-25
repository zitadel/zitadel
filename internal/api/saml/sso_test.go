package saml

import (
	"github.com/caos/zitadel/internal/api/saml/serviceprovider"
	"github.com/caos/zitadel/internal/api/saml/xml/md"
	"github.com/caos/zitadel/internal/api/saml/xml/xml_dsig"
	"net/http"
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
