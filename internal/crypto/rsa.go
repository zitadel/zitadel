package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privkey, &privkey.PublicKey, nil
}

func GenerateEncryptedKeyPair(bits int, alg EncryptionAlgorithm) (*CryptoValue, *CryptoValue, error) {
	privateKey, publicKey, err := GenerateKeyPair(bits)
	if err != nil {
		return nil, nil, err
	}
	return EncryptKeys(privateKey, publicKey, alg)
}

func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)
}

func PublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes, nil
}

func BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}
	return key, nil
}

var ErrEmpty = errors.New("cannot decode, empty data")

func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	if pub == nil {
		return nil, ErrEmpty
	}
	block, _ := pem.Decode(pub)
	if block == nil {
		return nil, ErrEmpty
	}
	ifc, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, err
	}
	return key, nil
}

func EncryptKeys(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, alg EncryptionAlgorithm) (*CryptoValue, *CryptoValue, error) {
	encryptedPrivateKey, err := Encrypt(PrivateKeyToBytes(privateKey), alg)
	if err != nil {
		return nil, nil, err
	}
	pubKey, err := PublicKeyToBytes(publicKey)
	if err != nil {
		return nil, nil, err
	}
	encryptedPublicKey, err := Encrypt(pubKey, alg)
	if err != nil {
		return nil, nil, err
	}
	return encryptedPrivateKey, encryptedPublicKey, nil
}
