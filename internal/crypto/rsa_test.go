package crypto

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"
)

func mustGenerateRSAPEM(t *testing.T) []byte {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
}

func mustGenerateECDSAPEM(t *testing.T) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
}

func mustGenerateEd25519PEM(t *testing.T) []byte {
	t.Helper()
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
}

func TestBytesToPublicKey_RSA(t *testing.T) {
	pemData := mustGenerateRSAPEM(t)
	key, err := BytesToPublicKey(pemData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := key.(*rsa.PublicKey); !ok {
		t.Fatalf("expected *rsa.PublicKey, got %T", key)
	}
}

func TestBytesToPublicKey_ECDSA(t *testing.T) {
	pemData := mustGenerateECDSAPEM(t)
	key, err := BytesToPublicKey(pemData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := key.(*ecdsa.PublicKey); !ok {
		t.Fatalf("expected *ecdsa.PublicKey, got %T", key)
	}
}

func TestBytesToPublicKey_Ed25519(t *testing.T) {
	pemData := mustGenerateEd25519PEM(t)
	key, err := BytesToPublicKey(pemData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := key.(ed25519.PublicKey); !ok {
		t.Fatalf("expected ed25519.PublicKey, got %T", key)
	}
}

func TestBytesToPublicKey_Empty(t *testing.T) {
	_, err := BytesToPublicKey(nil)
	if !errors.Is(err, ErrEmpty) {
		t.Fatalf("expected ErrEmpty, got %v", err)
	}
	_, err = BytesToPublicKey([]byte{})
	if !errors.Is(err, ErrEmpty) {
		t.Fatalf("expected ErrEmpty, got %v", err)
	}
}

func TestBytesToPublicKey_InvalidPEM(t *testing.T) {
	_, err := BytesToPublicKey([]byte("not a PEM block"))
	if !errors.Is(err, ErrEmpty) {
		t.Fatalf("expected ErrEmpty for non-PEM data, got %v", err)
	}
}

func TestBytesToPublicKey_MalformedDER(t *testing.T) {
	malformed := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: []byte("invalid DER data")})
	_, err := BytesToPublicKey(malformed)
	if err == nil {
		t.Fatal("expected parse error for malformed DER, got nil")
	}
}
