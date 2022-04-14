package signature_test

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
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
	"strings"
	"testing"
	"time"
)

func Test_SignatureValidate(t *testing.T) {
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

func Test_SignatureCreatePost(t *testing.T) {
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
