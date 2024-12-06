package saml

import (
	"encoding/xml"
	"testing"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp/providers/saml/requesttracker"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestProvider_Options(t *testing.T) {
	type fields struct {
		name        string
		rootURL     string
		metadata    []byte
		key         []byte
		certificate []byte
		options     []ProviderOpts
	}
	type want struct {
		err                           bool
		name                          string
		linkingAllowed                bool
		creationAllowed               bool
		autoCreation                  bool
		autoUpdate                    bool
		binding                       string
		nameIDFormat                  saml.NameIDFormat
		transientMappingAttributeName string
		withSignedRequest             bool
		requesttracker                samlsp.RequestTracker
		entityID                      string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{{
		name: "failed metadata",
		fields: fields{
			name:     "saml",
			rootURL:  "https://localhost:8080",
			metadata: []byte(">xml<"),
			options:  nil,
		},
		want: want{
			err: true,
		},
	},
		{
			name: "failed keypair cert",
			fields: fields{
				name:     "saml",
				rootURL:  "https://localhost:8080",
				key:      []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAxHd087RoEm9ywVWZ/H+tDWxQsmVvhfRz4jAq/RfU+OWXNH4J\njMMSHdFs0Q+WP98nNXRyc7fgbMb8NdmlB2yD4qLYapN5SDaBc5dh/3EnyFt53oSs\njTlKnQUPAeJr2qh/NY046CfyUyQMM4JR5OiQFo4TssfWnqdcgamGt0AEnk2lvbMZ\nKQdAqNS9lDzYbjMGavEQPTZE35mFXFQXjaooZXq+TIa7hbaq7/idH7cHNbLcPLgj\nfPQA8q+DYvnvhXlmq0LPQZH3Oiixf+SF2vRwrBzT2mqGD2OiOkUmhuPwyqEiiBHt\nfxklRtRU6WfLa1Gcb1PsV0uoBGpV3KybIl/GlwIDAQABAoIBAEQjDduLgOCL6Gem\n0X3hpdnW6/HC/jed/Sa//9jBECq2LYeWAqff64ON40hqOHi0YvvGA/+gEOSI6mWe\nsv5tIxxRz+6+cLybsq+tG96kluCE4TJMHy/nY7orS/YiWbd+4odnEApr+D3fbZ/b\nnZ1fDsHTyn8hkYx6jLmnWsJpIHDp7zxD76y7k2Bbg6DZrCGiVxngiLJk23dvz79W\np03lHLM7XE92aFwXQmhfxHGxrbuoB/9eY4ai5IHp36H4fw0vL6NXdNQAo/bhe0p9\nAYB7y0ZumF8Hg0Z/BmMeEzLy6HrYB+VE8cO93pNjhSyH+p2yDB/BlUyTiRLQAoM0\nVTmOZXECgYEA7NGlzpKNhyQEJihVqt0MW0LhKIO/xbBn+XgYfX6GpqPa/ucnMx5/\nVezpl3gK8IU4wPUhAyXXAHJiqNBcEeyxrw0MXLujDVMJgYaLysCLJdvMVgoY08mS\nK5IQivpbozpf4+0y3mOnA+Sy1kbfxv2X8xiWLODRQW3f3q/xoklwOR8CgYEA1GEe\nfaibOFTQAYcIVj77KXtBfYZsX3EGAyfAN9O7cKHq5oaxVstwnF47WxpuVtoKZxCZ\nbNm9D5WvQ9b+Ztpioe42tzwE7Bff/Osj868GcDdRPK7nFlh9N2yVn/D514dOYVwR\n4MBr1KrJzgRWt4QqS4H+to1GzudDTSNlG7gnK4kCgYBUi6AbOHzoYzZL/RhgcJwp\ntJ23nhmH1Su5h2OO4e3mbhcP66w19sxU+8iFN+kH5zfUw26utgKk+TE5vXExQQRK\nT2k7bg2PAzcgk80ybD0BHhA8I0yrx4m0nmfjhe/TPVLgh10iwgbtP+eM0i6v1vc5\nZWyvxu9N4ZEL6lpkqr0y1wKBgG/NAIQd8jhhTW7Aav8cAJQBsqQl038avJOEpYe+\nCnpsgoAAf/K0/f8TDCQVceh+t+MxtdK7fO9rWOxZjWsPo8Si5mLnUaAHoX4/OpnZ\nlYYVWMqdOEFnK+O1Yb7k2GFBdV2DXlX2dc1qavntBsls5ecB89id3pyk2aUN8Pf6\npYQhAoGAMGtrHFely9wyaxI0RTCyfmJbWZHGVGkv6ELK8wneJjdjl82XOBUGCg5q\naRCrTZ3dPitKwrUa6ibJCIFCIziiriBmjDvTHzkMvoJEap2TVxYNDR6IfINVsQ57\nlOsiC4A2uGq4Lbfld+gjoplJ5GX6qXtTgZ6m7eo0y7U6zm2tkN0=\n-----END RSA PRIVATE KEY-----\n"),
				metadata: []byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-08-27T12:40:58.803Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
				options:  nil,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "failed keypair key",
			fields: fields{
				name:        "saml",
				rootURL:     "https://localhost:8080",
				certificate: []byte("-----BEGIN CERTIFICATE-----\nMIIC2zCCAcOgAwIBAgIIAy/jm1gAAdEwDQYJKoZIhvcNAQELBQAwEjEQMA4GA1UE\nChMHWklUQURFTDAeFw0yMzA4MzAwNzExMTVaFw0yNDA4MjkwNzExMTVaMBIxEDAO\nBgNVBAoTB1pJVEFERUwwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDE\nd3TztGgSb3LBVZn8f60NbFCyZW+F9HPiMCr9F9T45Zc0fgmMwxId0WzRD5Y/3yc1\ndHJzt+Bsxvw12aUHbIPiothqk3lINoFzl2H/cSfIW3nehKyNOUqdBQ8B4mvaqH81\njTjoJ/JTJAwzglHk6JAWjhOyx9aep1yBqYa3QASeTaW9sxkpB0Co1L2UPNhuMwZq\n8RA9NkTfmYVcVBeNqihler5MhruFtqrv+J0ftwc1stw8uCN89ADyr4Ni+e+FeWar\nQs9Bkfc6KLF/5IXa9HCsHNPaaoYPY6I6RSaG4/DKoSKIEe1/GSVG1FTpZ8trUZxv\nU+xXS6gEalXcrJsiX8aXAgMBAAGjNTAzMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUE\nDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMA0GCSqGSIb3DQEBCwUAA4IBAQCx\n/dRNIj0N/16zJhZR/ahkc2AkvDXYxyr4JRT5wK9GQDNl/oaX3debRuSi/tfaXFIX\naJA6PxM4J49ZaiEpLrKfxMz5kAhjKchCBEMcH3mGt+iNZH7EOyTvHjpGrP2OZrsh\nO17yrvN3HuQxIU6roJlqtZz2iAADsoPtwOO4D7hupm9XTMkSnAmlMWOo/q46Jz89\n1sMxB+dXmH/zV0wgwh0omZfLV0u89mvdq269VhcjNBpBYSnN1ccqYWd5iwziob3I\nvaavGHGfkbvRUn/tKftYuTK30q03R+e9YbmlWZ0v695owh2e/apCzowQsCKfSVC8\nOxVyt5XkHq1tWwVyBmFp\n-----END CERTIFICATE-----\n"),
				metadata:    []byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-08-27T12:40:58.803Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
				options:     nil,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "failed url",
			fields: fields{
				name:        "saml",
				rootURL:     "%%",
				key:         []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAxHd087RoEm9ywVWZ/H+tDWxQsmVvhfRz4jAq/RfU+OWXNH4J\njMMSHdFs0Q+WP98nNXRyc7fgbMb8NdmlB2yD4qLYapN5SDaBc5dh/3EnyFt53oSs\njTlKnQUPAeJr2qh/NY046CfyUyQMM4JR5OiQFo4TssfWnqdcgamGt0AEnk2lvbMZ\nKQdAqNS9lDzYbjMGavEQPTZE35mFXFQXjaooZXq+TIa7hbaq7/idH7cHNbLcPLgj\nfPQA8q+DYvnvhXlmq0LPQZH3Oiixf+SF2vRwrBzT2mqGD2OiOkUmhuPwyqEiiBHt\nfxklRtRU6WfLa1Gcb1PsV0uoBGpV3KybIl/GlwIDAQABAoIBAEQjDduLgOCL6Gem\n0X3hpdnW6/HC/jed/Sa//9jBECq2LYeWAqff64ON40hqOHi0YvvGA/+gEOSI6mWe\nsv5tIxxRz+6+cLybsq+tG96kluCE4TJMHy/nY7orS/YiWbd+4odnEApr+D3fbZ/b\nnZ1fDsHTyn8hkYx6jLmnWsJpIHDp7zxD76y7k2Bbg6DZrCGiVxngiLJk23dvz79W\np03lHLM7XE92aFwXQmhfxHGxrbuoB/9eY4ai5IHp36H4fw0vL6NXdNQAo/bhe0p9\nAYB7y0ZumF8Hg0Z/BmMeEzLy6HrYB+VE8cO93pNjhSyH+p2yDB/BlUyTiRLQAoM0\nVTmOZXECgYEA7NGlzpKNhyQEJihVqt0MW0LhKIO/xbBn+XgYfX6GpqPa/ucnMx5/\nVezpl3gK8IU4wPUhAyXXAHJiqNBcEeyxrw0MXLujDVMJgYaLysCLJdvMVgoY08mS\nK5IQivpbozpf4+0y3mOnA+Sy1kbfxv2X8xiWLODRQW3f3q/xoklwOR8CgYEA1GEe\nfaibOFTQAYcIVj77KXtBfYZsX3EGAyfAN9O7cKHq5oaxVstwnF47WxpuVtoKZxCZ\nbNm9D5WvQ9b+Ztpioe42tzwE7Bff/Osj868GcDdRPK7nFlh9N2yVn/D514dOYVwR\n4MBr1KrJzgRWt4QqS4H+to1GzudDTSNlG7gnK4kCgYBUi6AbOHzoYzZL/RhgcJwp\ntJ23nhmH1Su5h2OO4e3mbhcP66w19sxU+8iFN+kH5zfUw26utgKk+TE5vXExQQRK\nT2k7bg2PAzcgk80ybD0BHhA8I0yrx4m0nmfjhe/TPVLgh10iwgbtP+eM0i6v1vc5\nZWyvxu9N4ZEL6lpkqr0y1wKBgG/NAIQd8jhhTW7Aav8cAJQBsqQl038avJOEpYe+\nCnpsgoAAf/K0/f8TDCQVceh+t+MxtdK7fO9rWOxZjWsPo8Si5mLnUaAHoX4/OpnZ\nlYYVWMqdOEFnK+O1Yb7k2GFBdV2DXlX2dc1qavntBsls5ecB89id3pyk2aUN8Pf6\npYQhAoGAMGtrHFely9wyaxI0RTCyfmJbWZHGVGkv6ELK8wneJjdjl82XOBUGCg5q\naRCrTZ3dPitKwrUa6ibJCIFCIziiriBmjDvTHzkMvoJEap2TVxYNDR6IfINVsQ57\nlOsiC4A2uGq4Lbfld+gjoplJ5GX6qXtTgZ6m7eo0y7U6zm2tkN0=\n-----END RSA PRIVATE KEY-----\n"),
				certificate: []byte("-----BEGIN CERTIFICATE-----\nMIIC2zCCAcOgAwIBAgIIAy/jm1gAAdEwDQYJKoZIhvcNAQELBQAwEjEQMA4GA1UE\nChMHWklUQURFTDAeFw0yMzA4MzAwNzExMTVaFw0yNDA4MjkwNzExMTVaMBIxEDAO\nBgNVBAoTB1pJVEFERUwwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDE\nd3TztGgSb3LBVZn8f60NbFCyZW+F9HPiMCr9F9T45Zc0fgmMwxId0WzRD5Y/3yc1\ndHJzt+Bsxvw12aUHbIPiothqk3lINoFzl2H/cSfIW3nehKyNOUqdBQ8B4mvaqH81\njTjoJ/JTJAwzglHk6JAWjhOyx9aep1yBqYa3QASeTaW9sxkpB0Co1L2UPNhuMwZq\n8RA9NkTfmYVcVBeNqihler5MhruFtqrv+J0ftwc1stw8uCN89ADyr4Ni+e+FeWar\nQs9Bkfc6KLF/5IXa9HCsHNPaaoYPY6I6RSaG4/DKoSKIEe1/GSVG1FTpZ8trUZxv\nU+xXS6gEalXcrJsiX8aXAgMBAAGjNTAzMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUE\nDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMA0GCSqGSIb3DQEBCwUAA4IBAQCx\n/dRNIj0N/16zJhZR/ahkc2AkvDXYxyr4JRT5wK9GQDNl/oaX3debRuSi/tfaXFIX\naJA6PxM4J49ZaiEpLrKfxMz5kAhjKchCBEMcH3mGt+iNZH7EOyTvHjpGrP2OZrsh\nO17yrvN3HuQxIU6roJlqtZz2iAADsoPtwOO4D7hupm9XTMkSnAmlMWOo/q46Jz89\n1sMxB+dXmH/zV0wgwh0omZfLV0u89mvdq269VhcjNBpBYSnN1ccqYWd5iwziob3I\nvaavGHGfkbvRUn/tKftYuTK30q03R+e9YbmlWZ0v695owh2e/apCzowQsCKfSVC8\nOxVyt5XkHq1tWwVyBmFp\n-----END CERTIFICATE-----\n"),
				metadata:    []byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-08-27T12:40:58.803Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
				options:     nil,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "default",
			fields: fields{
				name:        "saml",
				rootURL:     "https://localhost:8080",
				key:         []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAxHd087RoEm9ywVWZ/H+tDWxQsmVvhfRz4jAq/RfU+OWXNH4J\njMMSHdFs0Q+WP98nNXRyc7fgbMb8NdmlB2yD4qLYapN5SDaBc5dh/3EnyFt53oSs\njTlKnQUPAeJr2qh/NY046CfyUyQMM4JR5OiQFo4TssfWnqdcgamGt0AEnk2lvbMZ\nKQdAqNS9lDzYbjMGavEQPTZE35mFXFQXjaooZXq+TIa7hbaq7/idH7cHNbLcPLgj\nfPQA8q+DYvnvhXlmq0LPQZH3Oiixf+SF2vRwrBzT2mqGD2OiOkUmhuPwyqEiiBHt\nfxklRtRU6WfLa1Gcb1PsV0uoBGpV3KybIl/GlwIDAQABAoIBAEQjDduLgOCL6Gem\n0X3hpdnW6/HC/jed/Sa//9jBECq2LYeWAqff64ON40hqOHi0YvvGA/+gEOSI6mWe\nsv5tIxxRz+6+cLybsq+tG96kluCE4TJMHy/nY7orS/YiWbd+4odnEApr+D3fbZ/b\nnZ1fDsHTyn8hkYx6jLmnWsJpIHDp7zxD76y7k2Bbg6DZrCGiVxngiLJk23dvz79W\np03lHLM7XE92aFwXQmhfxHGxrbuoB/9eY4ai5IHp36H4fw0vL6NXdNQAo/bhe0p9\nAYB7y0ZumF8Hg0Z/BmMeEzLy6HrYB+VE8cO93pNjhSyH+p2yDB/BlUyTiRLQAoM0\nVTmOZXECgYEA7NGlzpKNhyQEJihVqt0MW0LhKIO/xbBn+XgYfX6GpqPa/ucnMx5/\nVezpl3gK8IU4wPUhAyXXAHJiqNBcEeyxrw0MXLujDVMJgYaLysCLJdvMVgoY08mS\nK5IQivpbozpf4+0y3mOnA+Sy1kbfxv2X8xiWLODRQW3f3q/xoklwOR8CgYEA1GEe\nfaibOFTQAYcIVj77KXtBfYZsX3EGAyfAN9O7cKHq5oaxVstwnF47WxpuVtoKZxCZ\nbNm9D5WvQ9b+Ztpioe42tzwE7Bff/Osj868GcDdRPK7nFlh9N2yVn/D514dOYVwR\n4MBr1KrJzgRWt4QqS4H+to1GzudDTSNlG7gnK4kCgYBUi6AbOHzoYzZL/RhgcJwp\ntJ23nhmH1Su5h2OO4e3mbhcP66w19sxU+8iFN+kH5zfUw26utgKk+TE5vXExQQRK\nT2k7bg2PAzcgk80ybD0BHhA8I0yrx4m0nmfjhe/TPVLgh10iwgbtP+eM0i6v1vc5\nZWyvxu9N4ZEL6lpkqr0y1wKBgG/NAIQd8jhhTW7Aav8cAJQBsqQl038avJOEpYe+\nCnpsgoAAf/K0/f8TDCQVceh+t+MxtdK7fO9rWOxZjWsPo8Si5mLnUaAHoX4/OpnZ\nlYYVWMqdOEFnK+O1Yb7k2GFBdV2DXlX2dc1qavntBsls5ecB89id3pyk2aUN8Pf6\npYQhAoGAMGtrHFely9wyaxI0RTCyfmJbWZHGVGkv6ELK8wneJjdjl82XOBUGCg5q\naRCrTZ3dPitKwrUa6ibJCIFCIziiriBmjDvTHzkMvoJEap2TVxYNDR6IfINVsQ57\nlOsiC4A2uGq4Lbfld+gjoplJ5GX6qXtTgZ6m7eo0y7U6zm2tkN0=\n-----END RSA PRIVATE KEY-----\n"),
				certificate: []byte("-----BEGIN CERTIFICATE-----\nMIIC2zCCAcOgAwIBAgIIAy/jm1gAAdEwDQYJKoZIhvcNAQELBQAwEjEQMA4GA1UE\nChMHWklUQURFTDAeFw0yMzA4MzAwNzExMTVaFw0yNDA4MjkwNzExMTVaMBIxEDAO\nBgNVBAoTB1pJVEFERUwwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDE\nd3TztGgSb3LBVZn8f60NbFCyZW+F9HPiMCr9F9T45Zc0fgmMwxId0WzRD5Y/3yc1\ndHJzt+Bsxvw12aUHbIPiothqk3lINoFzl2H/cSfIW3nehKyNOUqdBQ8B4mvaqH81\njTjoJ/JTJAwzglHk6JAWjhOyx9aep1yBqYa3QASeTaW9sxkpB0Co1L2UPNhuMwZq\n8RA9NkTfmYVcVBeNqihler5MhruFtqrv+J0ftwc1stw8uCN89ADyr4Ni+e+FeWar\nQs9Bkfc6KLF/5IXa9HCsHNPaaoYPY6I6RSaG4/DKoSKIEe1/GSVG1FTpZ8trUZxv\nU+xXS6gEalXcrJsiX8aXAgMBAAGjNTAzMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUE\nDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMA0GCSqGSIb3DQEBCwUAA4IBAQCx\n/dRNIj0N/16zJhZR/ahkc2AkvDXYxyr4JRT5wK9GQDNl/oaX3debRuSi/tfaXFIX\naJA6PxM4J49ZaiEpLrKfxMz5kAhjKchCBEMcH3mGt+iNZH7EOyTvHjpGrP2OZrsh\nO17yrvN3HuQxIU6roJlqtZz2iAADsoPtwOO4D7hupm9XTMkSnAmlMWOo/q46Jz89\n1sMxB+dXmH/zV0wgwh0omZfLV0u89mvdq269VhcjNBpBYSnN1ccqYWd5iwziob3I\nvaavGHGfkbvRUn/tKftYuTK30q03R+e9YbmlWZ0v695owh2e/apCzowQsCKfSVC8\nOxVyt5XkHq1tWwVyBmFp\n-----END CERTIFICATE-----\n"),
				metadata:    []byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-08-27T12:40:58.803Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
				options:     nil,
			},
			want: want{
				name:            "saml",
				linkingAllowed:  false,
				creationAllowed: false,
				autoCreation:    false,
				autoUpdate:      false,
				nameIDFormat:    saml.PersistentNameIDFormat,
			},
		},
		{
			name: "all set / true",
			fields: fields{
				name:        "saml",
				key:         []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAxHd087RoEm9ywVWZ/H+tDWxQsmVvhfRz4jAq/RfU+OWXNH4J\njMMSHdFs0Q+WP98nNXRyc7fgbMb8NdmlB2yD4qLYapN5SDaBc5dh/3EnyFt53oSs\njTlKnQUPAeJr2qh/NY046CfyUyQMM4JR5OiQFo4TssfWnqdcgamGt0AEnk2lvbMZ\nKQdAqNS9lDzYbjMGavEQPTZE35mFXFQXjaooZXq+TIa7hbaq7/idH7cHNbLcPLgj\nfPQA8q+DYvnvhXlmq0LPQZH3Oiixf+SF2vRwrBzT2mqGD2OiOkUmhuPwyqEiiBHt\nfxklRtRU6WfLa1Gcb1PsV0uoBGpV3KybIl/GlwIDAQABAoIBAEQjDduLgOCL6Gem\n0X3hpdnW6/HC/jed/Sa//9jBECq2LYeWAqff64ON40hqOHi0YvvGA/+gEOSI6mWe\nsv5tIxxRz+6+cLybsq+tG96kluCE4TJMHy/nY7orS/YiWbd+4odnEApr+D3fbZ/b\nnZ1fDsHTyn8hkYx6jLmnWsJpIHDp7zxD76y7k2Bbg6DZrCGiVxngiLJk23dvz79W\np03lHLM7XE92aFwXQmhfxHGxrbuoB/9eY4ai5IHp36H4fw0vL6NXdNQAo/bhe0p9\nAYB7y0ZumF8Hg0Z/BmMeEzLy6HrYB+VE8cO93pNjhSyH+p2yDB/BlUyTiRLQAoM0\nVTmOZXECgYEA7NGlzpKNhyQEJihVqt0MW0LhKIO/xbBn+XgYfX6GpqPa/ucnMx5/\nVezpl3gK8IU4wPUhAyXXAHJiqNBcEeyxrw0MXLujDVMJgYaLysCLJdvMVgoY08mS\nK5IQivpbozpf4+0y3mOnA+Sy1kbfxv2X8xiWLODRQW3f3q/xoklwOR8CgYEA1GEe\nfaibOFTQAYcIVj77KXtBfYZsX3EGAyfAN9O7cKHq5oaxVstwnF47WxpuVtoKZxCZ\nbNm9D5WvQ9b+Ztpioe42tzwE7Bff/Osj868GcDdRPK7nFlh9N2yVn/D514dOYVwR\n4MBr1KrJzgRWt4QqS4H+to1GzudDTSNlG7gnK4kCgYBUi6AbOHzoYzZL/RhgcJwp\ntJ23nhmH1Su5h2OO4e3mbhcP66w19sxU+8iFN+kH5zfUw26utgKk+TE5vXExQQRK\nT2k7bg2PAzcgk80ybD0BHhA8I0yrx4m0nmfjhe/TPVLgh10iwgbtP+eM0i6v1vc5\nZWyvxu9N4ZEL6lpkqr0y1wKBgG/NAIQd8jhhTW7Aav8cAJQBsqQl038avJOEpYe+\nCnpsgoAAf/K0/f8TDCQVceh+t+MxtdK7fO9rWOxZjWsPo8Si5mLnUaAHoX4/OpnZ\nlYYVWMqdOEFnK+O1Yb7k2GFBdV2DXlX2dc1qavntBsls5ecB89id3pyk2aUN8Pf6\npYQhAoGAMGtrHFely9wyaxI0RTCyfmJbWZHGVGkv6ELK8wneJjdjl82XOBUGCg5q\naRCrTZ3dPitKwrUa6ibJCIFCIziiriBmjDvTHzkMvoJEap2TVxYNDR6IfINVsQ57\nlOsiC4A2uGq4Lbfld+gjoplJ5GX6qXtTgZ6m7eo0y7U6zm2tkN0=\n-----END RSA PRIVATE KEY-----\n"),
				certificate: []byte("-----BEGIN CERTIFICATE-----\nMIIC2zCCAcOgAwIBAgIIAy/jm1gAAdEwDQYJKoZIhvcNAQELBQAwEjEQMA4GA1UE\nChMHWklUQURFTDAeFw0yMzA4MzAwNzExMTVaFw0yNDA4MjkwNzExMTVaMBIxEDAO\nBgNVBAoTB1pJVEFERUwwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDE\nd3TztGgSb3LBVZn8f60NbFCyZW+F9HPiMCr9F9T45Zc0fgmMwxId0WzRD5Y/3yc1\ndHJzt+Bsxvw12aUHbIPiothqk3lINoFzl2H/cSfIW3nehKyNOUqdBQ8B4mvaqH81\njTjoJ/JTJAwzglHk6JAWjhOyx9aep1yBqYa3QASeTaW9sxkpB0Co1L2UPNhuMwZq\n8RA9NkTfmYVcVBeNqihler5MhruFtqrv+J0ftwc1stw8uCN89ADyr4Ni+e+FeWar\nQs9Bkfc6KLF/5IXa9HCsHNPaaoYPY6I6RSaG4/DKoSKIEe1/GSVG1FTpZ8trUZxv\nU+xXS6gEalXcrJsiX8aXAgMBAAGjNTAzMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUE\nDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMA0GCSqGSIb3DQEBCwUAA4IBAQCx\n/dRNIj0N/16zJhZR/ahkc2AkvDXYxyr4JRT5wK9GQDNl/oaX3debRuSi/tfaXFIX\naJA6PxM4J49ZaiEpLrKfxMz5kAhjKchCBEMcH3mGt+iNZH7EOyTvHjpGrP2OZrsh\nO17yrvN3HuQxIU6roJlqtZz2iAADsoPtwOO4D7hupm9XTMkSnAmlMWOo/q46Jz89\n1sMxB+dXmH/zV0wgwh0omZfLV0u89mvdq269VhcjNBpBYSnN1ccqYWd5iwziob3I\nvaavGHGfkbvRUn/tKftYuTK30q03R+e9YbmlWZ0v695owh2e/apCzowQsCKfSVC8\nOxVyt5XkHq1tWwVyBmFp\n-----END CERTIFICATE-----\n"),
				metadata:    []byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-08-27T12:40:58.803Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
				options: []ProviderOpts{
					WithLinkingAllowed(),
					WithCreationAllowed(),
					WithAutoCreation(),
					WithAutoUpdate(),
					WithBinding("binding"),
					WithSignedRequest(),
					WithCustomRequestTracker(&requesttracker.RequestTracker{}),
					WithEntityID("entityID"),
					WithNameIDFormat(domain.SAMLNameIDFormatTransient),
					WithTransientMappingAttributeName("attribute"),
				},
			},
			want: want{
				name:                          "saml",
				linkingAllowed:                true,
				creationAllowed:               true,
				autoCreation:                  true,
				autoUpdate:                    true,
				binding:                       "binding",
				entityID:                      "entityID",
				nameIDFormat:                  saml.TransientNameIDFormat,
				transientMappingAttributeName: "attribute",
				withSignedRequest:             true,
				requesttracker:                &requesttracker.RequestTracker{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			provider, err := New(tt.fields.name, tt.fields.rootURL, tt.fields.metadata, tt.fields.certificate, tt.fields.key, tt.fields.options...)
			if tt.want.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				a.Equal(tt.want.name, provider.Name())
				a.Equal(tt.want.linkingAllowed, provider.IsLinkingAllowed())
				a.Equal(tt.want.creationAllowed, provider.IsCreationAllowed())
				a.Equal(tt.want.autoCreation, provider.IsAutoCreation())
				a.Equal(tt.want.autoUpdate, provider.IsAutoUpdate())
				a.Equal(tt.want.binding, provider.binding)
				a.Equal(tt.want.nameIDFormat, provider.nameIDFormat)
				a.Equal(tt.want.transientMappingAttributeName, provider.transientMappingAttributeName)
				a.Equal(tt.want.withSignedRequest, provider.spOptions.SignRequest)
				a.Equal(tt.want.requesttracker, provider.requestTracker)
				a.Equal(tt.want.entityID, provider.spOptions.EntityID)
			}
		})
	}
}

func TestParseMetadata(t *testing.T) {
	type args struct {
		metadata []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *saml.EntityDescriptor
		wantErr error
	}{
		{
			"invalid",
			args{
				metadata: []byte(`<Test></Test>`),
			},
			nil,
			xml.UnmarshalError("expected element type <EntityDescriptor> but have <Test>"),
		},
		{
			"valid entity descriptor",
			args{
				metadata: []byte(`<?xml version="1.0" encoding="UTF-8"?><EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" entityID="http://localhost:8000/metadata"><IDPSSODescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata"><SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="http://localhost:8000/sso"></SingleSignOnService></IDPSSODescriptor></EntityDescriptor>`),
			},
			&saml.EntityDescriptor{
				EntityID: "http://localhost:8000/metadata",
				IDPSSODescriptors: []saml.IDPSSODescriptor{
					{
						XMLName: xml.Name{
							Space: "urn:oasis:names:tc:SAML:2.0:metadata",
							Local: "IDPSSODescriptor",
						},
						SingleSignOnServices: []saml.Endpoint{
							{
								Binding:  saml.HTTPRedirectBinding,
								Location: "http://localhost:8000/sso",
							},
						},
					},
				},
			},
			nil,
		},
		{
			"valid entity descriptor, non utf-8",
			args{
				metadata: []byte(`<?xml version="1.0" encoding="windows-1252"?><EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" entityID="http://localhost:8000/metadata"><IDPSSODescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata"><SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="http://localhost:8000/sso"></SingleSignOnService></IDPSSODescriptor></EntityDescriptor>`),
			},
			&saml.EntityDescriptor{
				EntityID: "http://localhost:8000/metadata",
				IDPSSODescriptors: []saml.IDPSSODescriptor{
					{
						XMLName: xml.Name{
							Space: "urn:oasis:names:tc:SAML:2.0:metadata",
							Local: "IDPSSODescriptor",
						},
						SingleSignOnServices: []saml.Endpoint{
							{
								Binding:  saml.HTTPRedirectBinding,
								Location: "http://localhost:8000/sso",
							},
						},
					},
				},
			},
			nil,
		},
		{
			"entities descriptor without IDPSSODescriptor",
			args{
				metadata: []byte(`<?xml version="1.0" encoding="UTF-8"?><EntitiesDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata"><EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" entityID="http://localhost:8000/metadata"></EntityDescriptor></EntitiesDescriptor>`),
			},
			nil,
			zerrors.ThrowInternal(nil, "SAML-Ejoi3r2", "no entity found with IDPSSODescriptor"),
		},
		{
			"valid entities descriptor",
			args{
				metadata: []byte(`<?xml version="1.0" encoding="UTF-8"?><EntitiesDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata"><EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" entityID="http://localhost:8000/metadata"><IDPSSODescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata"><SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="http://localhost:8000/sso"></SingleSignOnService></IDPSSODescriptor></EntityDescriptor></EntitiesDescriptor>`),
			},
			&saml.EntityDescriptor{
				EntityID: "http://localhost:8000/metadata",
				IDPSSODescriptors: []saml.IDPSSODescriptor{
					{
						XMLName: xml.Name{
							Space: "urn:oasis:names:tc:SAML:2.0:metadata",
							Local: "IDPSSODescriptor",
						},
						SingleSignOnServices: []saml.Endpoint{
							{
								Binding:  saml.HTTPRedirectBinding,
								Location: "http://localhost:8000/sso",
							},
						},
					},
				},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMetadata(tt.args.metadata)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
