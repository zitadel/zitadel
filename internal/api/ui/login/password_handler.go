package login

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/pbkdf2"

	"github.com/zitadel/zitadel/internal/domain"
)

const (
	tmplPassword = "password"

	// passwordEncryptedPartCount is the number of dot-separated components in
	// an encrypted password payload: ciphertext.tag.iv.salt
	passwordEncryptedPartCount = 4
	// pbkdf2Iterations matches the client-side iteration count in password_encrypt.js
	pbkdf2Iterations = 100_000
	// pbkdf2KeyLen is the AES-256 key length in bytes
	pbkdf2KeyLen = 32
	// aesGCMNonceSize is the expected IV length for AES-GCM
	aesGCMNonceSize = 12
)

type passwordFormData struct {
	Password string `schema:"password"`
}

func (l *Login) renderPassword(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getUserData(r, authReq, translator, "Password.Title", "Password.Description", err)
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
		// The browser encrypted the password with AES-256-GCM, keyed via
		// PBKDF2-SHA256 using the auth-request ID as the passphrase.
		// Silently fall back to the raw value if decryption fails so that
		// a misconfigured client still reaches VerifyPassword and fails
		// there rather than producing an opaque 500.
		if plaintext, decErr := decryptLoginPassword(data.Password, authReq.ID); decErr == nil {
			data.Password = plaintext
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

// decryptLoginPassword reverses the AES-256-GCM encryption applied by
// password_encrypt.js. The payload format is hex(ciphertext).hex(tag).hex(iv).hex(salt),
// and the key is derived from the authRequestID passphrase using PBKDF2-SHA256.
func decryptLoginPassword(payload, authRequestID string) (string, error) {
	parts := strings.SplitN(payload, ".", passwordEncryptedPartCount)
	if len(parts) != passwordEncryptedPartCount {
		return "", fmt.Errorf("invalid encrypted password format: expected %d parts, got %d", passwordEncryptedPartCount, len(parts))
	}

	ciphertext, err := hex.DecodeString(parts[0])
	if err != nil {
		return "", fmt.Errorf("decode ciphertext: %w", err)
	}
	tag, err := hex.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("decode tag: %w", err)
	}
	iv, err := hex.DecodeString(parts[2])
	if err != nil {
		return "", fmt.Errorf("decode iv: %w", err)
	}
	salt, err := hex.DecodeString(parts[3])
	if err != nil {
		return "", fmt.Errorf("decode salt: %w", err)
	}

	if len(iv) != aesGCMNonceSize {
		return "", fmt.Errorf("invalid IV length: got %d bytes, expected %d", len(iv), aesGCMNonceSize)
	}

	key := pbkdf2.Key([]byte(authRequestID), salt, pbkdf2Iterations, pbkdf2KeyLen, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create GCM: %w", err)
	}

	// GCM Open expects ciphertext + tag appended
	combined := append(ciphertext, tag...) //nolint:gocritic // intentional append to new slice
	plaintext, err := aead.Open(nil, iv, combined, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}
	return string(plaintext), nil
}
