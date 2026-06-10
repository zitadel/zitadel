package crypto

import (
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

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
				g:                createMockGenerator(t, createMockCrypto(t)),
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
				g:                createMockGenerator(t, CreateMockEncryptionAlg(gomock.NewController(t))),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyCode(tt.args.creationDate, tt.args.expiry, tt.args.cryptoCode, tt.args.verificationCode, tt.args.g.Alg()); (err != nil) != tt.wantErr {
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
				alg:              CreateMockEncryptionAlg(gomock.NewController(t)),
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
				alg:              CreateMockEncryptionAlg(gomock.NewController(t)),
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
				alg:              CreateMockEncryptionAlg(gomock.NewController(t)),
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
				alg:              CreateMockEncryptionAlg(gomock.NewController(t)),
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
				alg:              CreateMockEncryptionAlg(gomock.NewController(t)),
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
				alg:              CreateMockEncryptionAlg(gomock.NewController(t)),
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

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name    string
		length  uint
		chars   []rune
		want    string
		wantErr bool
	}{
		{
			name:    "empty chars returns error",
			length:  5,
			chars:   nil,
			wantErr: true,
		},
		{
			name:   "zero length returns empty string",
			length: 0,
			chars:  []rune{'a', 'b'},
			want:   "",
		},
		{
			name:   "single charset char is used for all positions",
			length: 4,
			chars:  []rune{'x'},
			want:   "xxxx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRandomString(tt.length, tt.chars)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateRandomString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("GenerateRandomString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGenerateRandomString_UsesOnlyProvidedChars(t *testing.T) {
	chars := []rune("abcXYZ09")
	got, err := GenerateRandomString(32, chars)
	if err != nil {
		t.Fatalf("GenerateRandomString() error = %v", err)
	}
	if len([]rune(got)) != 32 {
		t.Fatalf("GenerateRandomString() length = %d, want %d", len([]rune(got)), 32)
	}
	allowed := make(map[rune]struct{}, len(chars))
	for _, r := range chars {
		allowed[r] = struct{}{}
	}
	for _, r := range got {
		if _, ok := allowed[r]; !ok {
			t.Fatalf("GenerateRandomString() produced rune %q not in charset", r)
		}
	}
}
