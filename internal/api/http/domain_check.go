package http

import (
	errorsAs "errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/zitadel/zitadel/internal/errors"
)

type CheckType int

const (
	CheckTypeHTTP CheckType = iota
	CheckTypeDNS

	HTTPPattern = "https://%s/.well-known/zitadel-challenge/%s.txt"
	DNSPattern  = "_zitadel-challenge.%s"
)

func ValidateDomain(domain, token, verifier string, checkType CheckType) error {
	switch checkType {
	case CheckTypeHTTP:
		return ValidateDomainHTTP(domain, token, verifier)
	case CheckTypeDNS:
		return ValidateDomainDNS(domain, verifier)
	default:
		return errors.ThrowInvalidArgument(nil, "HTTP-Iqd11", "Errors.Internal")
	}
}

func ValidateDomainHTTP(domain, token, verifier string) error {
	resp, err := http.Get(tokenUrlHTTP(domain, token))
	if err != nil {
		return errors.ThrowInternal(err, "HTTP-BH42h", "Errors.Internal")
	}
	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			return errors.ThrowNotFound(err, "ORG-F4zhw", "Errors.Org.DomainVerificationHTTPNotFound")
		}
		return errors.ThrowInternal(err, "HTTP-G2zsw", "Errors.Internal")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.ThrowInternal(err, "HTTP-HB432", "Errors.Internal")
	}
	if string(body) == verifier {
		return nil
	}
	return errors.ThrowNotFound(err, "ORG-GH422", "Errors.Org.DomainVerificationHTTPNoMatch")
}

func ValidateDomainDNS(domain, verifier string) error {
	txtRecords, err := net.LookupTXT(tokenUrlDNS(domain))
	if err != nil {
		var dnsError *net.DNSError
		if errorsAs.As(err, &dnsError) {
			if dnsError.IsNotFound {
				return errors.ThrowNotFound(err, "ORG-G241f", "Errors.Org.DomainVerificationTXTNotFound")
			}
			if dnsError.IsTimeout {
				return errors.ThrowNotFound(err, "ORG-K563l", "Errors.Org.DomainVerificationTimeout")
			}
		}
		return errors.ThrowInternal(err, "HTTP-Hwsw2", "Errors.Internal")
	}

	for _, record := range txtRecords {
		if record == verifier {
			return nil
		}
	}
	return errors.ThrowNotFound(err, "ORG-G28if", "Errors.Org.DomainVerificationTXTNoMatch")
}

func TokenUrl(domain, token string, checkType CheckType) (string, error) {
	switch checkType {
	case CheckTypeHTTP:
		return tokenUrlHTTP(domain, token), nil
	case CheckTypeDNS:
		return tokenUrlDNS(domain), nil
	default:
		return "", errors.ThrowInvalidArgument(nil, "HTTP-Iqd11", "")
	}
}

func tokenUrlHTTP(domain, token string) string {
	return fmt.Sprintf(HTTPPattern, domain, token)
}

func tokenUrlDNS(domain string) string {
	return fmt.Sprintf(DNSPattern, domain)
}
