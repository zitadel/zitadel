package crypto

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

type mockEncCrypto struct {
}

func (m *mockEncCrypto) Algorithm() string {
	return "enc"
}

func (m *mockEncCrypto) Encrypt(value []byte) ([]byte, error) {
	return value, nil
}

func (m *mockEncCrypto) Decrypt(value []byte, _ string) ([]byte, error) {
	return value, nil
}

func (m *mockEncCrypto) DecryptString(value []byte, _ string) (string, error) {
	return string(value), nil
}

func (m *mockEncCrypto) EncryptionKeyID() string {
	return "keyID"
}
func (m *mockEncCrypto) DecryptionKeyIDs() []string {
	return []string{"keyID"}
}

type mockHashCrypto struct {
}

func (m *mockHashCrypto) Algorithm() string {
	return "hash"
}

func (m *mockHashCrypto) Hash(value []byte) ([]byte, error) {
	return value, nil
}

func (m *mockHashCrypto) CompareHash(hashed, comparer []byte) error {
	if !bytes.Equal(hashed, comparer) {
		return errors.New("not equal")
	}
	return nil
}

type alg struct{}

func (a *alg) Algorithm() string {
	return "alg"
}

func TestCrypt(t *testing.T) {
	type args struct {
		value []byte
		c     EncryptionAlgorithm
	}
	tests := []struct {
		name    string
		args    args
		want    *CryptoValue
		wantErr bool
	}{
		{
			"encrypt",
			args{[]byte("test"), &mockEncCrypto{}},
			&CryptoValue{CryptoType: TypeEncryption, Algorithm: "enc", KeyID: "keyID", Crypted: []byte("test")},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Crypt(tt.args.value, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Crypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Crypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	type args struct {
		value []byte
		c     EncryptionAlgorithm
	}
	tests := []struct {
		name    string
		args    args
		want    *CryptoValue
		wantErr bool
	}{
		{
			"ok",
			args{[]byte("test"), &mockEncCrypto{}},
			&CryptoValue{CryptoType: TypeEncryption, Algorithm: "enc", KeyID: "keyID", Crypted: []byte("test")},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt(tt.args.value, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	type args struct {
		value *CryptoValue
		c     EncryptionAlgorithm
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"ok",
			args{&CryptoValue{CryptoType: TypeEncryption, Algorithm: "enc", KeyID: "keyID", Crypted: []byte("test")}, &mockEncCrypto{}},
			[]byte("test"),
			false,
		},
		{
			"wrong id",
			args{&CryptoValue{CryptoType: TypeEncryption, Algorithm: "enc", KeyID: "keyID2", Crypted: []byte("test")}, &mockEncCrypto{}},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypt(tt.args.value, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecryptString(t *testing.T) {
	type args struct {
		value *CryptoValue
		c     EncryptionAlgorithm
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"ok",
			args{&CryptoValue{CryptoType: TypeEncryption, Algorithm: "enc", KeyID: "keyID", Crypted: []byte("test")}, &mockEncCrypto{}},
			"test",
			false,
		},
		{
			"wrong id",
			args{&CryptoValue{CryptoType: TypeEncryption, Algorithm: "enc", KeyID: "keyID2", Crypted: []byte("test")}, &mockEncCrypto{}},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecryptString(tt.args.value, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecryptString() = %v, want %v", got, tt.want)
			}
		})
	}
}
