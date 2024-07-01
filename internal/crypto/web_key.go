package crypto

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"

	"github.com/go-jose/go-jose/v4"
	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/zerrors"
)

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

type WebKeyConfigType int

const (
	WebKeyConfigTypeRSA WebKeyConfigType = iota
	WebKeyConfigTypeECDSA
	WebKeyConfigTypeED25519
)

type RSABits int

const (
	RSABits2048 RSABits = 2048
	RSABits3072 RSABits = 3072
	RSABits4096 RSABits = 4096
)

type RSAHasher int

const (
	RSAHasherSHA256 RSAHasher = iota
	RSAHasherSHA384
	RSAHasherSHA512
)

type EllipticCurve int

const (
	EllipticCurveP256 EllipticCurve = iota
	EllipticCurveP384
	EllipticCurveP512
)

type WebKeyConfig interface {
	Alg() jose.SignatureAlgorithm
	Type() WebKeyConfigType // Type is needed to make Unmarshal work
}

func UnmarshalWebKeyConfig(data []byte, configType WebKeyConfigType) (config WebKeyConfig, err error) {
	switch configType {
	case WebKeyConfigTypeRSA:
		config = new(WebKeyRSAConfig)
	case WebKeyConfigTypeECDSA:
		config = new(WebKeyECDSAConfig)
	case WebKeyConfigTypeED25519:
		config = new(WebKeyED25519Config)
	default:
		return nil, zerrors.ThrowInternal(nil, "CRYPT-Eig8ho", "Errors.Internal")
	}
	if err = json.Unmarshal(data, config); err != nil {
		return nil, zerrors.ThrowInternal(nil, "CRYPT-waeR0N", "Errors.Internal")
	}
	return config, nil
}

type WebKeyRSAConfig struct {
	Bits   RSABits
	Hasher RSAHasher
}

func (c WebKeyRSAConfig) Alg() jose.SignatureAlgorithm {
	switch c.Hasher {
	case RSAHasherSHA256:
		return jose.RS256
	case RSAHasherSHA384:
		return jose.RS384
	case RSAHasherSHA512:
		return jose.RS512
	default:
		return jose.RS256
	}
}

func (WebKeyRSAConfig) Type() WebKeyConfigType {
	return WebKeyConfigTypeRSA
}

type WebKeyECDSAConfig struct {
	Curve EllipticCurve
}

func (c WebKeyECDSAConfig) Alg() jose.SignatureAlgorithm {
	switch c.Curve {
	case EllipticCurveP256:
		return jose.ES256
	case EllipticCurveP384:
		return jose.ES384
	case EllipticCurveP512:
		return jose.ES512
	default:
		return jose.ES256
	}
}

func (WebKeyECDSAConfig) Type() WebKeyConfigType {
	return WebKeyConfigTypeECDSA
}

func (c WebKeyECDSAConfig) GetCurve() elliptic.Curve {
	switch c.Curve {
	case EllipticCurveP256:
		return elliptic.P256()
	case EllipticCurveP384:
		return elliptic.P384()
	case EllipticCurveP512:
		return elliptic.P521()
	default:
		return elliptic.P256()
	}
}

type WebKeyED25519Config struct{}

func (WebKeyED25519Config) Alg() jose.SignatureAlgorithm {
	return jose.EdDSA
}

func (WebKeyED25519Config) Type() WebKeyConfigType {
	return WebKeyConfigTypeED25519
}

func GenerateEncryptedWebKey(keyID string, alg EncryptionAlgorithm, genConfig WebKeyConfig) (encryptedPrivate *CryptoValue, public *jose.JSONWebKey, err error) {
	return generateEncryptedWebKey(rand.Reader, keyID, alg, genConfig)
}

func generateEncryptedWebKey(reader io.Reader, keyID string, alg EncryptionAlgorithm, genConfig WebKeyConfig) (encryptedWebKey *CryptoValue, public *jose.JSONWebKey, err error) {
	var key any
	switch conf := genConfig.(type) {
	case WebKeyRSAConfig:
		key, err = rsa.GenerateKey(reader, int(conf.Bits))
	case WebKeyECDSAConfig:
		key, err = ecdsa.GenerateKey(conf.GetCurve(), reader)
	case WebKeyED25519Config:
		_, key, err = ed25519.GenerateKey(reader)
	default:
		err = zerrors.ThrowInternalf(nil, "CRYPT-aeW6x", "Errors.Internal")
	}
	if err != nil {
		return nil, nil, err
	}

	webKey := newJSONWebkey(key, keyID, genConfig.Alg())
	encryptedWebKey, err = EncryptJSON(webKey, alg)
	if err != nil {
		return nil, nil, err
	}
	return encryptedWebKey, gu.Ptr(webKey.Public()), err
}

func newJSONWebkey(key any, keyID string, algorithm jose.SignatureAlgorithm) *jose.JSONWebKey {
	return &jose.JSONWebKey{
		Key:       key,
		KeyID:     keyID,
		Algorithm: string(algorithm),
		Use:       KeyUsageSigning.String(),
	}
}
