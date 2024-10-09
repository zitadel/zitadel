package actions

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

var (
	ErrNoValidSignature = errors.New("no valid signature")
)

const (
	SigningHeader = "ZITADEL-Signature"
)

func ComputeSignature(payload []byte, signingKey string) string {
	return base64.RawStdEncoding.EncodeToString(computeSignature(payload, signingKey))
}

func computeSignature(payload []byte, signingKey string) []byte {
	mac := hmac.New(sha256.New, []byte(signingKey))
	mac.Write(payload)
	return mac.Sum(nil)
}

func ValidatePayload(payload []byte, sigHeader string, signingKey string) error {
	expectedSignature := computeSignature(payload, signingKey)
	decoded, err := base64.RawStdEncoding.DecodeString(sigHeader)
	if err != nil {
		return ErrNoValidSignature
	}
	if hmac.Equal(expectedSignature, decoded) {
		return nil
	}
	return ErrNoValidSignature
}
