package webkey

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/webkey/v3alpha"
)

func Test_generateWebKeyRequestToConfig(t *testing.T) {
	type args struct {
		req *v3alpha.GenerateWebKeyRequest
	}
	tests := []struct {
		name string
		args args
		want crypto.WebKeyConfig
	}{
		{
			name: "RSA",
			args: args{&v3alpha.GenerateWebKeyRequest{
				Config: &v3alpha.GenerateWebKeyRequest_Rsa{
					Rsa: &v3alpha.WebKeyRSAConfig{
						Bits:   v3alpha.WebKeyRSAConfig_RSA_BITS_3072,
						Hasher: v3alpha.WebKeyRSAConfig_RSA_HASHER_SHA384,
					},
				},
			}},
			want: &crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits3072,
				Hasher: crypto.RSAHasherSHA384,
			},
		},
		{
			name: "ECDSA",
			args: args{&v3alpha.GenerateWebKeyRequest{
				Config: &v3alpha.GenerateWebKeyRequest_Ecdsa{
					Ecdsa: &v3alpha.WebKeyECDSAConfig{
						Curve: v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_P384,
					},
				},
			}},
			want: &crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP384,
			},
		},
		{
			name: "ED25519",
			args: args{&v3alpha.GenerateWebKeyRequest{
				Config: &v3alpha.GenerateWebKeyRequest_Ed25519{
					Ed25519: &v3alpha.WebKeyED25519Config{},
				},
			}},
			want: &crypto.WebKeyED25519Config{},
		},
		{
			name: "default",
			args: args{&v3alpha.GenerateWebKeyRequest{
				Config: nil,
			}},
			want: &crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits2048,
				Hasher: crypto.RSAHasherSHA256,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateWebKeyRequestToConfig(tt.args.req)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_webKeyRSAConfigToCrypto(t *testing.T) {
	type args struct {
		config *v3alpha.WebKeyRSAConfig
	}
	tests := []struct {
		name string
		args args
		want *crypto.WebKeyRSAConfig
	}{
		{
			name: "unspecified",
			args: args{&v3alpha.WebKeyRSAConfig{
				Bits:   v3alpha.WebKeyRSAConfig_RSA_BITS_UNSPECIFIED,
				Hasher: v3alpha.WebKeyRSAConfig_RSA_HASHER_UNSPECIFIED,
			}},
			want: &crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits2048,
				Hasher: crypto.RSAHasherSHA256,
			},
		},
		{
			name: "2048, RSA256",
			args: args{&v3alpha.WebKeyRSAConfig{
				Bits:   v3alpha.WebKeyRSAConfig_RSA_BITS_2048,
				Hasher: v3alpha.WebKeyRSAConfig_RSA_HASHER_SHA256,
			}},
			want: &crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits2048,
				Hasher: crypto.RSAHasherSHA256,
			},
		},
		{
			name: "3072, RSA384",
			args: args{&v3alpha.WebKeyRSAConfig{
				Bits:   v3alpha.WebKeyRSAConfig_RSA_BITS_3072,
				Hasher: v3alpha.WebKeyRSAConfig_RSA_HASHER_SHA384,
			}},
			want: &crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits3072,
				Hasher: crypto.RSAHasherSHA384,
			},
		},
		{
			name: "4096, RSA512",
			args: args{&v3alpha.WebKeyRSAConfig{
				Bits:   v3alpha.WebKeyRSAConfig_RSA_BITS_4096,
				Hasher: v3alpha.WebKeyRSAConfig_RSA_HASHER_SHA512,
			}},
			want: &crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits4096,
				Hasher: crypto.RSAHasherSHA512,
			},
		},
		{
			name: "invalid",
			args: args{&v3alpha.WebKeyRSAConfig{
				Bits:   99,
				Hasher: 99,
			}},
			want: &crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits2048,
				Hasher: crypto.RSAHasherSHA256,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := webKeyRSAConfigToCrypto(tt.args.config)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_webKeyECDSAConfigToCrypto(t *testing.T) {
	type args struct {
		config *v3alpha.WebKeyECDSAConfig
	}
	tests := []struct {
		name string
		args args
		want *crypto.WebKeyECDSAConfig
	}{
		{
			name: "unspecified",
			args: args{&v3alpha.WebKeyECDSAConfig{
				Curve: v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_UNSPECIFIED,
			}},
			want: &crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP256,
			},
		},
		{
			name: "P256",
			args: args{&v3alpha.WebKeyECDSAConfig{
				Curve: v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_P256,
			}},
			want: &crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP256,
			},
		},
		{
			name: "P384",
			args: args{&v3alpha.WebKeyECDSAConfig{
				Curve: v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_P384,
			}},
			want: &crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP384,
			},
		},
		{
			name: "P512",
			args: args{&v3alpha.WebKeyECDSAConfig{
				Curve: v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_P512,
			}},
			want: &crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP512,
			},
		},
		{
			name: "invalid",
			args: args{&v3alpha.WebKeyECDSAConfig{
				Curve: 99,
			}},
			want: &crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP256,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := webKeyECDSAConfigToCrypto(tt.args.config)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_webKeyDetailsListToPb(t *testing.T) {
	list := []query.WebKeyDetails{
		{
			KeyID:        "key1",
			CreationDate: time.Unix(123, 456),
			ChangeDate:   time.Unix(789, 0),
			Sequence:     123,
			State:        domain.WebKeyStateActive,
			Config: &crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits3072,
				Hasher: crypto.RSAHasherSHA384,
			},
		},
		{
			KeyID:        "key2",
			CreationDate: time.Unix(123, 456),
			ChangeDate:   time.Unix(789, 0),
			Sequence:     123,
			State:        domain.WebKeyStateActive,
			Config:       &crypto.WebKeyED25519Config{},
		},
	}
	want := []*v3alpha.WebKeyDetails{
		{
			KeyId:       "key1",
			CreatedDate: &timestamppb.Timestamp{Seconds: 123, Nanos: 456},
			ChangeDate:  &timestamppb.Timestamp{Seconds: 789, Nanos: 0},
			Sequence:    123,
			State:       v3alpha.WebKeyState_STATE_ACTIVE,
			Config: &v3alpha.WebKeyDetails_Rsa{
				Rsa: &v3alpha.WebKeyRSAConfig{
					Bits:   v3alpha.WebKeyRSAConfig_RSA_BITS_3072,
					Hasher: v3alpha.WebKeyRSAConfig_RSA_HASHER_SHA384,
				},
			},
		},
		{
			KeyId:       "key2",
			CreatedDate: &timestamppb.Timestamp{Seconds: 123, Nanos: 456},
			ChangeDate:  &timestamppb.Timestamp{Seconds: 789, Nanos: 0},
			Sequence:    123,
			State:       v3alpha.WebKeyState_STATE_ACTIVE,
			Config: &v3alpha.WebKeyDetails_Ed25519{
				Ed25519: &v3alpha.WebKeyED25519Config{},
			},
		},
	}
	got := webKeyDetailsListToPb(list)
	assert.Equal(t, want, got)
}

func Test_webKeyDetailsToPb(t *testing.T) {
	type args struct {
		details *query.WebKeyDetails
	}
	tests := []struct {
		name string
		args args
		want *v3alpha.WebKeyDetails
	}{
		{
			name: "RSA",
			args: args{&query.WebKeyDetails{
				KeyID:        "keyID",
				CreationDate: time.Unix(123, 456),
				ChangeDate:   time.Unix(789, 0),
				Sequence:     123,
				State:        domain.WebKeyStateActive,
				Config: &crypto.WebKeyRSAConfig{
					Bits:   crypto.RSABits3072,
					Hasher: crypto.RSAHasherSHA384,
				},
			}},
			want: &v3alpha.WebKeyDetails{
				KeyId:       "keyID",
				CreatedDate: &timestamppb.Timestamp{Seconds: 123, Nanos: 456},
				ChangeDate:  &timestamppb.Timestamp{Seconds: 789, Nanos: 0},
				Sequence:    123,
				State:       v3alpha.WebKeyState_STATE_ACTIVE,
				Config: &v3alpha.WebKeyDetails_Rsa{
					Rsa: &v3alpha.WebKeyRSAConfig{
						Bits:   v3alpha.WebKeyRSAConfig_RSA_BITS_3072,
						Hasher: v3alpha.WebKeyRSAConfig_RSA_HASHER_SHA384,
					},
				},
			},
		},
		{
			name: "ECDSA",
			args: args{&query.WebKeyDetails{
				KeyID:        "keyID",
				CreationDate: time.Unix(123, 456),
				ChangeDate:   time.Unix(789, 0),
				Sequence:     123,
				State:        domain.WebKeyStateActive,
				Config: &crypto.WebKeyECDSAConfig{
					Curve: crypto.EllipticCurveP384,
				},
			}},
			want: &v3alpha.WebKeyDetails{
				KeyId:       "keyID",
				CreatedDate: &timestamppb.Timestamp{Seconds: 123, Nanos: 456},
				ChangeDate:  &timestamppb.Timestamp{Seconds: 789, Nanos: 0},
				Sequence:    123,
				State:       v3alpha.WebKeyState_STATE_ACTIVE,
				Config: &v3alpha.WebKeyDetails_Ecdsa{
					Ecdsa: &v3alpha.WebKeyECDSAConfig{
						Curve: v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_P384,
					},
				},
			},
		},
		{
			name: "ED25519",
			args: args{&query.WebKeyDetails{
				KeyID:        "keyID",
				CreationDate: time.Unix(123, 456),
				ChangeDate:   time.Unix(789, 0),
				Sequence:     123,
				State:        domain.WebKeyStateActive,
				Config:       &crypto.WebKeyED25519Config{},
			}},
			want: &v3alpha.WebKeyDetails{
				KeyId:       "keyID",
				CreatedDate: &timestamppb.Timestamp{Seconds: 123, Nanos: 456},
				ChangeDate:  &timestamppb.Timestamp{Seconds: 789, Nanos: 0},
				Sequence:    123,
				State:       v3alpha.WebKeyState_STATE_ACTIVE,
				Config: &v3alpha.WebKeyDetails_Ed25519{
					Ed25519: &v3alpha.WebKeyED25519Config{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := webKeyDetailsToPb(tt.args.details)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_webKeyStateToPb(t *testing.T) {
	type args struct {
		state domain.WebKeyState
	}
	tests := []struct {
		name string
		args args
		want v3alpha.WebKeyState
	}{
		{
			name: "unspecified",
			args: args{domain.WebKeyStateUnspecified},
			want: v3alpha.WebKeyState_STATE_UNSPECIFIED,
		},
		{
			name: "inactive",
			args: args{domain.WebKeyStateInactive},
			want: v3alpha.WebKeyState_STATE_INACTIVE,
		},
		{
			name: "active",
			args: args{domain.WebKeyStateActive},
			want: v3alpha.WebKeyState_STATE_ACTIVE,
		},
		{
			name: "removed",
			args: args{domain.WebKeyStateRemoved},
			want: v3alpha.WebKeyState_STATE_REMOVED,
		},
		{
			name: "invalid",
			args: args{99},
			want: v3alpha.WebKeyState_STATE_UNSPECIFIED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := webKeyStateToPb(tt.args.state)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_webKeyRSAConfigToPb(t *testing.T) {
	type args struct {
		config *crypto.WebKeyRSAConfig
	}
	tests := []struct {
		name string
		args args
		want *v3alpha.WebKeyRSAConfig
	}{
		{
			name: "2048, RSA256",
			args: args{&crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits2048,
				Hasher: crypto.RSAHasherSHA256,
			}},
			want: &v3alpha.WebKeyRSAConfig{
				Bits:   v3alpha.WebKeyRSAConfig_RSA_BITS_2048,
				Hasher: v3alpha.WebKeyRSAConfig_RSA_HASHER_SHA256,
			},
		},
		{
			name: "3072, RSA384",
			args: args{&crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits3072,
				Hasher: crypto.RSAHasherSHA384,
			}},
			want: &v3alpha.WebKeyRSAConfig{
				Bits:   v3alpha.WebKeyRSAConfig_RSA_BITS_3072,
				Hasher: v3alpha.WebKeyRSAConfig_RSA_HASHER_SHA384,
			},
		},
		{
			name: "4096, RSA512",
			args: args{&crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits4096,
				Hasher: crypto.RSAHasherSHA512,
			}},
			want: &v3alpha.WebKeyRSAConfig{
				Bits:   v3alpha.WebKeyRSAConfig_RSA_BITS_4096,
				Hasher: v3alpha.WebKeyRSAConfig_RSA_HASHER_SHA512,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := webKeyRSAConfigToPb(tt.args.config)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_webKeyECDSAConfigToPb(t *testing.T) {
	type args struct {
		config *crypto.WebKeyECDSAConfig
	}
	tests := []struct {
		name string
		args args
		want *v3alpha.WebKeyECDSAConfig
	}{
		{
			name: "P256",
			args: args{&crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP256,
			}},
			want: &v3alpha.WebKeyECDSAConfig{
				Curve: v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_P256,
			},
		},
		{
			name: "P384",
			args: args{&crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP384,
			}},
			want: &v3alpha.WebKeyECDSAConfig{
				Curve: v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_P384,
			},
		},
		{
			name: "P512",
			args: args{&crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP512,
			}},
			want: &v3alpha.WebKeyECDSAConfig{
				Curve: v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_P512,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := webKeyECDSAConfigToPb(tt.args.config)
			assert.Equal(t, tt.want, got)
		})
	}
}
