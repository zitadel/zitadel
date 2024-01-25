package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"
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

type CertificateInformations struct {
	SerialNumber *big.Int
	Organisation []string
	CommonName   string
	NotBefore    time.Time
	NotAfter     time.Time
	KeyUsage     x509.KeyUsage
	ExtKeyUsage  []x509.ExtKeyUsage
}

func GenerateEncryptedKeyPairWithCACertificate(bits int, keyAlg, certAlg EncryptionAlgorithm, informations *CertificateInformations) (*CryptoValue, *CryptoValue, *CryptoValue, error) {
	privateKey, publicKey, cert, err := GenerateCACertificate(bits, informations)
	if err != nil {
		return nil, nil, nil, err
	}
	encryptPriv, encryptPub, encryptCaCert, err := EncryptKeysAndCert(privateKey, publicKey, cert, keyAlg, certAlg)
	if err != nil {
		return nil, nil, nil, err
	}
	return encryptPriv, encryptPub, encryptCaCert, nil
}

func GenerateEncryptedKeyPairWithCertificate(bits int, keyAlg, certAlg EncryptionAlgorithm, caPrivateKey *rsa.PrivateKey, caCertificate []byte, informations *CertificateInformations) (*CryptoValue, *CryptoValue, *CryptoValue, error) {
	privateKey, publicKey, cert, err := GenerateCertificate(bits, caPrivateKey, caCertificate, informations)
	if err != nil {
		return nil, nil, nil, err
	}
	encryptPriv, encryptPub, encryptCaCert, err := EncryptKeysAndCert(privateKey, publicKey, cert, keyAlg, certAlg)
	if err != nil {
		return nil, nil, nil, err
	}
	return encryptPriv, encryptPub, encryptCaCert, nil
}

func GenerateCACertificate(bits int, informations *CertificateInformations) (*rsa.PrivateKey, *rsa.PublicKey, []byte, error) {
	return generateCertificate(bits, nil, nil, informations)
}

func GenerateCertificate(bits int, caPrivateKey *rsa.PrivateKey, ca []byte, informations *CertificateInformations) (*rsa.PrivateKey, *rsa.PublicKey, []byte, error) {
	return generateCertificate(bits, caPrivateKey, ca, informations)
}

func generateCertificate(bits int, caPrivateKey *rsa.PrivateKey, ca []byte, informations *CertificateInformations) (*rsa.PrivateKey, *rsa.PublicKey, []byte, error) {
	notBefore := time.Now()
	if !informations.NotBefore.IsZero() {
		notBefore = informations.NotBefore
	}
	cert := &x509.Certificate{
		SerialNumber: informations.SerialNumber,
		Subject: pkix.Name{
			CommonName:   informations.CommonName,
			Organization: informations.Organisation,
		},
		NotBefore:   notBefore,
		NotAfter:    informations.NotAfter,
		KeyUsage:    informations.KeyUsage,
		ExtKeyUsage: informations.ExtKeyUsage,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, nil, err
	}

	certBytes := make([]byte, 0)
	if ca == nil {
		cert.IsCA = true
		cert.BasicConstraintsValid = true

		certBytes, err = x509.CreateCertificate(rand.Reader, cert, cert, &certPrivKey.PublicKey, certPrivKey)
		if err != nil {
			return nil, nil, nil, err
		}
	} else {
		caCert, err := x509.ParseCertificate(ca)
		if err != nil {
			return nil, nil, nil, err
		}

		certBytes, err = x509.CreateCertificate(rand.Reader, cert, caCert, &certPrivKey.PublicKey, caPrivateKey)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	x509Cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, nil, nil, err
	}

	certPem, err := CertificateToBytes(x509Cert)
	if err != nil {
		return nil, nil, nil, err
	}

	return certPrivKey, &certPrivKey.PublicKey, certPem, nil
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
	if len(pub) == 0 {
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

func CertificateToBytes(cert *x509.Certificate) ([]byte, error) {
	certPem := new(bytes.Buffer)
	if err := pem.Encode(certPem, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}); err != nil {
		return nil, err
	}
	return certPem.Bytes(), nil
}

func BytesToCertificate(data []byte) ([]byte, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}
	return block.Bytes, nil
}

func EncryptKeysAndCert(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, cert []byte, keyAlg, certAlg EncryptionAlgorithm) (*CryptoValue, *CryptoValue, *CryptoValue, error) {
	encryptedPrivateKey, encryptedPublicKey, err := EncryptKeys(privateKey, publicKey, keyAlg)
	if err != nil {
		return nil, nil, nil, err
	}
	encryptedCertificate, err := Encrypt(cert, certAlg)
	if err != nil {
		return nil, nil, nil, err
	}
	return encryptedPrivateKey, encryptedPublicKey, encryptedCertificate, nil
}
