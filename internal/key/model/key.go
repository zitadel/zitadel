package model

import (
	"time"

	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type KeyPair struct {
	es_models.ObjectRoot

	Usage      KeyUsage
	Algorithm  string
	PrivateKey *Key
	PublicKey  *Key
}

type KeyUsage int32

const (
	KeyUsageSigning KeyUsage = iota
)

type Key struct {
	Key    *crypto.CryptoValue
	Expiry time.Time
}
