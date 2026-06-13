package login

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/zitadel/zitadel/internal/domain"
)

const (
	tmplPassword = "password"

	// passwordEncryptedPartCount is the number of dot-separated components:
	// base64(clientPubKey) . hex(ciphertext) . hex(gcm_tag) . hex(iv)
	passwordEncryptedPartCount = 4

	// Fixed sizes for AES-256-GCM components — enforced before ECDH/AES runs
	// to prevent DoS via large attacker-controlled allocations.
	encTagHexLen = 32  // 16-byte GCM tag → 32 hex chars
	encIVHexLen  = 24  // 12-byte IV → 24 hex chars
	// Ciphertext is password-length + GCM overhead; cap at a generous bound.
	// Max password length in Zitadel is 72 chars (bcrypt limit); 256 hex chars
	// (128 bytes) is more than sufficient. Empty password produces 0-length ciphertext.
	encMaxCiphertextHexLen = 256
	// P-256 uncompressed public key is 65 bytes; base64 = 88 chars.
	encMaxPubKeyBase64Len = 88
)

type passwordFormData struct {
	Password string `schema:"password"`
}

func (l *Login) renderPassword(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getUserData(r, authReq, translator, "Password.Title", "Password.Description", err)

	if l.config.PasswordEncryption.Enabled && authReq != nil {
		// Generate an ephemeral ECDH keypair for this auth request.
		// The public key is embedded in the page; the private key stays server-side.
		pubKey, genErr := l.encKeyStore.generate(authReq.ID)
		if genErr == nil {
			data.PasswordEncryptionEnabled = true
			data.PasswordEncryptionPublicKey = base64.StdEncoding.EncodeToString(pubKey.Bytes())
		}
		// On key generation failure we silently fall back to plaintext submission.
		// The feature is not available but auth still works normally.
	}

	funcs := map[string]interface{}{
		"showPasswordReset": func() bool {
			if authReq.LoginPolicy != nil {
				return !authReq.LoginPolicy.HidePasswordReset
			}
			return true
		},
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplPassword], data, funcs)
}

func (l *Login) handlePasswordCheck(w http.ResponseWriter, r *http.Request) {
	data := new(passwordFormData)
	authReq, err := l.ensureAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	if l.config.PasswordEncryption.Enabled {
		privKey, keyErr := l.encKeyStore.retrieve(authReq.ID)
		if keyErr != nil {
			// No ephemeral key found for this auth request — either it expired,
			// was already consumed, or the client never received one.
			if !l.config.PasswordEncryption.AllowPlaintextFallback {
				// Strict mode: reject the submission. Render password page with error.
				l.renderPassword(w, r, authReq, fmt.Errorf("encrypted password required"))
				return
			}
			// Fallback mode: proceed with raw password; VerifyPassword will
			// reject it if incorrect.
		} else {
			plaintext, decErr := decryptLoginPassword(data.Password, privKey)
			if decErr != nil {
				if !l.config.PasswordEncryption.AllowPlaintextFallback {
					l.renderPassword(w, r, authReq, fmt.Errorf("encrypted password required"))
					return
				}
				// Fallback: use raw value; VerifyPassword will reject it.
			} else {
				data.Password = plaintext
			}
		}
	}

	err = l.authRepo.VerifyPassword(setContext(r.Context(), authReq.UserOrgID), authReq.ID, authReq.UserID, authReq.UserOrgID, data.Password, authReq.AgentID, domain.BrowserInfoFromRequest(r))

	metadata, actionErr := l.runPostInternalAuthenticationActions(authReq, r, authMethodPassword, err)
	if err == nil && actionErr == nil && len(metadata) > 0 {
		err = l.bulkSetUserMetadata(r.Context(), authReq.UserID, authReq.UserOrgID, metadata)
	} else if actionErr != nil && err == nil {
		err = actionErr
	}

	if err != nil {
		if authReq.LoginPolicy.IgnoreUnknownUsernames {
			l.renderLogin(w, r, authReq, err)
			return
		}
		l.renderPassword(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

// decryptLoginPassword decrypts a payload produced by password_encrypt.js.
//
// The payload format is dot-separated:
//
//	base64(clientP256PubKey) . hex(ciphertext) . hex(gcm_tag) . hex(iv)
//
// The shared secret is derived via P-256 ECDH between the server's ephemeral
// private key and the client's ephemeral public key. The 32-byte shared secret
// is used directly as the AES-256-GCM key (after SHA-256 of the raw ECDH
// shared-X coordinate, matching the JS SubtleCrypto deriveBits path).
//
// Because the server private key is stored server-side only and deleted on
// first retrieval, a captured POST body — containing only the client public key
// and ciphertext — is insufficient to decrypt without the server private key.
func decryptLoginPassword(payload string, serverPrivKey *ecdh.PrivateKey) (string, error) {
	parts := strings.SplitN(payload, ".", passwordEncryptedPartCount)
	if len(parts) != passwordEncryptedPartCount {
		return "", fmt.Errorf("invalid encrypted password: expected %d dot-separated parts, got %d", passwordEncryptedPartCount, len(parts))
	}

	clientPubKeyB64, ciphertextHex, tagHex, ivHex := parts[0], parts[1], parts[2], parts[3]

	// Enforce fixed sizes for constant-size fields before any allocation.
	if len(tagHex) != encTagHexLen {
		return "", fmt.Errorf("invalid GCM tag length: got %d hex chars, expected %d", len(tagHex), encTagHexLen)
	}
	if len(ivHex) != encIVHexLen {
		return "", fmt.Errorf("invalid IV length: got %d hex chars, expected %d", len(ivHex), encIVHexLen)
	}
	if len(ciphertextHex) > encMaxCiphertextHexLen {
		return "", fmt.Errorf("ciphertext length out of bounds: %d hex chars", len(ciphertextHex))
	}
	if len(clientPubKeyB64) > encMaxPubKeyBase64Len {
		return "", fmt.Errorf("client public key too long: %d chars", len(clientPubKeyB64))
	}

	// Decode client public key.
	clientPubKeyBytes, err := base64.StdEncoding.DecodeString(clientPubKeyB64)
	if err != nil {
		return "", fmt.Errorf("decode client public key: %w", err)
	}
	clientPubKey, err := ecdh.P256().NewPublicKey(clientPubKeyBytes)
	if err != nil {
		return "", fmt.Errorf("parse client public key: %w", err)
	}

	// ECDH: derive shared secret.
	sharedSecret, err := serverPrivKey.ECDH(clientPubKey)
	if err != nil {
		return "", fmt.Errorf("ECDH: %w", err)
	}
	// Derive AES-256 key from the raw shared secret bytes via SHA-256.
	// This matches the Web Crypto API deriveBits("ECDH") + SHA-256 digest path
	// used in password_encrypt.js.
	aesKey := sha256.Sum256(sharedSecret)

	// Decode hex fields.
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", fmt.Errorf("decode ciphertext: %w", err)
	}
	tag, err := hex.DecodeString(tagHex)
	if err != nil {
		return "", fmt.Errorf("decode GCM tag: %w", err)
	}
	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		return "", fmt.Errorf("decode IV: %w", err)
	}

	// AES-256-GCM decrypt.
	block, err := aes.NewCipher(aesKey[:])
	if err != nil {
		return "", fmt.Errorf("create AES cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create GCM: %w", err)
	}

	// Web Crypto AES-GCM appends the tag to the ciphertext in the encrypted buffer.
	combined := append(ciphertext, tag...) //nolint:gocritic
	plaintext, err := aead.Open(nil, iv, combined, nil)
	if err != nil {
		return "", fmt.Errorf("AES-GCM decrypt: %w", err)
	}
	return string(plaintext), nil
}
