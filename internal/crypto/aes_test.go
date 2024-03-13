package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: refactor test style
func TestDecrypt_OK(t *testing.T) {
	encryptedpw, err := EncryptAESString("ThisIsMySecretPw", "passphrasewhichneedstobe32bytes!")
	assert.NoError(t, err)

	decryptedpw, err := DecryptAESString(encryptedpw, "passphrasewhichneedstobe32bytes!")
	assert.NoError(t, err)

	assert.Equal(t, "ThisIsMySecretPw", decryptedpw)
}
