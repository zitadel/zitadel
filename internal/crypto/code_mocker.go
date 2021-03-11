package crypto

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/mock/gomock"
	"testing"
)

func CreateMockEncryptionAlg(ctrl *gomock.Controller) EncryptionAlgorithm {
	mCrypto := NewMockEncryptionAlgorithm(ctrl)
	mCrypto.EXPECT().Algorithm().AnyTimes().Return("enc")
	mCrypto.EXPECT().EncryptionKeyID().AnyTimes().Return("id")
	mCrypto.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{"id"})
	mCrypto.EXPECT().Encrypt(gomock.Any()).AnyTimes().DoAndReturn(
		func(code []byte) ([]byte, error) {
			return code, nil
		},
	)
	mCrypto.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
		func(code []byte, keyID string) (string, error) {
			if keyID != "id" {
				return "", errors.ThrowInternal(nil, "id", "invalid key id")
			}
			return string(code), nil
		},
	)
	return mCrypto
}

func CreateMockHashAlg(ctrl *gomock.Controller) HashAlgorithm {
	mCrypto := NewMockHashAlgorithm(ctrl)
	mCrypto.EXPECT().Algorithm().AnyTimes().Return("hash")
	mCrypto.EXPECT().Hash(gomock.Any()).AnyTimes().DoAndReturn(
		func(code []byte) ([]byte, error) {
			return code, nil
		},
	)
	mCrypto.EXPECT().CompareHash(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
		func(hashed, comparer []byte) error {
			if string(hashed) != string(comparer) {
				return errors.ThrowInternal(nil, "id", "invalid")
			}
			return nil
		},
	)
	return mCrypto
}

func CreateMockCrypto(t *testing.T) Crypto {
	mCrypto := NewMockCrypto(gomock.NewController(t))
	mCrypto.EXPECT().Algorithm().AnyTimes().Return("crypto")
	return mCrypto
}

func CreateMockGenerator(t *testing.T, crypto Crypto) Generator {
	mGenerator := NewMockGenerator(gomock.NewController(t))
	mGenerator.EXPECT().Alg().AnyTimes().Return(crypto)
	mGenerator.EXPECT().Length().AnyTimes().Return(8)
	return mGenerator
}
