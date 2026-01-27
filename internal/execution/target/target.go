package target

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
)

type TargetType uint

const (
	TargetTypeWebhook TargetType = iota
	TargetTypeCall
	TargetTypeAsync
)

type PayloadType uint

const (
	PayloadTypeUnspecified PayloadType = iota
	PayloadTypeJSON
	PayloadTypeJWT
	PayloadTypeJWE
)

type Target struct {
	ExecutionID      string              `json:"execution_id,omitempty"`
	TargetID         string              `json:"target_id,omitempty"`
	TargetType       TargetType          `json:"target_type,omitempty"`
	Endpoint         string              `json:"endpoint,omitempty"`
	Timeout          time.Duration       `json:"timeout,omitempty"`
	InterruptOnError bool                `json:"interrupt_on_error,omitempty"`
	SigningKey       *crypto.CryptoValue `json:"signing_key,omitempty"`
	PayloadType      PayloadType         `json:"payload_type,omitempty"`
	EncryptionKey    []byte              `json:"encryption_key,omitempty"`
	EncryptionKeyID  string              `json:"encryption_key_id,omitempty"`
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
func (e *Target) GetSigningKey(alg crypto.EncryptionAlgorithm) (string, error) {
	if e.SigningKey == nil {
		return "", nil
	}
	return crypto.DecryptString(e.SigningKey, alg)
}
func (e *Target) GetPayloadType() PayloadType {
	return e.PayloadType
}
func (e *Target) GetEncryptionKey() []byte {
	return e.EncryptionKey
}

func (e *Target) GetEncryptionKeyID() string {
	return e.EncryptionKeyID
}
