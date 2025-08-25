package target

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type TargetType uint

const (
	TargetTypeWebhook TargetType = iota
	TargetTypeCall
	TargetTypeAsync
)

type Target struct {
	ExecutionID      string              `json:"execution_id,omitempty"`
	TargetID         string              `json:"target_id,omitempty"`
	TargetType       TargetType          `json:"target_type,omitempty"`
	Endpoint         string              `json:"endpoint,omitempty"`
	Timeout          time.Duration       `json:"timeout,omitempty"`
	InterruptOnError bool                `json:"interrupt_on_error,omitempty"`
	SigningKey       *crypto.CryptoValue `json:"signing_key,omitempty"`
}

func (e *Target) DecryptSigningKey(alg crypto.EncryptionAlgorithm) (string, error) {
	if e.SigningKey == nil {
		return "", nil
	}
	keyValue, err := crypto.DecryptString(e.SigningKey, alg)
	if err != nil {
		return "", zerrors.ThrowInternal(err, "QUERY-bxevy3YXwy", "Errors.Internal")
	}
	return keyValue, nil
}
