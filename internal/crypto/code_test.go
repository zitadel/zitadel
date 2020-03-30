package crypto

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Encrypted_OK(t *testing.T) {
	mCrypto := NewMockEncryptionAlgorithm(gomock.NewController(t))
	mCrypto.EXPECT().Algorithm().AnyTimes().Return("enc")
	mCrypto.EXPECT().EncryptionKeyID().AnyTimes().Return("id")
	mCrypto.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{"id"})
	mCrypto.EXPECT().Encrypt(gomock.Any()).DoAndReturn(
		func(code []byte) ([]byte, error) {
			return code, nil
		},
	)
	mCrypto.EXPECT().DecryptString(gomock.Any(), gomock.Any()).DoAndReturn(
		func(code []byte, _ string) (string, error) {
			return string(code), nil
		},
	)
	generator := NewEncryptionGenerator(6, 0, mCrypto, Digits)

	crypto, code, err := NewCode(generator)
	assert.NoError(t, err)

	decrypted, err := DecryptString(crypto, generator.alg)
	assert.NoError(t, err)
	assert.Equal(t, code, decrypted)
	assert.Equal(t, 6, len(decrypted))
}

func Test_Verify_Encrypted_OK(t *testing.T) {
	mCrypto := NewMockEncryptionAlgorithm(gomock.NewController(t))
	mCrypto.EXPECT().Algorithm().AnyTimes().Return("enc")
	mCrypto.EXPECT().EncryptionKeyID().AnyTimes().Return("id")
	mCrypto.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{"id"})
	mCrypto.EXPECT().DecryptString(gomock.Any(), gomock.Any()).DoAndReturn(
		func(code []byte, _ string) (string, error) {
			return string(code), nil
		},
	)
	creationDate := time.Now()
	code := &CryptoValue{
		CryptoType: TypeEncryption,
		Algorithm:  "enc",
		KeyID:      "id",
		Crypted:    []byte("code"),
	}
	generator := NewEncryptionGenerator(6, 0, mCrypto, Digits)

	err := VerifyCode(creationDate, 1*time.Hour, code, "code", generator)
	assert.NoError(t, err)
}
func Test_Verify_Encrypted_Invalid_Err(t *testing.T) {
	mCrypto := NewMockEncryptionAlgorithm(gomock.NewController(t))
	mCrypto.EXPECT().Algorithm().AnyTimes().Return("enc")
	mCrypto.EXPECT().EncryptionKeyID().AnyTimes().Return("id")
	mCrypto.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{"id"})
	mCrypto.EXPECT().DecryptString(gomock.Any(), gomock.Any()).DoAndReturn(
		func(code []byte, _ string) (string, error) {
			return string(code), nil
		},
	)
	creationDate := time.Now()
	code := &CryptoValue{
		CryptoType: TypeEncryption,
		Algorithm:  "enc",
		KeyID:      "id",
		Crypted:    []byte("code"),
	}
	generator := NewEncryptionGenerator(6, 0, mCrypto, Digits)

	err := VerifyCode(creationDate, 1*time.Hour, code, "wrong", generator)
	assert.Error(t, err)
}

func TestIsCodeExpired_Expired(t *testing.T) {
	creationDate := time.Date(2019, time.April, 1, 0, 0, 0, 0, time.UTC)
	expired := IsCodeExpired(creationDate, 1*time.Hour)
	assert.True(t, expired)
}

func TestIsCodeExpired_NotExpired(t *testing.T) {
	creationDate := time.Now()
	expired := IsCodeExpired(creationDate, 1*time.Hour)
	assert.False(t, expired)
}
