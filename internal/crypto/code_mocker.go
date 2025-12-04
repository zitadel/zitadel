package crypto

import (
	"testing"

	"go.uber.org/mock/gomock"

	cryptomock "github.com/zitadel/zitadel/internal/crypto/mock"
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

func createMockEncryptionAlgorithm(ctrl *gomock.Controller, encryptFunction func(c []byte) ([]byte, error)) *cryptomock.MockEncryptionAlgorithm {
	mCrypto := cryptomock.NewMockEncryptionAlgorithm(ctrl)
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
	mCrypto := cryptomock.NewMockEncryptionAlgorithm(gomock.NewController(t))
	mCrypto.EXPECT().Algorithm().AnyTimes().Return("crypto")
	return mCrypto
}

func createMockGenerator(t *testing.T, crypto EncryptionAlgorithm) Generator {
	mGenerator := NewMockGenerator(gomock.NewController(t))
	mGenerator.EXPECT().Alg().AnyTimes().Return(crypto)
	return mGenerator
}
