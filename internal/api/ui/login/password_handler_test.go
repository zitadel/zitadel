package login

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"
	"testing"
)

// ── helpers ──────────────────────────────────────────────────────────────────

// encryptForTest mirrors the algorithm in password_encrypt.js so we can
// produce valid test payloads without a browser.
func encryptForTest(t *testing.T, password string, serverPubKey *ecdh.PublicKey) string {
	t.Helper()
	clientPriv := mustGenerateECDHKey(t)
	aesKey := ecdhDerivedAESKey(t, clientPriv, serverPubKey)
	aead := mustNewGCM(t, aesKey)
	iv := mustRandBytes(t, 12)
	sealed := aead.Seal(nil, iv, []byte(password), nil)
	ciphertext := sealed[:len(sealed)-16]
	tag := sealed[len(sealed)-16:]
	clientPubB64 := base64.StdEncoding.EncodeToString(clientPriv.PublicKey().Bytes())
	return clientPubB64 + "." + hex.EncodeToString(ciphertext) + "." + hex.EncodeToString(tag) + "." + hex.EncodeToString(iv)
}

func mustGenerateECDHKey(t *testing.T) *ecdh.PrivateKey {
	t.Helper()
	k, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal("generate ECDH key:", err)
	}
	return k
}

func ecdhDerivedAESKey(t *testing.T, priv *ecdh.PrivateKey, pub *ecdh.PublicKey) [32]byte {
	t.Helper()
	shared, err := priv.ECDH(pub)
	if err != nil {
		t.Fatal("ECDH:", err)
	}
	return sha256.Sum256(shared)
}

func mustNewGCM(t *testing.T, aesKey [32]byte) cipher.AEAD {
	t.Helper()
	block, err := aes.NewCipher(aesKey[:])
	if err != nil {
		t.Fatal("new cipher:", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatal("new GCM:", err)
	}
	return aead
}

func mustRandBytes(t *testing.T, n int) []byte {
	t.Helper()
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		t.Fatal("rand bytes:", err)
	}
	return b
}

func generateServerKeyPair(t *testing.T) (*ecdh.PrivateKey, *ecdh.PublicKey) {
	t.Helper()
	priv, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal("generate server key:", err)
	}
	return priv, priv.PublicKey()
}

// mutatePayload splits a valid payload, applies fn to the parts slice, and
// rejoins. Keeps all per-case mutation logic in a single shared helper.
func mutatePayload(t *testing.T, password string, serverPub *ecdh.PublicKey, fn func(parts []string) []string) string {
	t.Helper()
	parts := strings.SplitN(encryptForTest(t, password, serverPub), ".", 4)
	parts = fn(parts)
	return strings.Join(parts, ".")
}

// ── decryptLoginPassword round-trip tests ────────────────────────────────────

func Test_decryptLoginPassword_roundTrip(t *testing.T) {
	serverPriv, serverPub := generateServerKeyPair(t)

	cases := []struct {
		name     string
		password string
	}{
		{"ascii", "MyP@ssw0rd!"},
		{"unicode", "pässwörد"},
		{"empty", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload := encryptForTest(t, tc.password, serverPub)
			got, err := decryptLoginPassword(payload, serverPriv)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.password {
				t.Errorf("got %q, want %q", got, tc.password)
			}
		})
	}
}

// ── decryptLoginPassword error path tests ────────────────────────────────────

func Test_decryptLoginPassword_errors(t *testing.T) {
	serverPriv, serverPub := generateServerKeyPair(t)
	wrongPriv, _ := generateServerKeyPair(t)

	cases := []struct {
		name    string
		payload string
		key     *ecdh.PrivateKey
	}{
		{
			name:    "wrong server key",
			payload: encryptForTest(t, "secret", serverPub),
			key:     wrongPriv,
		},
		{
			name:    "too few parts",
			payload: "aabb.ccdd",
			key:     serverPriv,
		},
		{
			name:    "invalid base64 client pubkey",
			payload: mutatePayload(t, "x", serverPub, func(p []string) []string { p[0] = "!!!notbase64!!!"; return p }),
			key:     serverPriv,
		},
		{
			name:    "invalid hex ciphertext",
			payload: mutatePayload(t, "x", serverPub, func(p []string) []string { p[1] = strings.Repeat("zz", len(p[1])/2); return p }),
			key:     serverPriv,
		},
		{
			name:    "tag too short",
			payload: mutatePayload(t, "x", serverPub, func(p []string) []string { p[2] = "aabb"; return p }),
			key:     serverPriv,
		},
		{
			name:    "IV too short",
			payload: mutatePayload(t, "x", serverPub, func(p []string) []string { p[3] = "aabb"; return p }),
			key:     serverPriv,
		},
		{
			name:    "ciphertext too long",
			payload: mutatePayload(t, "x", serverPub, func(p []string) []string { p[1] = strings.Repeat("aa", encMaxCiphertextHexLen/2+1); return p }),
			key:     serverPriv,
		},
		{
			name:    "client pubkey too long",
			payload: mutatePayload(t, "x", serverPub, func(p []string) []string { p[0] = strings.Repeat("A", encMaxPubKeyBase64Len+1); return p }),
			key:     serverPriv,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := decryptLoginPassword(tc.payload, tc.key)
			if err == nil {
				t.Errorf("expected error for %q, got nil", tc.name)
			}
		})
	}
}

// ── passwordEncKeyStore tests ─────────────────────────────────────────────────

func Test_passwordEncKeyStore_generateAndRetrieve(t *testing.T) {
	store := newPasswordEncKeyStore()
	defer store.close()

	pubKey, err := store.generate("req-1")
	if err != nil {
		t.Fatal(err)
	}
	if pubKey == nil {
		t.Fatal("expected non-nil public key")
	}
	priv, err := store.retrieve("req-1")
	if err != nil {
		t.Fatal(err)
	}
	if priv == nil {
		t.Fatal("expected non-nil private key")
	}
}

func Test_passwordEncKeyStore_singleUse(t *testing.T) {
	store := newPasswordEncKeyStore()
	defer store.close()

	if _, err := store.generate("req-2"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.retrieve("req-2"); err != nil {
		t.Fatal("first retrieve should succeed:", err)
	}
	if _, err := store.retrieve("req-2"); err == nil {
		t.Error("expected error on second retrieve — key must be deleted after first use")
	}
}

func Test_passwordEncKeyStore_unknownKey(t *testing.T) {
	store := newPasswordEncKeyStore()
	defer store.close()

	if _, err := store.retrieve("nonexistent-req"); err == nil {
		t.Error("expected error for unknown auth request ID")
	}
}

func Test_passwordEncKeyStore_keypairConsistency(t *testing.T) {
	store := newPasswordEncKeyStore()
	defer store.close()

	pubKey, err := store.generate("req-3")
	if err != nil {
		t.Fatal(err)
	}
	privKey, err := store.retrieve("req-3")
	if err != nil {
		t.Fatal(err)
	}

	// Verify via ECDH: both sides must derive the same shared secret.
	clientPriv, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	serverShared, err := privKey.ECDH(clientPriv.PublicKey())
	if err != nil {
		t.Fatal("server ECDH:", err)
	}
	clientShared, err := clientPriv.ECDH(pubKey)
	if err != nil {
		t.Fatal("client ECDH:", err)
	}
	if string(serverShared) != string(clientShared) {
		t.Error("server and client shared secrets do not match — keypair is inconsistent")
	}
}
