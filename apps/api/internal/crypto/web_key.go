package crypto

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"

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

//go:generate enumer -type WebKeyConfigType -trimprefix WebKeyConfigType -text -json -linecomment
type WebKeyConfigType int

const (
	WebKeyConfigTypeUnspecified WebKeyConfigType = iota //
	WebKeyConfigTypeRSA
	WebKeyConfigTypeECDSA
	WebKeyConfigTypeED25519
)

//go:generate enumer -type RSABits -trimprefix RSABits -text -json -linecomment
type RSABits int

const (
	RSABitsUnspecified RSABits = 0 //
	RSABits2048        RSABits = 2048
	RSABits3072        RSABits = 3072
	RSABits4096        RSABits = 4096
)

type RSAHasher int

//go:generate enumer -type RSAHasher -trimprefix RSAHasher -text -json -linecomment
const (
	RSAHasherUnspecified RSAHasher = iota //
	RSAHasherSHA256
	RSAHasherSHA384
	RSAHasherSHA512
)

type EllipticCurve int

//go:generate enumer -type EllipticCurve -trimprefix EllipticCurve -text -json -linecomment
const (
	EllipticCurveUnspecified EllipticCurve = iota //
	EllipticCurveP256
	EllipticCurveP384
	EllipticCurveP512
)

type WebKeyConfig interface {
	Alg() jose.SignatureAlgorithm
	Type() WebKeyConfigType // Type is needed to make Unmarshal work
	IsValid() error
}

func UnmarshalWebKeyConfig(data []byte, configType WebKeyConfigType) (config WebKeyConfig, err error) {
	switch configType {
	case WebKeyConfigTypeUnspecified:
		return nil, zerrors.ThrowInternal(nil, "CRYPT-Ii3AiH", "Errors.Internal")
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
		return nil, zerrors.ThrowInternal(err, "CRYPT-waeR0N", "Errors.Internal")
	}
	return config, nil
}

type WebKeyRSAConfig struct {
	Bits   RSABits
	Hasher RSAHasher
}

func (c WebKeyRSAConfig) Alg() jose.SignatureAlgorithm {
	switch c.Hasher {
	case RSAHasherUnspecified:
		return ""
	case RSAHasherSHA256:
		return jose.RS256
	case RSAHasherSHA384:
		return jose.RS384
	case RSAHasherSHA512:
		return jose.RS512
	default:
		return ""
	}
}

func (WebKeyRSAConfig) Type() WebKeyConfigType {
	return WebKeyConfigTypeRSA
}

func (c WebKeyRSAConfig) IsValid() error {
	if !c.Bits.IsARSABits() || c.Bits == RSABitsUnspecified {
		return zerrors.ThrowInvalidArgument(nil, "CRYPTO-eaz3T", "Errors.WebKey.Config")
	}
	if !c.Hasher.IsARSAHasher() || c.Hasher == RSAHasherUnspecified {
		return zerrors.ThrowInvalidArgument(nil, "CRYPTO-ODie7", "Errors.WebKey.Config")
	}
	return nil
}

type WebKeyECDSAConfig struct {
	Curve EllipticCurve
}

func (c WebKeyECDSAConfig) Alg() jose.SignatureAlgorithm {
	switch c.Curve {
	case EllipticCurveUnspecified:
		return ""
	case EllipticCurveP256:
		return jose.ES256
	case EllipticCurveP384:
		return jose.ES384
	case EllipticCurveP512:
		return jose.ES512
	default:
		return ""
	}
}

func (WebKeyECDSAConfig) Type() WebKeyConfigType {
	return WebKeyConfigTypeECDSA
}

func (c WebKeyECDSAConfig) IsValid() error {
	if !c.Curve.IsAEllipticCurve() || c.Curve == EllipticCurveUnspecified {
		return zerrors.ThrowInvalidArgument(nil, "CRYPTO-Ii2ai", "Errors.WebKey.Config")
	}
	return nil
}

func (c WebKeyECDSAConfig) GetCurve() elliptic.Curve {
	switch c.Curve {
	case EllipticCurveUnspecified:
		return nil
	case EllipticCurveP256:
		return elliptic.P256()
	case EllipticCurveP384:
		return elliptic.P384()
	case EllipticCurveP512:
		return elliptic.P521()
	default:
		return nil
	}
}

type WebKeyED25519Config struct{}

func (WebKeyED25519Config) Alg() jose.SignatureAlgorithm {
	return jose.EdDSA
}

func (WebKeyED25519Config) Type() WebKeyConfigType {
	return WebKeyConfigTypeED25519
}

func (WebKeyED25519Config) IsValid() error {
	return nil
}

func GenerateEncryptedWebKey(keyID string, alg EncryptionAlgorithm, genConfig WebKeyConfig) (encryptedPrivate *CryptoValue, public *jose.JSONWebKey, err error) {
	private, public, err := generateWebKey(keyID, genConfig)
	if err != nil {
		return nil, nil, err
	}
	encryptedPrivate, err = EncryptJSON(private, alg)
	if err != nil {
		return nil, nil, err
	}
	return encryptedPrivate, public, nil
}

func generateWebKey(keyID string, genConfig WebKeyConfig) (private, public *jose.JSONWebKey, err error) {
	if err = genConfig.IsValid(); err != nil {
		return nil, nil, err
	}
	var key any
	switch conf := genConfig.(type) {
	case *WebKeyRSAConfig:
		key, err = rsa.GenerateKey(rand.Reader, int(conf.Bits))
	case *WebKeyECDSAConfig:
		key, err = ecdsa.GenerateKey(conf.GetCurve(), rand.Reader)
	case *WebKeyED25519Config:
		_, key, err = ed25519.GenerateKey(rand.Reader)
	default:
		return nil, nil, fmt.Errorf("unknown webkey config type %T", genConfig)
	}
	if err != nil {
		return nil, nil, err
	}

	private = newJSONWebkey(key, keyID, genConfig.Alg())
	return private, gu.Ptr(private.Public()), err
}

func newJSONWebkey(key any, keyID string, algorithm jose.SignatureAlgorithm) *jose.JSONWebKey {
	return &jose.JSONWebKey{
		Key:       key,
		KeyID:     keyID,
		Algorithm: string(algorithm),
		Use:       KeyUsageSigning.String(),
	}
}
