package pem

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func EncodeCertificate(data []byte) ([]byte, error) {
	certPem := new(bytes.Buffer)
	if err := pem.Encode(certPem, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: data,
	}); err != nil {
		return nil, err
	}
	return certPem.Bytes(), nil
}

func EncodeKey(key *rsa.PrivateKey) ([]byte, error) {
	keyPem := new(bytes.Buffer)
	if err := pem.Encode(keyPem, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}); err != nil {
		return nil, err
	}
	return keyPem.Bytes(), nil
}

func DecodeKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func DecodeCertificate(data []byte) ([]byte, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	return block.Bytes, nil
}
