package receiver

import "github.com/zitadel/zitadel/internal/crypto"

type Verifiable struct {
	IsVerified bool
	Code       *crypto.CryptoValue
}
