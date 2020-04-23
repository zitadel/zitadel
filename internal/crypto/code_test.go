package crypto

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/caos/zitadel/internal/errors"
)

func createMockEncryptionAlg(t *testing.T) EncryptionAlgorithm {
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
		func(code []byte, keyID string) (string, error) {
			if keyID != "id" {
				return "", errors.ThrowInternal(nil, "id", "invalid key id")
			}
			return string(code), nil
		},
	)
	return mCrypto
}

func createMockHashAlg(t *testing.T) HashAlgorithm {
	mCrypto := NewMockHashAlgorithm(gomock.NewController(t))
	mCrypto.EXPECT().Algorithm().AnyTimes().Return("hash")
	mCrypto.EXPECT().Hash(gomock.Any()).DoAndReturn(
		func(code []byte) ([]byte, error) {
			return code, nil
		},
	)
	mCrypto.EXPECT().CompareHash(gomock.Any(), gomock.Any()).DoAndReturn(
		func(hashed, comparer []byte) error {
			if string(hashed) != string(comparer) {
				return errors.ThrowInternal(nil, "id", "invalid")
			}
			return nil
		},
	)
	return mCrypto
}

func createMockCrypto(t *testing.T) Crypto {
	mCrypto := NewMockCrypto(gomock.NewController(t))
	mCrypto.EXPECT().Algorithm().AnyTimes().Return("crypto")
	return mCrypto
}

func createMockGenerator(t *testing.T, crypto Crypto) Generator {
	mGenerator := NewMockGenerator(gomock.NewController(t))
	mGenerator.EXPECT().Alg().AnyTimes().Return(crypto)
	return mGenerator
}

func TestIsCodeExpired(t *testing.T) {
	type args struct {
		creationDate time.Time
		expiry       time.Duration
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"not expired",
			args{
				creationDate: time.Now(),
				expiry:       time.Duration(5 * time.Minute),
			},
			false,
		},
		{
			"never expires",
			args{
				creationDate: time.Now().Add(-5 * time.Minute),
				expiry:       0,
			},
			false,
		},
		{
			"expired",
			args{
				creationDate: time.Now().Add(-5 * time.Minute),
				expiry:       time.Duration(5 * time.Minute),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCodeExpired(tt.args.creationDate, tt.args.expiry); got != tt.want {
				t.Errorf("IsCodeExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyCode(t *testing.T) {
	type args struct {
		creationDate     time.Time
		expiry           time.Duration
		cryptoCode       *CryptoValue
		verificationCode string
		g                Generator
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"expired",
			args{
				creationDate:     time.Now().Add(-5 * time.Minute),
				expiry:           5 * time.Minute,
				cryptoCode:       nil,
				verificationCode: "",
				g:                nil,
			},
			true,
		},
		{
			"unsupported alg err",
			args{
				creationDate:     time.Now(),
				expiry:           5 * time.Minute,
				cryptoCode:       nil,
				verificationCode: "code",
				g:                createMockGenerator(t, createMockCrypto(t)),
			},
			true,
		},
		{
			"encryption alg ok",
			args{
				creationDate: time.Now(),
				expiry:       5 * time.Minute,
				cryptoCode: &CryptoValue{
					CryptoType: TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				},
				verificationCode: "code",
				g:                createMockGenerator(t, createMockEncryptionAlg(t)),
			},
			false,
		},
		{
			"hash alg ok",
			args{
				creationDate: time.Now(),
				expiry:       5 * time.Minute,
				cryptoCode: &CryptoValue{
					CryptoType: TypeHash,
					Algorithm:  "hash",
					Crypted:    []byte("code"),
				},
				verificationCode: "code",
				g:                createMockGenerator(t, createMockHashAlg(t)),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyCode(tt.args.creationDate, tt.args.expiry, tt.args.cryptoCode, tt.args.verificationCode, tt.args.g); (err != nil) != tt.wantErr {
				t.Errorf("VerifyCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyEncryptedCode(t *testing.T) {
	type args struct {
		cryptoCode       *CryptoValue
		verificationCode string
		alg              EncryptionAlgorithm
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"nil error",
			args{
				cryptoCode:       nil,
				verificationCode: "",
				alg:              createMockEncryptionAlg(t),
			},
			true,
		},
		{
			"wrong cryptotype error",
			args{
				cryptoCode: &CryptoValue{
					CryptoType: TypeHash,
					Crypted:    nil,
				},
				verificationCode: "",
				alg:              createMockEncryptionAlg(t),
			},
			true,
		},
		{
			"wrong algorithm error",
			args{
				cryptoCode: &CryptoValue{
					CryptoType: TypeEncryption,
					Algorithm:  "enc2",
					Crypted:    nil,
				},
				verificationCode: "",
				alg:              createMockEncryptionAlg(t),
			},
			true,
		},
		{
			"wrong key id error",
			args{
				cryptoCode: &CryptoValue{
					CryptoType: TypeEncryption,
					Algorithm:  "enc",
					Crypted:    nil,
				},
				verificationCode: "wrong",
				alg:              createMockEncryptionAlg(t),
			},
			true,
		},
		{
			"wrong verification code error",
			args{
				cryptoCode: &CryptoValue{
					CryptoType: TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				},
				verificationCode: "wrong",
				alg:              createMockEncryptionAlg(t),
			},
			true,
		},
		{
			"verification code ok",
			args{
				cryptoCode: &CryptoValue{
					CryptoType: TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				},
				verificationCode: "code",
				alg:              createMockEncryptionAlg(t),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyEncryptedCode(tt.args.cryptoCode, tt.args.verificationCode, tt.args.alg); (err != nil) != tt.wantErr {
				t.Errorf("verifyEncryptedCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyHashedCode(t *testing.T) {
	type args struct {
		cryptoCode       *CryptoValue
		verificationCode string
		alg              HashAlgorithm
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{
			"nil error",
			args{
				cryptoCode:       nil,
				verificationCode: "",
				alg:              createMockHashAlg(t),
			},
			true,
		},
		{
			"wrong cryptotype error",
			args{
				cryptoCode: &CryptoValue{
					CryptoType: TypeEncryption,
					Crypted:    nil,
				},
				verificationCode: "",
				alg:              createMockHashAlg(t),
			},
			true,
		},
		{
			"wrong algorithm error",
			args{
				cryptoCode: &CryptoValue{
					CryptoType: TypeHash,
					Algorithm:  "hash2",
					Crypted:    nil,
				},
				verificationCode: "",
				alg:              createMockHashAlg(t),
			},
			true,
		},
		{
			"wrong verification code error",
			args{
				cryptoCode: &CryptoValue{
					CryptoType: TypeHash,
					Algorithm:  "hash",
					Crypted:    []byte("code"),
				},
				verificationCode: "wrong",
				alg:              createMockHashAlg(t),
			},
			true,
		},
		{
			"verification code ok",
			args{
				cryptoCode: &CryptoValue{
					CryptoType: TypeHash,
					Algorithm:  "hash",
					Crypted:    []byte("code"),
				},
				verificationCode: "code",
				alg:              createMockHashAlg(t),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyHashedCode(tt.args.cryptoCode, tt.args.verificationCode, tt.args.alg); (err != nil) != tt.wantErr {
				t.Errorf("verifyHashedCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
