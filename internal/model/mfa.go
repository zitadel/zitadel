package model

import "github.com/caos/zitadel/internal/crypto"

type Multifactors struct {
	OTP OTP
}

type OTP struct {
	Issuer    string
	CryptoMFA crypto.EncryptionAlgorithm
}
