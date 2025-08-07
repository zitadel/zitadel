package crypto

import (
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func CreateMockEncryptionAlg(ctrl *gomock.Controller) EncryptionAlgorithm {
	return createMockEncryptionAlgorithm(
		ctrl,
		func(code []byte) ([]byte, error) {
			return code, nil
		},
	)
}

// CreateMockEncryptionAlgWithCode compares the length of the value to be encrypted with the length of the provided code.
// It will return an error if they do not match.
// The provided code will be used to encrypt in favor of the value passed to the encryption.
// This function is intended to be used where the passed value is not in control, but where the returned encryption requires a static value.
func CreateMockEncryptionAlgWithCode(ctrl *gomock.Controller, code string) EncryptionAlgorithm {
	return createMockEncryptionAlgorithm(
		ctrl,
		func(c []byte) ([]byte, error) {
			if len(c) != len(code) {
				return nil, zerrors.ThrowInvalidArgumentf(nil, "id", "invalid code length - expected %d, got %d", len(code), len(c))
			}
			return []byte(code), nil
		},
	)
}

func createMockEncryptionAlgorithm(ctrl *gomock.Controller, encryptFunction func(c []byte) ([]byte, error)) *MockEncryptionAlgorithm {
	mCrypto := NewMockEncryptionAlgorithm(ctrl)
	mCrypto.EXPECT().Algorithm().AnyTimes().Return("enc")
	mCrypto.EXPECT().EncryptionKeyID().AnyTimes().Return("id")
	mCrypto.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{"id"})
	mCrypto.EXPECT().Encrypt(gomock.Any()).AnyTimes().DoAndReturn(
		encryptFunction,
	)
	mCrypto.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
		func(code []byte, keyID string) (string, error) {
			if keyID != "id" {
				return "", zerrors.ThrowInternal(nil, "id", "invalid key id")
			}
			return string(code), nil
		},
	)
	mCrypto.EXPECT().Decrypt(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
		func(code []byte, keyID string) ([]byte, error) {
			if keyID != "id" {
				return nil, zerrors.ThrowInternal(nil, "id", "invalid key id")
			}
			return code, nil
		},
	)
	return mCrypto
}

func createMockCrypto(t *testing.T) EncryptionAlgorithm {
	mCrypto := NewMockEncryptionAlgorithm(gomock.NewController(t))
	mCrypto.EXPECT().Algorithm().AnyTimes().Return("crypto")
	return mCrypto
}

func createMockGenerator(t *testing.T, crypto EncryptionAlgorithm) Generator {
	mGenerator := NewMockGenerator(gomock.NewController(t))
	mGenerator.EXPECT().Alg().AnyTimes().Return(crypto)
	return mGenerator
}
