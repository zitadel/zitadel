package signature_test

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"github.com/beevik/etree"
	"github.com/caos/zitadel/internal/api/saml/signature"
	saml_xml "github.com/caos/zitadel/internal/api/saml/xml"
	"github.com/caos/zitadel/internal/api/saml/xml/samlp"
	"github.com/caos/zitadel/internal/api/saml/xml/xml_dsig"
	"github.com/caos/zitadel/internal/crypto"
	dsig "github.com/russellhaering/goxmldsig"
	"math/big"
	"math/rand"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestSignature_ValidatePost(t *testing.T) {
	type args struct {
		authRequest string
		certificate string
	}

	tests := []struct {
		name string
		args args
		err  bool
	}{
		{
			name: "crewjam",
			args: args{
				authRequest: "<samlp:AuthnRequest xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" ID=\"id-6594153690cfdc336afa98e39145e9716bb88e0f\" Version=\"2.0\" IssueInstant=\"2022-04-13T09:45:38.162Z\" Destination=\"http://localhost:50002/saml/SSO\" AssertionConsumerServiceURL=\"http://localhost:8000/saml/acs\" ProtocolBinding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"><saml:Issuer Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:entity\">http://localhost:8000/saml/metadata</saml:Issuer><ds:Signature xmlns:ds=\"http://www.w3.org/2000/09/xmldsig#\"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/><ds:SignatureMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#rsa-sha1\"/><ds:Reference URI=\"#id-6594153690cfdc336afa98e39145e9716bb88e0f\"><ds:Transforms><ds:Transform Algorithm=\"http://www.w3.org/2000/09/xmldsig#enveloped-signature\"/><ds:Transform Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/></ds:Transforms><ds:DigestMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#sha1\"/><ds:DigestValue>uJqDdJyOLqWh7H0JPGuH9+/FbcY=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>XQLDINOWHYaGtt53WetQ7G3twWjkMAR+BKoO91LQWTYsv7NxxvVN9F1PyMxEJQyBnfMWnbckokIF0F6jieNNByWm6yX7EG9WOvc9ymb6e7Cx+zElWdvBKbO0Is3cQWiDTn3teexZjcVD1fAdHgpiTA1vYPYaR40pVwVbTQ+qSf7SuYHQifoDzGFIrPGIDZ9PkQ995fQbo8rdLwmJxAVjkMd/24z3iC3h5h6WqVT/9xcGSr8Rr8SQ9rKLZBgpwz2me0PcT97NHHhHR8K6Am0cXTFaKhfZxQQB601RgO6/38BoOtXdIXh+veBIHBlbNagEkJfnzvBNQ6KAYQ58EnVC3w==</ds:SignatureValue><ds:KeyInfo><ds:X509Data><ds:X509Certificate>MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQwNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVoLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+ula1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4WsfawWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oaLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=</ds:X509Certificate></ds:X509Data></ds:KeyInfo></ds:Signature><samlp:NameIDPolicy Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\" AllowCreate=\"true\"/></samlp:AuthnRequest>",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			err: false,
		},
		{
			name: "crewjam without keyinfo",
			args: args{
				authRequest: "<samlp:AuthnRequest xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" ID=\"id-6594153690cfdc336afa98e39145e9716bb88e0f\" Version=\"2.0\" IssueInstant=\"2022-04-13T09:45:38.162Z\" Destination=\"http://localhost:50002/saml/SSO\" AssertionConsumerServiceURL=\"http://localhost:8000/saml/acs\" ProtocolBinding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"><saml:Issuer Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:entity\">http://localhost:8000/saml/metadata</saml:Issuer><ds:Signature xmlns:ds=\"http://www.w3.org/2000/09/xmldsig#\"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/><ds:SignatureMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#rsa-sha1\"/><ds:Reference URI=\"#id-6594153690cfdc336afa98e39145e9716bb88e0f\"><ds:Transforms><ds:Transform Algorithm=\"http://www.w3.org/2000/09/xmldsig#enveloped-signature\"/><ds:Transform Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/></ds:Transforms><ds:DigestMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#sha1\"/><ds:DigestValue>uJqDdJyOLqWh7H0JPGuH9+/FbcY=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>XQLDINOWHYaGtt53WetQ7G3twWjkMAR+BKoO91LQWTYsv7NxxvVN9F1PyMxEJQyBnfMWnbckokIF0F6jieNNByWm6yX7EG9WOvc9ymb6e7Cx+zElWdvBKbO0Is3cQWiDTn3teexZjcVD1fAdHgpiTA1vYPYaR40pVwVbTQ+qSf7SuYHQifoDzGFIrPGIDZ9PkQ995fQbo8rdLwmJxAVjkMd/24z3iC3h5h6WqVT/9xcGSr8Rr8SQ9rKLZBgpwz2me0PcT97NHHhHR8K6Am0cXTFaKhfZxQQB601RgO6/38BoOtXdIXh+veBIHBlbNagEkJfnzvBNQ6KAYQ58EnVC3w==</ds:SignatureValue></ds:Signature><samlp:NameIDPolicy Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\" AllowCreate=\"true\"/></samlp:AuthnRequest>",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			err: false,
		},
		{
			name: "crewjam without ID",
			args: args{
				authRequest: "<samlp:AuthnRequest xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" Version=\"2.0\" IssueInstant=\"2022-04-13T09:45:38.162Z\" Destination=\"http://localhost:50002/saml/SSO\" AssertionConsumerServiceURL=\"http://localhost:8000/saml/acs\" ProtocolBinding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"><saml:Issuer Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:entity\">http://localhost:8000/saml/metadata</saml:Issuer><ds:Signature xmlns:ds=\"http://www.w3.org/2000/09/xmldsig#\"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/><ds:SignatureMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#rsa-sha1\"/><ds:Reference URI=\"#id-6594153690cfdc336afa98e39145e9716bb88e0f\"><ds:Transforms><ds:Transform Algorithm=\"http://www.w3.org/2000/09/xmldsig#enveloped-signature\"/><ds:Transform Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/></ds:Transforms><ds:DigestMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#sha1\"/><ds:DigestValue>uJqDdJyOLqWh7H0JPGuH9+/FbcY=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>XQLDINOWHYaGtt53WetQ7G3twWjkMAR+BKoO91LQWTYsv7NxxvVN9F1PyMxEJQyBnfMWnbckokIF0F6jieNNByWm6yX7EG9WOvc9ymb6e7Cx+zElWdvBKbO0Is3cQWiDTn3teexZjcVD1fAdHgpiTA1vYPYaR40pVwVbTQ+qSf7SuYHQifoDzGFIrPGIDZ9PkQ995fQbo8rdLwmJxAVjkMd/24z3iC3h5h6WqVT/9xcGSr8Rr8SQ9rKLZBgpwz2me0PcT97NHHhHR8K6Am0cXTFaKhfZxQQB601RgO6/38BoOtXdIXh+veBIHBlbNagEkJfnzvBNQ6KAYQ58EnVC3w==</ds:SignatureValue><ds:KeyInfo><ds:X509Data><ds:X509Certificate>MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQwNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVoLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+ula1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4WsfawWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oaLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=</ds:X509Certificate></ds:X509Data></ds:KeyInfo></ds:Signature><samlp:NameIDPolicy Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\" AllowCreate=\"true\"/></samlp:AuthnRequest>",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			err: true,
		},
		{
			name: "crewjam without signature",
			args: args{
				authRequest: "<samlp:AuthnRequest xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" ID=\"id-6594153690cfdc336afa98e39145e9716bb88e0f\" Version=\"2.0\" IssueInstant=\"2022-04-13T09:45:38.162Z\" Destination=\"http://localhost:50002/saml/SSO\" AssertionConsumerServiceURL=\"http://localhost:8000/saml/acs\" ProtocolBinding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"><saml:Issuer Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:entity\">http://localhost:8000/saml/metadata</saml:Issuer><samlp:NameIDPolicy Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\" AllowCreate=\"true\"/></samlp:AuthnRequest>",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			err: true,
		},
		{
			name: "crewjam without digestValue",
			args: args{
				authRequest: "<samlp:AuthnRequest xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" ID=\"id-6594153690cfdc336afa98e39145e9716bb88e0f\" Version=\"2.0\" IssueInstant=\"2022-04-13T09:45:38.162Z\" Destination=\"http://localhost:50002/saml/SSO\" AssertionConsumerServiceURL=\"http://localhost:8000/saml/acs\" ProtocolBinding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"><saml:Issuer Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:entity\">http://localhost:8000/saml/metadata</saml:Issuer><ds:Signature xmlns:ds=\"http://www.w3.org/2000/09/xmldsig#\"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/><ds:SignatureMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#rsa-sha1\"/><ds:Reference URI=\"#id-6594153690cfdc336afa98e39145e9716bb88e0f\"><ds:Transforms><ds:Transform Algorithm=\"http://www.w3.org/2000/09/xmldsig#enveloped-signature\"/><ds:Transform Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/></ds:Transforms><ds:DigestMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#sha1\"/></ds:Reference></ds:SignedInfo><ds:SignatureValue>XQLDINOWHYaGtt53WetQ7G3twWjkMAR+BKoO91LQWTYsv7NxxvVN9F1PyMxEJQyBnfMWnbckokIF0F6jieNNByWm6yX7EG9WOvc9ymb6e7Cx+zElWdvBKbO0Is3cQWiDTn3teexZjcVD1fAdHgpiTA1vYPYaR40pVwVbTQ+qSf7SuYHQifoDzGFIrPGIDZ9PkQ995fQbo8rdLwmJxAVjkMd/24z3iC3h5h6WqVT/9xcGSr8Rr8SQ9rKLZBgpwz2me0PcT97NHHhHR8K6Am0cXTFaKhfZxQQB601RgO6/38BoOtXdIXh+veBIHBlbNagEkJfnzvBNQ6KAYQ58EnVC3w==</ds:SignatureValue><ds:KeyInfo><ds:X509Data><ds:X509Certificate>MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQwNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVoLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+ula1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4WsfawWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oaLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=</ds:X509Certificate></ds:X509Data></ds:KeyInfo></ds:Signature><samlp:NameIDPolicy Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\" AllowCreate=\"true\"/></samlp:AuthnRequest>",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			err: true,
		},
		{
			name: "crewjam without signatureValue",
			args: args{
				authRequest: "<samlp:AuthnRequest xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" ID=\"id-6594153690cfdc336afa98e39145e9716bb88e0f\" Version=\"2.0\" IssueInstant=\"2022-04-13T09:45:38.162Z\" Destination=\"http://localhost:50002/saml/SSO\" AssertionConsumerServiceURL=\"http://localhost:8000/saml/acs\" ProtocolBinding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"><saml:Issuer Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:entity\">http://localhost:8000/saml/metadata</saml:Issuer><ds:Signature xmlns:ds=\"http://www.w3.org/2000/09/xmldsig#\"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/><ds:SignatureMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#rsa-sha1\"/><ds:Reference URI=\"#id-6594153690cfdc336afa98e39145e9716bb88e0f\"><ds:Transforms><ds:Transform Algorithm=\"http://www.w3.org/2000/09/xmldsig#enveloped-signature\"/><ds:Transform Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/></ds:Transforms><ds:DigestMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#sha1\"/><ds:DigestValue>uJqDdJyOLqWh7H0JPGuH9+/FbcY=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:KeyInfo><ds:X509Data><ds:X509Certificate>MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQwNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVoLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+ula1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4WsfawWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oaLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=</ds:X509Certificate></ds:X509Data></ds:KeyInfo></ds:Signature><samlp:NameIDPolicy Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\" AllowCreate=\"true\"/></samlp:AuthnRequest>",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			err: true,
		},
		{
			name: "crewjam wrong digestValue",
			args: args{
				authRequest: "<samlp:AuthnRequest xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" ID=\"id-6594153690cfdc336afa98e39145e9716bb88e0f\" Version=\"2.0\" IssueInstant=\"2022-04-13T09:45:38.162Z\" Destination=\"http://localhost:50002/saml/SSO\" AssertionConsumerServiceURL=\"http://localhost:8000/saml/acs\" ProtocolBinding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"><saml:Issuer Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:entity\">http://localhost:8000/saml/metadata</saml:Issuer><ds:Signature xmlns:ds=\"http://www.w3.org/2000/09/xmldsig#\"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/><ds:SignatureMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#rsa-sha1\"/><ds:Reference URI=\"#id-6594153690cfdc336afa98e39145e9716bb88e0f\"><ds:Transforms><ds:Transform Algorithm=\"http://www.w3.org/2000/09/xmldsig#enveloped-signature\"/><ds:Transform Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/></ds:Transforms><ds:DigestMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#sha1\"/><ds:DigestValue>wrongdigest</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>XQLDINOWHYaGtt53WetQ7G3twWjkMAR+BKoO91LQWTYsv7NxxvVN9F1PyMxEJQyBnfMWnbckokIF0F6jieNNByWm6yX7EG9WOvc9ymb6e7Cx+zElWdvBKbO0Is3cQWiDTn3teexZjcVD1fAdHgpiTA1vYPYaR40pVwVbTQ+qSf7SuYHQifoDzGFIrPGIDZ9PkQ995fQbo8rdLwmJxAVjkMd/24z3iC3h5h6WqVT/9xcGSr8Rr8SQ9rKLZBgpwz2me0PcT97NHHhHR8K6Am0cXTFaKhfZxQQB601RgO6/38BoOtXdIXh+veBIHBlbNagEkJfnzvBNQ6KAYQ58EnVC3w==</ds:SignatureValue><ds:KeyInfo><ds:X509Data><ds:X509Certificate>MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQwNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVoLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+ula1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4WsfawWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oaLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=</ds:X509Certificate></ds:X509Data></ds:KeyInfo></ds:Signature><samlp:NameIDPolicy Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\" AllowCreate=\"true\"/></samlp:AuthnRequest>",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			err: true,
		},
		{
			name: "crewjam wrong signatureValue",
			args: args{
				authRequest: "<samlp:AuthnRequest xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" ID=\"id-6594153690cfdc336afa98e39145e9716bb88e0f\" Version=\"2.0\" IssueInstant=\"2022-04-13T09:45:38.162Z\" Destination=\"http://localhost:50002/saml/SSO\" AssertionConsumerServiceURL=\"http://localhost:8000/saml/acs\" ProtocolBinding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"><saml:Issuer Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:entity\">http://localhost:8000/saml/metadata</saml:Issuer><ds:Signature xmlns:ds=\"http://www.w3.org/2000/09/xmldsig#\"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/><ds:SignatureMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#rsa-sha1\"/><ds:Reference URI=\"#id-6594153690cfdc336afa98e39145e9716bb88e0f\"><ds:Transforms><ds:Transform Algorithm=\"http://www.w3.org/2000/09/xmldsig#enveloped-signature\"/><ds:Transform Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/></ds:Transforms><ds:DigestMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#sha1\"/><ds:DigestValue>uJqDdJyOLqWh7H0JPGuH9+/FbcY=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>wrongsignature</ds:SignatureValue><ds:KeyInfo><ds:X509Data><ds:X509Certificate>MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQwNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVoLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+ula1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4WsfawWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oaLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=</ds:X509Certificate></ds:X509Data></ds:KeyInfo></ds:Signature><samlp:NameIDPolicy Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\" AllowCreate=\"true\"/></samlp:AuthnRequest>",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			err: true,
		},
		{
			name: "crewjam wrong digestAlgorithm",
			args: args{
				authRequest: "<samlp:AuthnRequest xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" ID=\"id-6594153690cfdc336afa98e39145e9716bb88e0f\" Version=\"2.0\" IssueInstant=\"2022-04-13T09:45:38.162Z\" Destination=\"http://localhost:50002/saml/SSO\" AssertionConsumerServiceURL=\"http://localhost:8000/saml/acs\" ProtocolBinding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"><saml:Issuer Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:entity\">http://localhost:8000/saml/metadata</saml:Issuer><ds:Signature xmlns:ds=\"http://www.w3.org/2000/09/xmldsig#\"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/><ds:SignatureMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#rsa-sha1\"/><ds:Reference URI=\"#id-6594153690cfdc336afa98e39145e9716bb88e0f\"><ds:Transforms><ds:Transform Algorithm=\"http://www.w3.org/2000/09/xmldsig#enveloped-signature\"/><ds:Transform Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/></ds:Transforms><ds:DigestMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#sha256\"/><ds:DigestValue>uJqDdJyOLqWh7H0JPGuH9+/FbcY=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>XQLDINOWHYaGtt53WetQ7G3twWjkMAR+BKoO91LQWTYsv7NxxvVN9F1PyMxEJQyBnfMWnbckokIF0F6jieNNByWm6yX7EG9WOvc9ymb6e7Cx+zElWdvBKbO0Is3cQWiDTn3teexZjcVD1fAdHgpiTA1vYPYaR40pVwVbTQ+qSf7SuYHQifoDzGFIrPGIDZ9PkQ995fQbo8rdLwmJxAVjkMd/24z3iC3h5h6WqVT/9xcGSr8Rr8SQ9rKLZBgpwz2me0PcT97NHHhHR8K6Am0cXTFaKhfZxQQB601RgO6/38BoOtXdIXh+veBIHBlbNagEkJfnzvBNQ6KAYQ58EnVC3w==</ds:SignatureValue><ds:KeyInfo><ds:X509Data><ds:X509Certificate>MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQwNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVoLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+ula1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4WsfawWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oaLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=</ds:X509Certificate></ds:X509Data></ds:KeyInfo></ds:Signature><samlp:NameIDPolicy Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\" AllowCreate=\"true\"/></samlp:AuthnRequest>",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			err: true,
		},
		{
			name: "crewjam wrong signatureAlgorithm",
			args: args{
				authRequest: "<samlp:AuthnRequest xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" ID=\"id-6594153690cfdc336afa98e39145e9716bb88e0f\" Version=\"2.0\" IssueInstant=\"2022-04-13T09:45:38.162Z\" Destination=\"http://localhost:50002/saml/SSO\" AssertionConsumerServiceURL=\"http://localhost:8000/saml/acs\" ProtocolBinding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"><saml:Issuer Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:entity\">http://localhost:8000/saml/metadata</saml:Issuer><ds:Signature xmlns:ds=\"http://www.w3.org/2000/09/xmldsig#\"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/><ds:SignatureMethod Algorithm=\"http://www.w3.org/2001/04/xmldsig-more#rsa-sha256\"/><ds:Reference URI=\"#id-6594153690cfdc336afa98e39145e9716bb88e0f\"><ds:Transforms><ds:Transform Algorithm=\"http://www.w3.org/2000/09/xmldsig#enveloped-signature\"/><ds:Transform Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/></ds:Transforms><ds:DigestMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#sha1\"/><ds:DigestValue>uJqDdJyOLqWh7H0JPGuH9+/FbcY=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>XQLDINOWHYaGtt53WetQ7G3twWjkMAR+BKoO91LQWTYsv7NxxvVN9F1PyMxEJQyBnfMWnbckokIF0F6jieNNByWm6yX7EG9WOvc9ymb6e7Cx+zElWdvBKbO0Is3cQWiDTn3teexZjcVD1fAdHgpiTA1vYPYaR40pVwVbTQ+qSf7SuYHQifoDzGFIrPGIDZ9PkQ995fQbo8rdLwmJxAVjkMd/24z3iC3h5h6WqVT/9xcGSr8Rr8SQ9rKLZBgpwz2me0PcT97NHHhHR8K6Am0cXTFaKhfZxQQB601RgO6/38BoOtXdIXh+veBIHBlbNagEkJfnzvBNQ6KAYQ58EnVC3w==</ds:SignatureValue><ds:KeyInfo><ds:X509Data><ds:X509Certificate>MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQwNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVoLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+ula1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4WsfawWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oaLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=</ds:X509Certificate></ds:X509Data></ds:KeyInfo></ds:Signature><samlp:NameIDPolicy Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\" AllowCreate=\"true\"/></samlp:AuthnRequest>",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			err: true,
		},
		{
			name: "crewjam unknown digestAlgorithm",
			args: args{
				authRequest: "<samlp:AuthnRequest xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" ID=\"id-6594153690cfdc336afa98e39145e9716bb88e0f\" Version=\"2.0\" IssueInstant=\"2022-04-13T09:45:38.162Z\" Destination=\"http://localhost:50002/saml/SSO\" AssertionConsumerServiceURL=\"http://localhost:8000/saml/acs\" ProtocolBinding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"><saml:Issuer Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:entity\">http://localhost:8000/saml/metadata</saml:Issuer><ds:Signature xmlns:ds=\"http://www.w3.org/2000/09/xmldsig#\"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/><ds:SignatureMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#rsa-sha1\"/><ds:Reference URI=\"#id-6594153690cfdc336afa98e39145e9716bb88e0f\"><ds:Transforms><ds:Transform Algorithm=\"http://www.w3.org/2000/09/xmldsig#enveloped-signature\"/><ds:Transform Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/></ds:Transforms><ds:DigestMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#sha300\"/><ds:DigestValue>uJqDdJyOLqWh7H0JPGuH9+/FbcY=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>XQLDINOWHYaGtt53WetQ7G3twWjkMAR+BKoO91LQWTYsv7NxxvVN9F1PyMxEJQyBnfMWnbckokIF0F6jieNNByWm6yX7EG9WOvc9ymb6e7Cx+zElWdvBKbO0Is3cQWiDTn3teexZjcVD1fAdHgpiTA1vYPYaR40pVwVbTQ+qSf7SuYHQifoDzGFIrPGIDZ9PkQ995fQbo8rdLwmJxAVjkMd/24z3iC3h5h6WqVT/9xcGSr8Rr8SQ9rKLZBgpwz2me0PcT97NHHhHR8K6Am0cXTFaKhfZxQQB601RgO6/38BoOtXdIXh+veBIHBlbNagEkJfnzvBNQ6KAYQ58EnVC3w==</ds:SignatureValue><ds:KeyInfo><ds:X509Data><ds:X509Certificate>MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQwNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVoLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+ula1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4WsfawWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oaLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=</ds:X509Certificate></ds:X509Data></ds:KeyInfo></ds:Signature><samlp:NameIDPolicy Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\" AllowCreate=\"true\"/></samlp:AuthnRequest>",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			err: true,
		},
		{
			name: "crewjam unknown signatureAlgorithm",
			args: args{
				authRequest: "<samlp:AuthnRequest xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" ID=\"id-6594153690cfdc336afa98e39145e9716bb88e0f\" Version=\"2.0\" IssueInstant=\"2022-04-13T09:45:38.162Z\" Destination=\"http://localhost:50002/saml/SSO\" AssertionConsumerServiceURL=\"http://localhost:8000/saml/acs\" ProtocolBinding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"><saml:Issuer Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:entity\">http://localhost:8000/saml/metadata</saml:Issuer><ds:Signature xmlns:ds=\"http://www.w3.org/2000/09/xmldsig#\"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/><ds:SignatureMethod Algorithm=\"http://www.w3.org/2001/04/xmldsig-more#rsa-sha300\"/><ds:Reference URI=\"#id-6594153690cfdc336afa98e39145e9716bb88e0f\"><ds:Transforms><ds:Transform Algorithm=\"http://www.w3.org/2000/09/xmldsig#enveloped-signature\"/><ds:Transform Algorithm=\"http://www.w3.org/2001/10/xml-exc-c14n#\"/></ds:Transforms><ds:DigestMethod Algorithm=\"http://www.w3.org/2000/09/xmldsig#sha1\"/><ds:DigestValue>uJqDdJyOLqWh7H0JPGuH9+/FbcY=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>XQLDINOWHYaGtt53WetQ7G3twWjkMAR+BKoO91LQWTYsv7NxxvVN9F1PyMxEJQyBnfMWnbckokIF0F6jieNNByWm6yX7EG9WOvc9ymb6e7Cx+zElWdvBKbO0Is3cQWiDTn3teexZjcVD1fAdHgpiTA1vYPYaR40pVwVbTQ+qSf7SuYHQifoDzGFIrPGIDZ9PkQ995fQbo8rdLwmJxAVjkMd/24z3iC3h5h6WqVT/9xcGSr8Rr8SQ9rKLZBgpwz2me0PcT97NHHhHR8K6Am0cXTFaKhfZxQQB601RgO6/38BoOtXdIXh+veBIHBlbNagEkJfnzvBNQ6KAYQ58EnVC3w==</ds:SignatureValue><ds:KeyInfo><ds:X509Data><ds:X509Certificate>MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQwNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVoLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+ula1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4WsfawWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oaLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=</ds:X509Certificate></ds:X509Data></ds:KeyInfo></ds:Signature><samlp:NameIDPolicy Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\" AllowCreate=\"true\"/></samlp:AuthnRequest>",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := etree.NewDocument()
			if err := doc.ReadFromBytes([]byte(tt.args.authRequest)); err != nil {
				t.Errorf("failed to read request")
				return
			}

			certStrs := []string{strings.ReplaceAll(tt.args.certificate, "\n", "")}
			certs, err := signature.ParseCertificates(certStrs)
			if err != nil {
				t.Errorf("failed to read request")
				return
			}

			if err := signature.ValidatePost(certs, doc.Root()); (err != nil) != tt.err {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.err)
				return
			}
		})
	}
}

func TestSignature_ValidateRedirect(t *testing.T) {
	type args struct {
		sigAlg      string
		element     string
		sig         string
		certificate string
	}
	tests := []struct {
		name string
		args args
		err  bool
	}{
		{
			"crewjam 1",
			args{
				sigAlg:      dsig.RSASHA1SignatureMethod,
				element:     "SAMLRequest=nJJBj9MwEIX%2FijX3NI6blsbaRCpbISotbLUpHLhNkwm15NjFMwH236O2i7RIqIe92vP5ved5d4yjP9n1JMfwRD8mYlG%2FRx%2FYni9qmFKwEdmxDTgSW%2Blsu%2F70YM1MW2SmJC4GeIWcbjOnFCV20YPabmpwfXZYlVgdBq1xoHmvl4hVvyBcvSv7eVkO1Wquq6EgA%2BorJXYx1GBmGtSWeaJtYMEgNRhtTKbLzOh9sbTFwhbVbFksvoHaEIsLKBfyKHKyee5jh%2F4YWexCa23ys%2B%2B8bR9Brf9Guo%2BBp5FSS%2Bmn6%2BjL08N%2F6JXW%2Bgpjx6B2L9neu9C78P32RxyuQ2w%2F7ve7bPfY7qG5LMNekiX1IaYR5fYj5xPXZ8Nl1FIQJ8%2FQ3PA5kmCPgnf5K6nmpQSfcaTtZhe9657fIC8JAzsKAmrtffx1nwiFapA0EeTNVfLfqjV%2FAgAA%2F%2F8%3D&RelayState=E_4Asq22pZf6vLuXaNStyUJBb3sBsqemdaPNZGoFPXSNu0iSsoAYRVPH&SigAlg=http%3A%2F%2Fwww.w3.org%2F2000%2F09%2Fxmldsig%23rsa-sha1",
				sig:         "bKLb42HR9q+hGKpnHTTYAufQSd5uVDnTAsDy6qT2VhqQskdGJ6vlqnp7WWeK1ANBE96Rh+YKUhO5hbJzca45+uECITQ1NzEqFxKtYB55lxUJY6O4XJPDNbtNsBwtGHHH/alF9DrD/LecK9l5Yo8ewdLzWR9KxH2XzU7uZy8QOhLhNMvRSgsOIABkivE3cv+Ik1IqqjxSkUtoTyBBRbRT+jSgqczF9B2Gn/YO8my5KukovTYkqqJ/HJulST5gzVuhx42FDsjVXE1uhHva4VDndzCRCw21CgPbR+DC/AxokNKnVGodXFg/xEp3M0JKGoZWzXUXnI5EHd/r1mSaizeJCg==",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			false,
		},
		{
			"crewjam 2",
			args{
				sigAlg:      dsig.RSASHA1SignatureMethod,
				element:     "SAMLRequest=nJJBj9owEIX%2FijX3kMEkLFibSHRRVaRtF21oD70NyVAsJTb1TNruv6%2BArbSVKg692vP5ved590JDf3KrUY%2Fhmb%2BPLGp%2BDX0Qd76oYEzBRRIvLtDA4rR1zerjo7MTdCTCSX0M8AY53WZOKWpsYw9ms67Ad9ns7m7f0X66KBBxVnazAy%2BXRYmHtixKtMv5dMbzZbEA84WT%2BBgqsBMEsxEZeRNEKWgFFq3NsMjsdIcLVxbOFhOcl1%2FBrFnUB9ILeVQ9uTzvY0v9MYq6EhFtfvadN80TmNWfSA8xyDhwajj98C1%2Ffn78B71AxCtMrYDZvmZ750Pnw7fbH7G%2FDon7sNtts%2B1Ts4P6sgx3SZbM%2B5gG0tuPnE98lx0uo46Den2B%2BobPgZU6UrrP30jVryX4RANv1tvY%2B%2FblP%2BQ1URDPQcGs%2Bj7%2BfEhMyhVoGhny%2Bir5d9Xq3wEAAP%2F%2F&RelayState=BXpR_zsuXNUyartSlxS5yDV_5JXGdLn-7EJJEaqVfhPySY_nYrbkCobT&SigAlg=http%3A%2F%2Fwww.w3.org%2F2000%2F09%2Fxmldsig%23rsa-sha1",
				sig:         "JqMteNAaoSWQcn+EbuTFrSt1E6B2rMmOuQioioidxywQoVGNipzUc9RRgh/QZ6ICtV0NHwbDSLd3P61WEttteBHFCQTHGVgb/BWYFI1Ck2Dp+svTdGgVZ9ymLOlx6chJoDoR57TxoIjEBQGHHAxTE+UU6g/tI8hVYgm+myMH8EBmxhW1dhsZOi9WeSpptPxfoiv39vm3ZoEe/ya8duoX6A5f8I5gYRifbVlLakBCeKkGL9Iz2jXDMlR/23KnpD7wweY2PDgEQNbtQPXIqxrwYfneNJ9W4s7vTlZl7vZpVgE6Cxx3AgD+QryW6AvsRVKc6TtOHlFGp+Np+SACanduWQ==",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			false,
		},
		{
			"crewjam 3",
			args{
				sigAlg:      dsig.RSASHA1SignatureMethod,
				element:     "SAMLRequest=nJJBj9MwEIX%2FijX3NI6bpVtrE6lshai0sNWmcOA2SSbUkmMXzwTYf4%2FaLtIioRz2as%2Fn957n3TGO%2FmQ3kxzDE%2F2YiEX9Hn1ge76oYErBRmTHNuBIbKWzzebTgzULbZGZkrgY4BVymmdOKUrsoge121bg%2BmypjSnW5ZrawRTlUKx1O7RDOSzJEK3Klvr%2BlnBVgPpKiV0MFZiFBrVjnmgXWDBIBUYbk%2BkyM8VBr62%2BscVqod8tv4HaEosLKBfyKHKyee5jh%2F4YWeyN1trkZ9950zyC2vyNdB8DTyOlhtJP19GXp4f%2F0Lda6yuMHYPav2R770Lvwvf5j2ivQ2w%2FHg77bP%2FYHKC%2BLMNekiX1IaYRZf6R84nrs%2BEyaimIk2eoZ3yOJNij4F3%2BSqp%2BKcFnHGm33Ufvuuc3yEvCwI6CgNp4H3%2FdJ0KhCiRNBHl9lfy3avWfAAAA%2F%2F8%3D&RelayState=Dz6_QBnMSgS0h6MRcC7Pz15HXRPGnT1XOMik7EydB8Fa4U7TwH50HxKa&SigAlg=http%3A%2F%2Fwww.w3.org%2F2000%2F09%2Fxmldsig%23rsa-sha1",
				sig:         "CgY3gl23C1nSAgakKE7c+Dha4oPdB2/XhV9H5TNARFzhgd0LafG7wxEWRgSqPPupuHSO1DH2eTuVjsXi5HiPc/DpGC/t3s5JOI1WWukr7OT1Q9GNv3trAuACtrVOLoTk8H0w197sPb6fLsBrciL5QTVu0e0NgztusHnH2ZvJl6XO0CHvGLoKE8+l9fygwpLrbSPrU33HN0WEiSaoKflqeGagPMS0gAFkeFkWOFtqlXN69xACMa76ciDKGU68GYSRvm8RqU8jhpRotLypimEsTmebcwramDCVWHnVd9DS7MTnW019PuucaI3BRMqFsQtje+DdcH735vZvj9AI8qGRQA==",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			false,
		},
		{
			"crewjam wrong relay state",
			args{
				sigAlg:      dsig.RSASHA1SignatureMethod,
				element:     "SAMLRequest=nJJBj9MwEIX%2FijX3NI6blsbaRCpbISotbLUpHLhNkwm15NjFMwH236O2i7RIqIe92vP5ved5d4yjP9n1JMfwRD8mYlG%2FRx%2FYni9qmFKwEdmxDTgSW%2Blsu%2F70YM1MW2SmJC4GeIWcbjOnFCV20YPabmpwfXZYlVgdBq1xoHmvl4hVvyBcvSv7eVkO1Wquq6EgA%2BorJXYx1GBmGtSWeaJtYMEgNRhtTKbLzOh9sbTFwhbVbFksvoHaEIsLKBfyKHKyee5jh%2F4YWexCa23ys%2B%2B8bR9Brf9Guo%2BBp5FSS%2Bmn6%2BjL08N%2F6JXW%2Bgpjx6B2L9neu9C78P32RxyuQ2w%2F7ve7bPfY7qG5LMNekiX1IaYR5fYj5xPXZ8Nl1FIQJ8%2FQ3PA5kmCPgnf5K6nmpQSfcaTtZhe9657fIC8JAzsKAmrtffx1nwiFapA0EeTNVfLfqjV%2FAgAA%2F%2F8%3D&RelayState=test&SigAlg=http%3A%2F%2Fwww.w3.org%2F2000%2F09%2Fxmldsig%23rsa-sha1",
				sig:         "bKLb42HR9q+hGKpnHTTYAufQSd5uVDnTAsDy6qT2VhqQskdGJ6vlqnp7WWeK1ANBE96Rh+YKUhO5hbJzca45+uECITQ1NzEqFxKtYB55lxUJY6O4XJPDNbtNsBwtGHHH/alF9DrD/LecK9l5Yo8ewdLzWR9KxH2XzU7uZy8QOhLhNMvRSgsOIABkivE3cv+Ik1IqqjxSkUtoTyBBRbRT+jSgqczF9B2Gn/YO8my5KukovTYkqqJ/HJulST5gzVuhx42FDsjVXE1uhHva4VDndzCRCw21CgPbR+DC/AxokNKnVGodXFg/xEp3M0JKGoZWzXUXnI5EHd/r1mSaizeJCg==",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			true,
		},
		{
			"crewjam wrong request",
			args{
				sigAlg:      dsig.RSASHA1SignatureMethod,
				element:     "SAMLRequest=test&RelayState=E_4Asq22pZf6vLuXaNStyUJBb3sBsqemdaPNZGoFPXSNu0iSsoAYRVPH&SigAlg=http%3A%2F%2Fwww.w3.org%2F2000%2F09%2Fxmldsig%23rsa-sha1",
				sig:         "bKLb42HR9q+hGKpnHTTYAufQSd5uVDnTAsDy6qT2VhqQskdGJ6vlqnp7WWeK1ANBE96Rh+YKUhO5hbJzca45+uECITQ1NzEqFxKtYB55lxUJY6O4XJPDNbtNsBwtGHHH/alF9DrD/LecK9l5Yo8ewdLzWR9KxH2XzU7uZy8QOhLhNMvRSgsOIABkivE3cv+Ik1IqqjxSkUtoTyBBRbRT+jSgqczF9B2Gn/YO8my5KukovTYkqqJ/HJulST5gzVuhx42FDsjVXE1uhHva4VDndzCRCw21CgPbR+DC/AxokNKnVGodXFg/xEp3M0JKGoZWzXUXnI5EHd/r1mSaizeJCg==",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			true,
		},
		{
			"crewjam wrong sigalg",
			args{
				sigAlg:      dsig.RSASHA1SignatureMethod,
				element:     "SAMLRequest=nJJBj9MwEIX%2FijX3NI6blsbaRCpbISotbLUpHLhNkwm15NjFMwH236O2i7RIqIe92vP5ved5d4yjP9n1JMfwRD8mYlG%2FRx%2FYni9qmFKwEdmxDTgSW%2Blsu%2F70YM1MW2SmJC4GeIWcbjOnFCV20YPabmpwfXZYlVgdBq1xoHmvl4hVvyBcvSv7eVkO1Wquq6EgA%2BorJXYx1GBmGtSWeaJtYMEgNRhtTKbLzOh9sbTFwhbVbFksvoHaEIsLKBfyKHKyee5jh%2F4YWexCa23ys%2B%2B8bR9Brf9Guo%2BBp5FSS%2Bmn6%2BjL08N%2F6JXW%2Bgpjx6B2L9neu9C78P32RxyuQ2w%2F7ve7bPfY7qG5LMNekiX1IaYR5fYj5xPXZ8Nl1FIQJ8%2FQ3PA5kmCPgnf5K6nmpQSfcaTtZhe9657fIC8JAzsKAmrtffx1nwiFapA0EeTNVfLfqjV%2FAgAA%2F%2F8%3D&RelayState=E_4Asq22pZf6vLuXaNStyUJBb3sBsqemdaPNZGoFPXSNu0iSsoAYRVPH&SigAlg=http%3A%2F%2Fwww.w3.org%2F2001%2F04%2Fxmldsig-more%23rsa-sha256",
				sig:         "bKLb42HR9q+hGKpnHTTYAufQSd5uVDnTAsDy6qT2VhqQskdGJ6vlqnp7WWeK1ANBE96Rh+YKUhO5hbJzca45+uECITQ1NzEqFxKtYB55lxUJY6O4XJPDNbtNsBwtGHHH/alF9DrD/LecK9l5Yo8ewdLzWR9KxH2XzU7uZy8QOhLhNMvRSgsOIABkivE3cv+Ik1IqqjxSkUtoTyBBRbRT+jSgqczF9B2Gn/YO8my5KukovTYkqqJ/HJulST5gzVuhx42FDsjVXE1uhHva4VDndzCRCw21CgPbR+DC/AxokNKnVGodXFg/xEp3M0JKGoZWzXUXnI5EHd/r1mSaizeJCg==",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			true,
		},
		{
			"crewjam wrong sig",
			args{
				sigAlg:      dsig.RSASHA1SignatureMethod,
				element:     "SAMLRequest=nJJBj9MwEIX%2FijX3NI6blsbaRCpbISotbLUpHLhNkwm15NjFMwH236O2i7RIqIe92vP5ved5d4yjP9n1JMfwRD8mYlG%2FRx%2FYni9qmFKwEdmxDTgSW%2Blsu%2F70YM1MW2SmJC4GeIWcbjOnFCV20YPabmpwfXZYlVgdBq1xoHmvl4hVvyBcvSv7eVkO1Wquq6EgA%2BorJXYx1GBmGtSWeaJtYMEgNRhtTKbLzOh9sbTFwhbVbFksvoHaEIsLKBfyKHKyee5jh%2F4YWexCa23ys%2B%2B8bR9Brf9Guo%2BBp5FSS%2Bmn6%2BjL08N%2F6JXW%2Bgpjx6B2L9neu9C78P32RxyuQ2w%2F7ve7bPfY7qG5LMNekiX1IaYR5fYj5xPXZ8Nl1FIQJ8%2FQ3PA5kmCPgnf5K6nmpQSfcaTtZhe9657fIC8JAzsKAmrtffx1nwiFapA0EeTNVfLfqjV%2FAgAA%2F%2F8%3D&RelayState=E_4Asq22pZf6vLuXaNStyUJBb3sBsqemdaPNZGoFPXSNu0iSsoAYRVPH&SigAlg=http%3A%2F%2Fwww.w3.org%2F2000%2F09%2Fxmldsig%23rsa-sha1",
				sig:         "test",
				certificate: "MIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=",
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regex := regexp.MustCompile(`\s+`)
			certStr := regex.ReplaceAllString(tt.args.certificate, "")
			certStr = strings.ReplaceAll(certStr, "\n", "")

			certBytes, err := base64.StdEncoding.DecodeString(certStr)
			if err != nil {
				t.Errorf("failed to parse PEM block containing the public key")
				return
			}
			parsedCert, err := x509.ParseCertificate(certBytes)
			if err != nil {
				t.Errorf("failed to parse certificate: " + err.Error())
				return
			}

			signatureValue, err := base64.StdEncoding.DecodeString(tt.args.sig)
			if err != nil {
				t.Errorf("failed to decode sig: " + err.Error())
				return
			}

			if err := signature.ValidateRedirect(tt.args.sigAlg, []byte(tt.args.element), signatureValue, parsedCert.PublicKey); (err != nil) != tt.err {
				t.Errorf("ValidateRedirect() error = %v, wantErr %v", err, tt.err)
				return
			}
		})
	}
}

func TestSignature_CreateRedirect(t *testing.T) {
	type args struct {
		certificate     string
		key             string
		sigAlg          string
		query           string
		signingContextF func()
	}
	tests := []struct {
		name string
		args args
		err  bool
	}{
		{
			"successful 1",
			args{
				"-----BEGIN CERTIFICATE-----\nMIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=\n-----END CERTIFICATE-----\n",
				"-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC7XKdCRxUZXjdq\nVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWc\nWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2\ngIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+a\nEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7\nLfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v\n4cxTNPn/AgMBAAECggEAF+rV9yH30Ysza8GwrXCR9qDN1Dp3QmmsavnXkonEvPoq\nEr2T3o0//6mBp6CLDboMQGQBjblJwl+3Y6PgZolvHAMOsMdHfYNPEo7FSzUBzEw+\nqRrs5HkMyvoPgfV6X8F97W3tiD4Q/AmHkMILl+MxbnfPXM54gWqPuwIqxY1uaCk5\nREwyb7WBon3rd58ceOI1SLRjod6SbqWBMMSN3cJ+5VEPObFjw/RlhNQ5rBI8G5Kt\nso2zBU5C4BB2CvqlWy98WDKJkTvWHbiTjZCy8BQ+gQ6UJM2vaNELFOVpuMGQnMIi\noWiX10Jg2e1gP9j3TdrohlGF8M3+TXjSFKNmeX0DUQKBgQDx7UazUWS5RtkgnjH9\nw2xH2xkstJVD7nAS8VTxNwcrgjVXPvTJha9El904obUjyRX7ppb02tuH5ML/bZh6\n9lL4bP5+SHcJ10e4q8CK/KAGHD6BYAbaGXRq0CoSk5a3vv5XPdob4T5qKCIHFpnu\nMfbvdbEoameLOyRYOGu/yVZIiwKBgQDGQs7FRTisHV0xooiRmlvYF0dcd19qpLed\nqhgJNqBPOTEvvGvJNRoi39haEY3cuTqsxZ5FAlFlVFMUUozz+d0xBLLInoVY/Y4h\nhSdGmdw/A6oHodLqyEp3N5RZNdLlh8/nDS3xXzMotAl75bW5kc2ttcRhRdtyNJ9Z\nup0PgppO3QKBgEC45upAQz8iCiKkz+EA8C4FGqYQJcLHvmoC8GOcAioMqrKNoDVt\ns2cZbdChynEpcd0iQ058YrDnbZeiPWHgFnBp0Gf+gQI7+u8X2+oTDci0s7Au/YZJ\nuxB8YlUX8QF1clvqqzg8OVNzKy9UR5gm+9YyWVPjq5HfH6kOZx0nAxNjAoGAERt8\nqgsCC9/wxbKnpCC0oh3IG5N1WUdjTKh7sHfVN2DQ/LR+fHsniTDVg1gWbKBTDsty\nj7PWgC7ZiFxjKz45NtyX7LW4/efLFttdezsVhR500nnFMFseCdFy7Iu3afThHKfH\nehdj27RFSTqWBrAtFjsj+dzERcOCqIRwvwDe/cUCgYEA5+1mzVXDVjKsWylKJPk+\nZZA4LUfvmTj3VLNDZrlSAI/xEikCFio0QWEA2TQYTAwbXTrKwQSeHQRhv7OTc1h+\nMhpAgvs189ze5J4jiNmULEkkrO+Cxxnw8tyV+UFRZtzW9gUoVBwXiZ/Wbl9sfnlO\nwLJHc0j6OltPcPJmxHP8gQI=\n-----END PRIVATE KEY-----\n",
				dsig.RSASHA1SignatureMethod,
				"test1",
				nil,
			},
			false,
		},
		{
			"successful 2",
			args{
				"-----BEGIN CERTIFICATE-----\nMIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=\n-----END CERTIFICATE-----\n",
				"-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC7XKdCRxUZXjdq\nVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWc\nWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2\ngIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+a\nEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7\nLfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v\n4cxTNPn/AgMBAAECggEAF+rV9yH30Ysza8GwrXCR9qDN1Dp3QmmsavnXkonEvPoq\nEr2T3o0//6mBp6CLDboMQGQBjblJwl+3Y6PgZolvHAMOsMdHfYNPEo7FSzUBzEw+\nqRrs5HkMyvoPgfV6X8F97W3tiD4Q/AmHkMILl+MxbnfPXM54gWqPuwIqxY1uaCk5\nREwyb7WBon3rd58ceOI1SLRjod6SbqWBMMSN3cJ+5VEPObFjw/RlhNQ5rBI8G5Kt\nso2zBU5C4BB2CvqlWy98WDKJkTvWHbiTjZCy8BQ+gQ6UJM2vaNELFOVpuMGQnMIi\noWiX10Jg2e1gP9j3TdrohlGF8M3+TXjSFKNmeX0DUQKBgQDx7UazUWS5RtkgnjH9\nw2xH2xkstJVD7nAS8VTxNwcrgjVXPvTJha9El904obUjyRX7ppb02tuH5ML/bZh6\n9lL4bP5+SHcJ10e4q8CK/KAGHD6BYAbaGXRq0CoSk5a3vv5XPdob4T5qKCIHFpnu\nMfbvdbEoameLOyRYOGu/yVZIiwKBgQDGQs7FRTisHV0xooiRmlvYF0dcd19qpLed\nqhgJNqBPOTEvvGvJNRoi39haEY3cuTqsxZ5FAlFlVFMUUozz+d0xBLLInoVY/Y4h\nhSdGmdw/A6oHodLqyEp3N5RZNdLlh8/nDS3xXzMotAl75bW5kc2ttcRhRdtyNJ9Z\nup0PgppO3QKBgEC45upAQz8iCiKkz+EA8C4FGqYQJcLHvmoC8GOcAioMqrKNoDVt\ns2cZbdChynEpcd0iQ058YrDnbZeiPWHgFnBp0Gf+gQI7+u8X2+oTDci0s7Au/YZJ\nuxB8YlUX8QF1clvqqzg8OVNzKy9UR5gm+9YyWVPjq5HfH6kOZx0nAxNjAoGAERt8\nqgsCC9/wxbKnpCC0oh3IG5N1WUdjTKh7sHfVN2DQ/LR+fHsniTDVg1gWbKBTDsty\nj7PWgC7ZiFxjKz45NtyX7LW4/efLFttdezsVhR500nnFMFseCdFy7Iu3afThHKfH\nehdj27RFSTqWBrAtFjsj+dzERcOCqIRwvwDe/cUCgYEA5+1mzVXDVjKsWylKJPk+\nZZA4LUfvmTj3VLNDZrlSAI/xEikCFio0QWEA2TQYTAwbXTrKwQSeHQRhv7OTc1h+\nMhpAgvs189ze5J4jiNmULEkkrO+Cxxnw8tyV+UFRZtzW9gUoVBwXiZ/Wbl9sfnlO\nwLJHc0j6OltPcPJmxHP8gQI=\n-----END PRIVATE KEY-----\n",
				dsig.RSASHA1SignatureMethod,
				"test2",
				nil,
			},
			false,
		},
		{
			"successful 3",
			args{
				"-----BEGIN CERTIFICATE-----\nMIICvDCCAaQCCQD6E8ZGsQ2usjANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVt\neXNlcnZpY2UuZXhhbXBsZS5jb20wHhcNMjIwMjE3MTQwNjM5WhcNMjMwMjE3MTQw\nNjM5WjAgMR4wHAYDVQQDDBVteXNlcnZpY2UuZXhhbXBsZS5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7XKdCRxUZXjdqVqwwwOJqc1Ch0nOSmk+U\nerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWcWAHJloqZ7GBS7NpDhzV8\nG+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2gIfsYPs3TTq1sq7oCs5q\nLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+aEkyRh07oMpXBEobGisfF\n2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7Lfgq7oxmv/8LFi4Zopr5\nnyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v4cxTNPn/AgMBAAEwDQYJ\nKoZIhvcNAQELBQADggEBAJYxROWSOZbOzXzafdGjQKsMgN948G/hHwVuZneyAcVo\nLMFTs1Weya9Z+snMp1u0AdDGmQTS9zGnD7syDYGOmgigOLcMvLMoWf5tCQBbEukW\n8O7DPjRR0XypChGSsHsqLGO0B0HaTel0HdP9Si827OCkc9Q+WbsFG/8/4ToGWL+u\nla1WuLawozoj8umPi9D8iXCoW35y2STU+WFQG7W+Kfdu+2CYz/0tGdwVqNG4Wsfa\nwWchrS00vGFKjm/fJc876gAfxiMH1I9fZvYSAxAZ3sVI//Ml2sUdgf067ywQ75oa\nLSS2NImmz5aos3vuWmOXhILd7iTU+BD8Uv6vWbI7I1M=\n-----END CERTIFICATE-----\n",
				"-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC7XKdCRxUZXjdq\nVqwwwOJqc1Ch0nOSmk+UerkUqlviWHdeLR+FolHKjqLzCBloAz4xVc0DFfR76gWc\nWAHJloqZ7GBS7NpDhzV8G+cXQ+bTU0Lu2e73zCQb30XUdKhWiGfDKaU+1xg9CD/2\ngIfsYPs3TTq1sq7oCs5qLdUHaVL5kcRaHKdnTi7cs5i9xzs3TsUnXcrJPwydjp+a\nEkyRh07oMpXBEobGisfF2p1MA6pVW2gjmywf7D5iYEFELQhM7poqPN3/kfBvU1n7\nLfgq7oxmv/8LFi4Zopr5nyqsz26XPtUy1WqTzgznAmP+nN0oBTERFVbXXdRa3k2v\n4cxTNPn/AgMBAAECggEAF+rV9yH30Ysza8GwrXCR9qDN1Dp3QmmsavnXkonEvPoq\nEr2T3o0//6mBp6CLDboMQGQBjblJwl+3Y6PgZolvHAMOsMdHfYNPEo7FSzUBzEw+\nqRrs5HkMyvoPgfV6X8F97W3tiD4Q/AmHkMILl+MxbnfPXM54gWqPuwIqxY1uaCk5\nREwyb7WBon3rd58ceOI1SLRjod6SbqWBMMSN3cJ+5VEPObFjw/RlhNQ5rBI8G5Kt\nso2zBU5C4BB2CvqlWy98WDKJkTvWHbiTjZCy8BQ+gQ6UJM2vaNELFOVpuMGQnMIi\noWiX10Jg2e1gP9j3TdrohlGF8M3+TXjSFKNmeX0DUQKBgQDx7UazUWS5RtkgnjH9\nw2xH2xkstJVD7nAS8VTxNwcrgjVXPvTJha9El904obUjyRX7ppb02tuH5ML/bZh6\n9lL4bP5+SHcJ10e4q8CK/KAGHD6BYAbaGXRq0CoSk5a3vv5XPdob4T5qKCIHFpnu\nMfbvdbEoameLOyRYOGu/yVZIiwKBgQDGQs7FRTisHV0xooiRmlvYF0dcd19qpLed\nqhgJNqBPOTEvvGvJNRoi39haEY3cuTqsxZ5FAlFlVFMUUozz+d0xBLLInoVY/Y4h\nhSdGmdw/A6oHodLqyEp3N5RZNdLlh8/nDS3xXzMotAl75bW5kc2ttcRhRdtyNJ9Z\nup0PgppO3QKBgEC45upAQz8iCiKkz+EA8C4FGqYQJcLHvmoC8GOcAioMqrKNoDVt\ns2cZbdChynEpcd0iQ058YrDnbZeiPWHgFnBp0Gf+gQI7+u8X2+oTDci0s7Au/YZJ\nuxB8YlUX8QF1clvqqzg8OVNzKy9UR5gm+9YyWVPjq5HfH6kOZx0nAxNjAoGAERt8\nqgsCC9/wxbKnpCC0oh3IG5N1WUdjTKh7sHfVN2DQ/LR+fHsniTDVg1gWbKBTDsty\nj7PWgC7ZiFxjKz45NtyX7LW4/efLFttdezsVhR500nnFMFseCdFy7Iu3afThHKfH\nehdj27RFSTqWBrAtFjsj+dzERcOCqIRwvwDe/cUCgYEA5+1mzVXDVjKsWylKJPk+\nZZA4LUfvmTj3VLNDZrlSAI/xEikCFio0QWEA2TQYTAwbXTrKwQSeHQRhv7OTc1h+\nMhpAgvs189ze5J4jiNmULEkkrO+Cxxnw8tyV+UFRZtzW9gUoVBwXiZ/Wbl9sfnlO\nwLJHc0j6OltPcPJmxHP8gQI=\n-----END PRIVATE KEY-----\n",
				dsig.RSASHA1SignatureMethod,
				"test3",
				nil,
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsCert, err := tls.X509KeyPair([]byte(tt.args.certificate), []byte(tt.args.key))
			if err != nil {
				t.Errorf("failed to tls certificate")
				return
			}
			signingContext, err := signature.GetSigningContext(tlsCert, tt.args.sigAlg)
			if err != nil {
				t.Errorf("failed to create signing context")
				return
			}

			sig, err := signature.CreateRedirect(signingContext, tt.args.query)
			if (err != nil) != tt.err {
				t.Errorf("CreateRedirect() error = %v, wantErr %v", err, tt.err)
				return
			}
			priv := tlsCert.PrivateKey.(*rsa.PrivateKey)
			if err := signature.ValidateRedirect(tt.args.sigAlg, []byte(tt.args.query), sig, priv.Public()); (err != nil) != tt.err {
				t.Errorf("ValidateRedirect() error = %v, wantErr %v", err, tt.err)
				return
			}
		})
	}
}

func TestSignature_CreatePost(t *testing.T) {
	type args struct {
		response           string
		singatureAlgorithm string
		canonicalizer      dsig.Canonicalizer
	}
	type res struct {
		canonicalization string
		signatureMethod  string
		digestMethod     string
		certificate      string
		transforms       []string
		err              bool
	}

	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "SAMLTool example response SHA1",
			args: args{
				response:           "<samlp:Response xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" ID=\"_8e8dc5f69a98cc4c1ff3427e5ce34606fd672f91e6\" Version=\"2.0\" IssueInstant=\"2014-07-17T01:01:48Z\" Destination=\"http://sp.example.com/demo1/index.php?acs\" InResponseTo=\"ONELOGIN_4fee3b046395c4e751011e97f8900b5273d56685\">\n  <saml:Issuer>http://idp.example.com/metadata.php</saml:Issuer>\n  <samlp:Status>\n    <samlp:StatusCode Value=\"urn:oasis:names:tc:SAML:2.0:status:Success\"/>\n  </samlp:Status>\n  <saml:Assertion xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xs=\"http://www.w3.org/2001/XMLSchema\" ID=\"_d71a3a8e9fcc45c9e9d248ef7049393fc8f04e5f75\" Version=\"2.0\" IssueInstant=\"2014-07-17T01:01:48Z\">\n    <saml:Issuer>http://idp.example.com/metadata.php</saml:Issuer>\n    <saml:Subject>\n      <saml:NameID SPNameQualifier=\"http://sp.example.com/demo1/metadata.php\" Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\">_ce3d2948b4cf20146dee0a0b3dd6f69b6cf86f62d7</saml:NameID>\n      <saml:SubjectConfirmation Method=\"urn:oasis:names:tc:SAML:2.0:cm:bearer\">\n        <saml:SubjectConfirmationData NotOnOrAfter=\"2024-01-18T06:21:48Z\" Recipient=\"http://sp.example.com/demo1/index.php?acs\" InResponseTo=\"ONELOGIN_4fee3b046395c4e751011e97f8900b5273d56685\"/>\n      </saml:SubjectConfirmation>\n    </saml:Subject>\n    <saml:Conditions NotBefore=\"2014-07-17T01:01:18Z\" NotOnOrAfter=\"2024-01-18T06:21:48Z\">\n      <saml:AudienceRestriction>\n        <saml:Audience>http://sp.example.com/demo1/metadata.php</saml:Audience>\n      </saml:AudienceRestriction>\n    </saml:Conditions>\n    <saml:AuthnStatement AuthnInstant=\"2014-07-17T01:01:48Z\" SessionNotOnOrAfter=\"2024-07-17T09:01:48Z\" SessionIndex=\"_be9967abd904ddcae3c0eb4189adbe3f71e327cf93\">\n      <saml:AuthnContext>\n        <saml:AuthnContextClassRef>urn:oasis:names:tc:SAML:2.0:ac:classes:Password</saml:AuthnContextClassRef>\n      </saml:AuthnContext>\n    </saml:AuthnStatement>\n    <saml:AttributeStatement>\n      <saml:Attribute Name=\"uid\" NameFormat=\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\">\n        <saml:AttributeValue xsi:type=\"xs:string\">test</saml:AttributeValue>\n      </saml:Attribute>\n      <saml:Attribute Name=\"mail\" NameFormat=\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\">\n        <saml:AttributeValue xsi:type=\"xs:string\">test@example.com</saml:AttributeValue>\n      </saml:Attribute>\n      <saml:Attribute Name=\"eduPersonAffiliation\" NameFormat=\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\">\n        <saml:AttributeValue xsi:type=\"xs:string\">users</saml:AttributeValue>\n        <saml:AttributeValue xsi:type=\"xs:string\">examplerole1</saml:AttributeValue>\n      </saml:Attribute>\n    </saml:AttributeStatement>\n  </saml:Assertion>\n</samlp:Response>",
				singatureAlgorithm: dsig.RSASHA1SignatureMethod,
				canonicalizer:      dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList(""),
			},
			res: res{
				canonicalization: dsig.CanonicalXML10ExclusiveAlgorithmId.String(),
				signatureMethod:  dsig.RSASHA1SignatureMethod,
				digestMethod:     xml_dsig.DigestMethodSHA1,
				transforms:       []string{dsig.EnvelopedSignatureAltorithmId.String(), dsig.CanonicalXML10ExclusiveAlgorithmId.String()},
				err:              false,
			},
		},
		{
			name: "SAMLTool example response SHA256",
			args: args{
				response:           "<samlp:Response xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" ID=\"_8e8dc5f69a98cc4c1ff3427e5ce34606fd672f91e6\" Version=\"2.0\" IssueInstant=\"2014-07-17T01:01:48Z\" Destination=\"http://sp.example.com/demo1/index.php?acs\" InResponseTo=\"ONELOGIN_4fee3b046395c4e751011e97f8900b5273d56685\">\n  <saml:Issuer>http://idp.example.com/metadata.php</saml:Issuer>\n  <samlp:Status>\n    <samlp:StatusCode Value=\"urn:oasis:names:tc:SAML:2.0:status:Success\"/>\n  </samlp:Status>\n  <saml:Assertion xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xs=\"http://www.w3.org/2001/XMLSchema\" ID=\"_d71a3a8e9fcc45c9e9d248ef7049393fc8f04e5f75\" Version=\"2.0\" IssueInstant=\"2014-07-17T01:01:48Z\">\n    <saml:Issuer>http://idp.example.com/metadata.php</saml:Issuer>\n    <saml:Subject>\n      <saml:NameID SPNameQualifier=\"http://sp.example.com/demo1/metadata.php\" Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\">_ce3d2948b4cf20146dee0a0b3dd6f69b6cf86f62d7</saml:NameID>\n      <saml:SubjectConfirmation Method=\"urn:oasis:names:tc:SAML:2.0:cm:bearer\">\n        <saml:SubjectConfirmationData NotOnOrAfter=\"2024-01-18T06:21:48Z\" Recipient=\"http://sp.example.com/demo1/index.php?acs\" InResponseTo=\"ONELOGIN_4fee3b046395c4e751011e97f8900b5273d56685\"/>\n      </saml:SubjectConfirmation>\n    </saml:Subject>\n    <saml:Conditions NotBefore=\"2014-07-17T01:01:18Z\" NotOnOrAfter=\"2024-01-18T06:21:48Z\">\n      <saml:AudienceRestriction>\n        <saml:Audience>http://sp.example.com/demo1/metadata.php</saml:Audience>\n      </saml:AudienceRestriction>\n    </saml:Conditions>\n    <saml:AuthnStatement AuthnInstant=\"2014-07-17T01:01:48Z\" SessionNotOnOrAfter=\"2024-07-17T09:01:48Z\" SessionIndex=\"_be9967abd904ddcae3c0eb4189adbe3f71e327cf93\">\n      <saml:AuthnContext>\n        <saml:AuthnContextClassRef>urn:oasis:names:tc:SAML:2.0:ac:classes:Password</saml:AuthnContextClassRef>\n      </saml:AuthnContext>\n    </saml:AuthnStatement>\n    <saml:AttributeStatement>\n      <saml:Attribute Name=\"uid\" NameFormat=\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\">\n        <saml:AttributeValue xsi:type=\"xs:string\">test</saml:AttributeValue>\n      </saml:Attribute>\n      <saml:Attribute Name=\"mail\" NameFormat=\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\">\n        <saml:AttributeValue xsi:type=\"xs:string\">test@example.com</saml:AttributeValue>\n      </saml:Attribute>\n      <saml:Attribute Name=\"eduPersonAffiliation\" NameFormat=\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\">\n        <saml:AttributeValue xsi:type=\"xs:string\">users</saml:AttributeValue>\n        <saml:AttributeValue xsi:type=\"xs:string\">examplerole1</saml:AttributeValue>\n      </saml:Attribute>\n    </saml:AttributeStatement>\n  </saml:Assertion>\n</samlp:Response>",
				singatureAlgorithm: dsig.RSASHA256SignatureMethod,
				canonicalizer:      dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList(""),
			},
			res: res{
				canonicalization: dsig.CanonicalXML10ExclusiveAlgorithmId.String(),
				signatureMethod:  dsig.RSASHA256SignatureMethod,
				digestMethod:     xml_dsig.DigestMethodSHA256,
				transforms:       []string{dsig.EnvelopedSignatureAltorithmId.String(), dsig.CanonicalXML10ExclusiveAlgorithmId.String()},
				err:              false,
			},
		},
		/* currently not supported by the xmlsig-library
		{
				name: "SAMLTool example response SHA512",
				args: args{
					response:           "<samlp:Response xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" ID=\"_8e8dc5f69a98cc4c1ff3427e5ce34606fd672f91e6\" Version=\"2.0\" IssueInstant=\"2014-07-17T01:01:48Z\" Destination=\"http://sp.example.com/demo1/index.php?acs\" InResponseTo=\"ONELOGIN_4fee3b046395c4e751011e97f8900b5273d56685\">\n  <saml:Issuer>http://idp.example.com/metadata.php</saml:Issuer>\n  <samlp:Status>\n    <samlp:StatusCode Value=\"urn:oasis:names:tc:SAML:2.0:status:Success\"/>\n  </samlp:Status>\n  <saml:Assertion xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xs=\"http://www.w3.org/2001/XMLSchema\" ID=\"_d71a3a8e9fcc45c9e9d248ef7049393fc8f04e5f75\" Version=\"2.0\" IssueInstant=\"2014-07-17T01:01:48Z\">\n    <saml:Issuer>http://idp.example.com/metadata.php</saml:Issuer>\n    <saml:Subject>\n      <saml:NameID SPNameQualifier=\"http://sp.example.com/demo1/metadata.php\" Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\">_ce3d2948b4cf20146dee0a0b3dd6f69b6cf86f62d7</saml:NameID>\n      <saml:SubjectConfirmation Method=\"urn:oasis:names:tc:SAML:2.0:cm:bearer\">\n        <saml:SubjectConfirmationData NotOnOrAfter=\"2024-01-18T06:21:48Z\" Recipient=\"http://sp.example.com/demo1/index.php?acs\" InResponseTo=\"ONELOGIN_4fee3b046395c4e751011e97f8900b5273d56685\"/>\n      </saml:SubjectConfirmation>\n    </saml:Subject>\n    <saml:Conditions NotBefore=\"2014-07-17T01:01:18Z\" NotOnOrAfter=\"2024-01-18T06:21:48Z\">\n      <saml:AudienceRestriction>\n        <saml:Audience>http://sp.example.com/demo1/metadata.php</saml:Audience>\n      </saml:AudienceRestriction>\n    </saml:Conditions>\n    <saml:AuthnStatement AuthnInstant=\"2014-07-17T01:01:48Z\" SessionNotOnOrAfter=\"2024-07-17T09:01:48Z\" SessionIndex=\"_be9967abd904ddcae3c0eb4189adbe3f71e327cf93\">\n      <saml:AuthnContext>\n        <saml:AuthnContextClassRef>urn:oasis:names:tc:SAML:2.0:ac:classes:Password</saml:AuthnContextClassRef>\n      </saml:AuthnContext>\n    </saml:AuthnStatement>\n    <saml:AttributeStatement>\n      <saml:Attribute Name=\"uid\" NameFormat=\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\">\n        <saml:AttributeValue xsi:type=\"xs:string\">test</saml:AttributeValue>\n      </saml:Attribute>\n      <saml:Attribute Name=\"mail\" NameFormat=\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\">\n        <saml:AttributeValue xsi:type=\"xs:string\">test@example.com</saml:AttributeValue>\n      </saml:Attribute>\n      <saml:Attribute Name=\"eduPersonAffiliation\" NameFormat=\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\">\n        <saml:AttributeValue xsi:type=\"xs:string\">users</saml:AttributeValue>\n        <saml:AttributeValue xsi:type=\"xs:string\">examplerole1</saml:AttributeValue>\n      </saml:Attribute>\n    </saml:AttributeStatement>\n  </saml:Assertion>\n</samlp:Response>",
					singatureAlgorithm: dsig.RSASHA512SignatureMethod,
					canonicalizer:      dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList(""),
				},
				res: res{
					canonicalization: dsig.CanonicalXML10ExclusiveAlgorithmId.String(),
					signatureMethod:  dsig.RSASHA512SignatureMethod,
					digestMethod:     xml_dsig.DigestMethodSHA512,
					transforms:       []string{dsig.EnvelopedSignatureAltorithmId.String(), dsig.CanonicalXML10ExclusiveAlgorithmId.String()},
					err:              false,
				},
			},*/
		{
			name: "SAMLTool example response unknown signing method",
			args: args{
				response:           "<samlp:Response xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" ID=\"_8e8dc5f69a98cc4c1ff3427e5ce34606fd672f91e6\" Version=\"2.0\" IssueInstant=\"2014-07-17T01:01:48Z\" Destination=\"http://sp.example.com/demo1/index.php?acs\" InResponseTo=\"ONELOGIN_4fee3b046395c4e751011e97f8900b5273d56685\">\n  <saml:Issuer>http://idp.example.com/metadata.php</saml:Issuer>\n  <samlp:Status>\n    <samlp:StatusCode Value=\"urn:oasis:names:tc:SAML:2.0:status:Success\"/>\n  </samlp:Status>\n  <saml:Assertion xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xs=\"http://www.w3.org/2001/XMLSchema\" ID=\"_d71a3a8e9fcc45c9e9d248ef7049393fc8f04e5f75\" Version=\"2.0\" IssueInstant=\"2014-07-17T01:01:48Z\">\n    <saml:Issuer>http://idp.example.com/metadata.php</saml:Issuer>\n    <saml:Subject>\n      <saml:NameID SPNameQualifier=\"http://sp.example.com/demo1/metadata.php\" Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\">_ce3d2948b4cf20146dee0a0b3dd6f69b6cf86f62d7</saml:NameID>\n      <saml:SubjectConfirmation Method=\"urn:oasis:names:tc:SAML:2.0:cm:bearer\">\n        <saml:SubjectConfirmationData NotOnOrAfter=\"2024-01-18T06:21:48Z\" Recipient=\"http://sp.example.com/demo1/index.php?acs\" InResponseTo=\"ONELOGIN_4fee3b046395c4e751011e97f8900b5273d56685\"/>\n      </saml:SubjectConfirmation>\n    </saml:Subject>\n    <saml:Conditions NotBefore=\"2014-07-17T01:01:18Z\" NotOnOrAfter=\"2024-01-18T06:21:48Z\">\n      <saml:AudienceRestriction>\n        <saml:Audience>http://sp.example.com/demo1/metadata.php</saml:Audience>\n      </saml:AudienceRestriction>\n    </saml:Conditions>\n    <saml:AuthnStatement AuthnInstant=\"2014-07-17T01:01:48Z\" SessionNotOnOrAfter=\"2024-07-17T09:01:48Z\" SessionIndex=\"_be9967abd904ddcae3c0eb4189adbe3f71e327cf93\">\n      <saml:AuthnContext>\n        <saml:AuthnContextClassRef>urn:oasis:names:tc:SAML:2.0:ac:classes:Password</saml:AuthnContextClassRef>\n      </saml:AuthnContext>\n    </saml:AuthnStatement>\n    <saml:AttributeStatement>\n      <saml:Attribute Name=\"uid\" NameFormat=\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\">\n        <saml:AttributeValue xsi:type=\"xs:string\">test</saml:AttributeValue>\n      </saml:Attribute>\n      <saml:Attribute Name=\"mail\" NameFormat=\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\">\n        <saml:AttributeValue xsi:type=\"xs:string\">test@example.com</saml:AttributeValue>\n      </saml:Attribute>\n      <saml:Attribute Name=\"eduPersonAffiliation\" NameFormat=\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\">\n        <saml:AttributeValue xsi:type=\"xs:string\">users</saml:AttributeValue>\n        <saml:AttributeValue xsi:type=\"xs:string\">examplerole1</saml:AttributeValue>\n      </saml:Attribute>\n    </saml:AttributeStatement>\n  </saml:Assertion>\n</samlp:Response>",
				singatureAlgorithm: "unknown",
				canonicalizer:      dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList(""),
			},
			res: res{
				canonicalization: dsig.CanonicalXML10ExclusiveAlgorithmId.String(),
				signatureMethod:  dsig.RSASHA512SignatureMethod,
				digestMethod:     xml_dsig.DigestMethodSHA512,
				transforms:       []string{dsig.EnvelopedSignatureAltorithmId.String(), dsig.CanonicalXML10ExclusiveAlgorithmId.String()},
				err:              true,
			},
		},
		{
			name: "SAMLTool example response structure-error",
			args: args{
				response:           "samlp:Response xmlns:samlp=\"urn:oasis:names:tc:SAML:2.0:protocol\" xmlns:saml=\"urn:oasis:names:tc:SAML:2.0:assertion\" ID=\"_8e8dc5f69a98cc4c1ff3427e5ce34606fd672f91e6\" Version=\"2.0\" IssueInstant=\"2014-07-17T01:01:48Z\" Destination=\"http://sp.example.com/demo1/index.php?acs\" InResponseTo=\"ONELOGIN_4fee3b046395c4e751011e97f8900b5273d56685\">\n  <saml:Issuer>http://idp.example.com/metadata.php</saml:Issuer>\n  <samlp:Status>\n    <samlp:StatusCode Value=\"urn:oasis:names:tc:SAML:2.0:status:Success\"/>\n  </samlp:Status>\n  <saml:Assertion xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xs=\"http://www.w3.org/2001/XMLSchema\" ID=\"_d71a3a8e9fcc45c9e9d248ef7049393fc8f04e5f75\" Version=\"2.0\" IssueInstant=\"2014-07-17T01:01:48Z\">\n    <saml:Issuer>http://idp.example.com/metadata.php</saml:Issuer>\n    <saml:Subject>\n      <saml:NameID SPNameQualifier=\"http://sp.example.com/demo1/metadata.php\" Format=\"urn:oasis:names:tc:SAML:2.0:nameid-format:transient\">_ce3d2948b4cf20146dee0a0b3dd6f69b6cf86f62d7</saml:NameID>\n      <saml:SubjectConfirmation Method=\"urn:oasis:names:tc:SAML:2.0:cm:bearer\">\n        <saml:SubjectConfirmationData NotOnOrAfter=\"2024-01-18T06:21:48Z\" Recipient=\"http://sp.example.com/demo1/index.php?acs\" InResponseTo=\"ONELOGIN_4fee3b046395c4e751011e97f8900b5273d56685\"/>\n      </saml:SubjectConfirmation>\n    </saml:Subject>\n    <saml:Conditions NotBefore=\"2014-07-17T01:01:18Z\" NotOnOrAfter=\"2024-01-18T06:21:48Z\">\n      <saml:AudienceRestriction>\n        <saml:Audience>http://sp.example.com/demo1/metadata.php</saml:Audience>\n      </saml:AudienceRestriction>\n    </saml:Conditions>\n    <saml:AuthnStatement AuthnInstant=\"2014-07-17T01:01:48Z\" SessionNotOnOrAfter=\"2024-07-17T09:01:48Z\" SessionIndex=\"_be9967abd904ddcae3c0eb4189adbe3f71e327cf93\">\n      <saml:AuthnContext>\n        <saml:AuthnContextClassRef>urn:oasis:names:tc:SAML:2.0:ac:classes:Password</saml:AuthnContextClassRef>\n      </saml:AuthnContext>\n    </saml:AuthnStatement>\n    <saml:AttributeStatement>\n      <saml:Attribute Name=\"uid\" NameFormat=\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\">\n        <saml:AttributeValue xsi:type=\"xs:string\">test</saml:AttributeValue>\n      </saml:Attribute>\n      <saml:Attribute Name=\"mail\" NameFormat=\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\">\n        <saml:AttributeValue xsi:type=\"xs:string\">test@example.com</saml:AttributeValue>\n      </saml:Attribute>\n      <saml:Attribute Name=\"eduPersonAffiliation\" NameFormat=\"urn:oasis:names:tc:SAML:2.0:attrname-format:basic\">\n        <saml:AttributeValue xsi:type=\"xs:string\">users</saml:AttributeValue>\n        <saml:AttributeValue xsi:type=\"xs:string\">examplerole1</saml:AttributeValue>\n      </saml:Attribute>\n    </saml:AttributeStatement>\n  </saml:Assertion>\n</samlp:Response>",
				singatureAlgorithm: dsig.RSASHA512SignatureMethod,
				canonicalizer:      dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList(""),
			},
			res: res{
				canonicalization: dsig.CanonicalXML10ExclusiveAlgorithmId.String(),
				signatureMethod:  dsig.RSASHA512SignatureMethod,
				digestMethod:     xml_dsig.DigestMethodSHA512,
				transforms:       []string{dsig.EnvelopedSignatureAltorithmId.String(), dsig.CanonicalXML10ExclusiveAlgorithmId.String()},
				err:              true,
			},
		},
	}
	caKeyGen, caCertGen, err := newCACertAndKey()
	if err != nil {
		t.Error("Create() failed to create ca cert")
		return
	}
	caCertGenData, err := crypto.BytesToCertificate(caCertGen)
	if err != nil {
		t.Error("Create() failed to decode ca cert")
		return
	}
	keyGen, certGen, err := newCertAndKey(caKeyGen, caCertGenData)
	if err != nil {
		t.Error("Create() failed to create cert")
		return
	}
	certGenData, err := crypto.BytesToCertificate(certGen)
	if err != nil {
		t.Error("Create() failed to decode cert")
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &samlp.ResponseType{}
			respBytes := []byte(tt.args.response)
			decoder := xml.NewDecoder(bytes.NewReader(respBytes))
			if err = decoder.Decode(resp); err != nil {
				if err := xml.Unmarshal(respBytes, resp); err != nil {
					if (err != nil) != tt.res.err {
						t.Errorf("Create() failed to unmarshall response")
					}
					return
				}
			}

			_, signer, err := signature.GetSigningContextAndSigner(certGenData, keyGen, tt.args.singatureAlgorithm)
			if err != nil {
				if (err != nil) != tt.res.err {
					t.Errorf("Create() failed to create signContext and signer")
				}
				return
			}

			sig, err := signature.Create(signer, resp)
			if err != nil {
				if (err != nil) != tt.res.err {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.res.err)
				}
				return
			}
			resp.Signature = sig

			respStr, err := saml_xml.Marshal(resp)
			if err != nil {
				if (err != nil) != tt.res.err {
					t.Errorf("Create() marshall response for signing")
				}
				return
			}

			doc := etree.NewDocument()
			if err := doc.ReadFromBytes([]byte(respStr)); err != nil {
				if (err != nil) != tt.res.err {
					t.Errorf("Cert() failed to read response")
				}
				return
			}
			certStr := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(string(certGen), "-----BEGIN CERTIFICATE-----", ""), "-----END CERTIFICATE-----", ""), "\n", "")

			certs, err := signature.ParseCertificates([]string{certStr})
			if err != nil {
				if (err != nil) != tt.res.err {
					t.Errorf("Create() failed to parse certificate")
				}
				return
			}
			if err := signature.ValidatePost(certs, doc.Root()); err != nil {
				if (err != nil) != tt.res.err {
					t.Errorf("Create() failed to validate created signature")
				}
				return
			}
		})
	}
}

func newCertAndKey(
	caPrivateKey *rsa.PrivateKey,
	caCertificate []byte,
) (*rsa.PrivateKey, []byte, error) {
	now := time.Now().UTC()
	after := now.Add(time.Minute * 5)

	privateCrypto, _, certificateCrypto, err := crypto.GenerateCertificate(4096, caPrivateKey, caCertificate, &crypto.CertificateInformations{
		SerialNumber: big.NewInt(int64(rand.Intn(50000))),
		Organisation: []string{"caos AG"},
		CommonName:   "test",
		NotBefore:    &now,
		NotAfter:     &after,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	})
	if err != nil {
		return nil, nil, err
	}
	return privateCrypto, certificateCrypto, nil

}

func newCACertAndKey() (
	*rsa.PrivateKey,
	[]byte,
	error,
) {
	now := time.Now().UTC()
	after := now.Add(time.Minute * 5)
	privateCrypto, _, certificateCrypto, err := crypto.GenerateCACertificate(4096, &crypto.CertificateInformations{
		SerialNumber: big.NewInt(int64(rand.Intn(50000))),
		Organisation: []string{"caos AG"},
		CommonName:   "Zitadel SAML CA",
		NotBefore:    &now,
		NotAfter:     &after,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
	})
	if err != nil {
		return nil, nil, err
	}
	return privateCrypto, certificateCrypto, nil

}
