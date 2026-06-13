package login

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strings"
	"testing"

	"golang.org/x/crypto/pbkdf2"
)

// encryptForTest mirrors the algorithm in password_encrypt.js so we can
// produce test vectors without a browser.
func encryptForTest(t *testing.T, password, authRequestID string) string {
	t.Helper()

	salt := make([]byte, 16)
	iv := make([]byte, aesGCMNonceSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		t.Fatal(err)
	}

	key := pbkdf2.Key([]byte(authRequestID), salt, pbkdf2Iterations, pbkdf2KeyLen, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatal(err)
	}

	// Seal appends ciphertext+tag
	sealed := aead.Seal(nil, iv, []byte(password), nil)
	ciphertext := sealed[:len(sealed)-16]
	tag := sealed[len(sealed)-16:]

	return hex.EncodeToString(ciphertext) + "." +
		hex.EncodeToString(tag) + "." +
		hex.EncodeToString(iv) + "." +
		hex.EncodeToString(salt)
}

func Test_decryptLoginPassword(t *testing.T) {
	const authRequestID = "test-auth-request-id-123"

	tests := []struct {
		name        string
		password    string
		authID      string
		wantErr     bool
	}{
		{
			name:     "simple ascii password",
			password: "MyP@ssw0rd!",
			authID:   authRequestID,
		},
		{
			name:     "unicode password",
			password: "pässwörد",
			authID:   authRequestID,
		},
		{
			name:     "empty password",
			password: "",
			authID:   authRequestID,
		},
		{
			name:    "wrong auth request id",
			password: "secret",
			authID:  "wrong-id",
			wantErr: true,
		},
		{
			name:    "malformed payload - too few parts",
			password: "aabbcc.ddeeff",
			authID:  authRequestID,
			wantErr: true,
		},
		{
			name:    "malformed payload - invalid hex",
			password: "zzzz.aabb.ccdd.eeff",
			authID:  authRequestID,
			wantErr: true,
		},
		{
			name:    "plaintext password not accepted",
			password: "plaintext-no-dots",
			authID:  authRequestID,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var payload string
			if !tt.wantErr || tt.name == "wrong auth request id" {
				// Encrypt with the correct key; wrong-ID test decrypts with wrong key
				payload = encryptForTest(t, tt.password, authRequestID)
			} else {
				payload = tt.password // pass the raw string as-is for format error cases
			}

			got, err := decryptLoginPassword(payload, tt.authID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("decryptLoginPassword() expected error, got nil (result=%q)", got)
				}
				return
			}
			if err != nil {
				t.Errorf("decryptLoginPassword() unexpected error: %v", err)
				return
			}
			if got != tt.password {
				t.Errorf("decryptLoginPassword() = %q, want %q", got, tt.password)
			}
		})
	}
}

func Test_decryptLoginPassword_payloadFormat(t *testing.T) {
	// Verify the expected dot-separated hex format is enforced
	cases := []string{
		"",
		"single",
		"a.b",
		"a.b.c",
		"a.b.c.d.e", // too many parts — SplitN(4) keeps last part intact, still 4 parts
	}
	for _, c := range cases {
		// The last case (5 apparent segments) actually passes SplitN with n=4,
		// because SplitN stops at 4 parts; test only the genuinely bad ones.
		parts := strings.SplitN(c, ".", passwordEncryptedPartCount)
		if len(parts) == passwordEncryptedPartCount {
			continue // format check passes, hex decode will fail separately
		}
		_, err := decryptLoginPassword(c, "any-id")
		if err == nil {
			t.Errorf("decryptLoginPassword(%q) expected format error, got nil", c)
		}
	}
}
