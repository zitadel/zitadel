package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type WebAuthNToken struct {
	es_models.ObjectRoot

	WebAuthNTokenID        string
	CredentialCreationData []byte
	State                  MfaState
	Challenge              string
	AllowedCredentialIDs   [][]byte
	UserVerification       UserVerificationRequirement
	KeyID                  []byte
	PublicKey              []byte
	AttestationType        string
	AAGUID                 []byte
	SignCount              uint32
}

type WebAuthNMethod int32

const (
	WebAuthNMethodUnspecified WebAuthNMethod = iota
	WebAuthNMethodU2F
	WebAuthNMethodPasswordless
)

type UserVerificationRequirement int32

const (
	UserVerificationRequirementUnspecified UserVerificationRequirement = iota
	UserVerificationRequirementRequired
	UserVerificationRequirementPreferred
	UserVerificationRequirementDiscouraged
)
