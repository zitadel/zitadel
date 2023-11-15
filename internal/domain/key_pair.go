package domain

import (
	"time"

	"github.com/zitadel/zitadel/v2/internal/crypto"
	es_models "github.com/zitadel/zitadel/v2/internal/eventstore/v1/models"
)

type KeyPair struct {
	es_models.ObjectRoot

	Usage       KeyUsage
	Algorithm   string
	PrivateKey  *Key
	PublicKey   *Key
	Certificate *Key
}

type KeyUsage int32

const (
	KeyUsageSigning KeyUsage = iota
	KeyUsageSAMLMetadataSigning
	KeyUsageSAMLResponseSinging
	KeyUsageSAMLCA
)

func (u KeyUsage) String() string {
	switch u {
	case KeyUsageSigning:
		return "sig"
	case KeyUsageSAMLCA:
		return "saml_ca"
	case KeyUsageSAMLResponseSinging:
		return "saml_response_sig"
	case KeyUsageSAMLMetadataSigning:
		return "saml_metadata_sig"
	}
	return ""
}

type Key struct {
	Key    *crypto.CryptoValue
	Expiry time.Time
}

func (k *KeyPair) IsValid() bool {
	return k.Algorithm != "" &&
		k.PrivateKey != nil && k.PrivateKey.IsValid() &&
		k.PublicKey != nil && k.PublicKey.IsValid() &&
		k.Certificate != nil && k.Certificate.IsValid()
}

func (k *Key) IsValid() bool {
	return k.Key != nil
}
