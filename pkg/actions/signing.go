package actions

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNoValidSignature = errors.New("no valid signature")
	ErrInvalidHeader    = errors.New("webhook has invalid Zitadel-Signature header")
	ErrNotSigned        = errors.New("webhook has no Zitadel-Signature header")
	ErrTooOld           = errors.New("timestamp wasn't within tolerance")
)

const (
	SigningHeader           = "ZITADEL-Signature"
	signingTimestamp        = "t"
	signingVersion   string = "v1"
	DefaultTolerance        = 300 * time.Second
	partSeparator           = ","
)

func ComputeSignatureHeader(t time.Time, payload []byte, signingKey ...string) string {
	parts := []string{
		fmt.Sprintf("%s=%d", signingTimestamp, t.Unix()),
	}
	for _, k := range signingKey {
		parts = append(parts, fmt.Sprintf("%s=%s", signingVersion, hex.EncodeToString(computeSignature(t, payload, k))))
	}
	return strings.Join(parts, partSeparator)
}

func computeSignature(t time.Time, payload []byte, signingKey string) []byte {
	mac := hmac.New(sha256.New, []byte(signingKey))
	mac.Write([]byte(fmt.Sprintf("%d", t.Unix())))
	mac.Write([]byte("."))
	mac.Write(payload)
	return mac.Sum(nil)
}

func ValidatePayload(payload []byte, header string, signingKey string) error {
	return ValidatePayloadWithTolerance(payload, header, signingKey, DefaultTolerance)
}

func ValidatePayloadWithTolerance(payload []byte, header string, signingKey string, tolerance time.Duration) error {
	return validatePayload(payload, header, signingKey, tolerance, true)
}

func validatePayload(payload []byte, sigHeader string, signingKey string, tolerance time.Duration, enforceTolerance bool) error {
	header, err := parseSignatureHeader(sigHeader)
	if err != nil {
		return err
	}

	expectedSignature := computeSignature(header.timestamp, payload, signingKey)
	expiredTimestamp := time.Since(header.timestamp) > tolerance
	if enforceTolerance && expiredTimestamp {
		return ErrTooOld
	}

	for _, sig := range header.signatures {
		if hmac.Equal(expectedSignature, sig) {
			return nil
		}
	}
	return ErrNoValidSignature
}

type signedHeader struct {
	timestamp  time.Time
	signatures [][]byte
}

func parseSignatureHeader(header string) (*signedHeader, error) {
	sh := &signedHeader{}
	if header == "" {
		return sh, ErrNotSigned
	}

	pairs := strings.Split(header, ",")
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			return sh, ErrInvalidHeader
		}
		switch parts[0] {
		case signingTimestamp:
			timestamp, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return sh, ErrInvalidHeader
			}
			sh.timestamp = time.Unix(timestamp, 0)

		case signingVersion:
			sig, err := hex.DecodeString(parts[1])
			if err != nil {
				continue
			}
			sh.signatures = append(sh.signatures, sig)
		default:
			continue
		}
	}

	if len(sh.signatures) == 0 {
		return sh, ErrNoValidSignature
	}
	return sh, nil
}
