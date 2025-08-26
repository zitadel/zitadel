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
	ExecutionID      string              `json:"execution_id"`
	TargetID         string              `json:"target_id"`
	TargetType       TargetType          `json:"target_type"`
	Endpoint         string              `json:"endpoint"`
	Timeout          time.Duration       `json:"timeout"`
	InterruptOnError bool                `json:"interrupt_on_error"`
	SigningKey       *crypto.CryptoValue `json:"signing_key"`
	SigningKeyDec    string
}

func (e *Target) GetExecutionID() string {
	return e.ExecutionID
}
func (e *Target) GetTargetID() string {
	return e.TargetID
}
func (e *Target) IsInterruptOnError() bool {
	return e.InterruptOnError
}
func (e *Target) GetEndpoint() string {
	return e.Endpoint
}
func (e *Target) GetTargetType() TargetType {
	return e.TargetType
}
func (e *Target) GetTimeout() time.Duration {
	return e.Timeout
}
func (e *Target) GetSigningKey(alg crypto.EncryptionAlgorithm) string {
	if e.SigningKeyDec == "" && e.SigningKey != nil {
		e.decryptSigningKey(alg)
	}
	return e.SigningKeyDec
}

func (t *Target) decryptSigningKey(alg crypto.EncryptionAlgorithm) error {
	if t.SigningKey == nil {
		return nil
	}
	keyValue, err := crypto.DecryptString(t.SigningKey, alg)
	if err != nil {
		return zerrors.ThrowInternal(err, "QUERY-bxevy3YXwy", "Errors.Internal")
	}
	t.SigningKeyDec = keyValue
	return nil
}
