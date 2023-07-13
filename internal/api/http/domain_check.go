package http

import (
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
	return errors.ThrowInvalidArgument(err, "HTTP-GH422", "Errors.Internal")
}

func ValidateDomainDNS(domain, verifier string) error {
	txtRecords, err := net.LookupTXT(tokenUrlDNS(domain))
	if err != nil {
		return errors.ThrowInternal(err, "HTTP-Hwsw2", "Errors.Internal")
	}

	for _, record := range txtRecords {
		if record == verifier {
			return nil
		}
	}
	return errors.ThrowInvalidArgument(err, "HTTP-G241f", "Errors.Internal")
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
