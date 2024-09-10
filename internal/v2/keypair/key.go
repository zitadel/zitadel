package keypair

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
)

type Key struct {
	Key    *crypto.CryptoValue `json:"key"`
	Expiry time.Time           `json:"expiry"`
}
