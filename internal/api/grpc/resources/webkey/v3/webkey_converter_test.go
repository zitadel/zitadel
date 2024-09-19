package webkey

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/v2/internal/crypto"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/query"
	object "github.com/zitadel/zitadel/v2/pkg/grpc/object/v3alpha"
	resource_object "github.com/zitadel/zitadel/v2/pkg/grpc/resources/object/v3alpha"
	webkey "github.com/zitadel/zitadel/v2/pkg/grpc/resources/webkey/v3alpha"
)

func Test_createWebKeyRequestToConfig(t *testing.T) {
	type args struct {
		req *webkey.CreateWebKeyRequest
	}
	tests := []struct {
		name string
		args args
		want crypto.WebKeyConfig
	}{
		{
			name: "RSA",
			args: args{&webkey.CreateWebKeyRequest{
				Key: &webkey.WebKey{
					Config: &webkey.WebKey_Rsa{
						Rsa: &webkey.WebKeyRSAConfig{
							Bits:   webkey.WebKeyRSAConfig_RSA_BITS_3072,
							Hasher: webkey.WebKeyRSAConfig_RSA_HASHER_SHA384,
						},
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
			args: args{&webkey.CreateWebKeyRequest{
				Key: &webkey.WebKey{
					Config: &webkey.WebKey_Ecdsa{
						Ecdsa: &webkey.WebKeyECDSAConfig{
							Curve: webkey.WebKeyECDSAConfig_ECDSA_CURVE_P384,
						},
					},
				},
			}},
			want: &crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP384,
			},
		},
		{
			name: "ED25519",
			args: args{&webkey.CreateWebKeyRequest{
				Key: &webkey.WebKey{
					Config: &webkey.WebKey_Ed25519{
						Ed25519: &webkey.WebKeyED25519Config{},
					},
				},
			}},
			want: &crypto.WebKeyED25519Config{},
		},
		{
			name: "default",
			args: args{&webkey.CreateWebKeyRequest{}},
			want: &crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits2048,
				Hasher: crypto.RSAHasherSHA256,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createWebKeyRequestToConfig(tt.args.req)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_webKeyRSAConfigToCrypto(t *testing.T) {
	type args struct {
		config *webkey.WebKeyRSAConfig
	}
	tests := []struct {
		name string
		args args
		want *crypto.WebKeyRSAConfig
	}{
		{
			name: "unspecified",
			args: args{&webkey.WebKeyRSAConfig{
				Bits:   webkey.WebKeyRSAConfig_RSA_BITS_UNSPECIFIED,
				Hasher: webkey.WebKeyRSAConfig_RSA_HASHER_UNSPECIFIED,
			}},
			want: &crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits2048,
				Hasher: crypto.RSAHasherSHA256,
			},
		},
		{
			name: "2048, RSA256",
			args: args{&webkey.WebKeyRSAConfig{
				Bits:   webkey.WebKeyRSAConfig_RSA_BITS_2048,
				Hasher: webkey.WebKeyRSAConfig_RSA_HASHER_SHA256,
			}},
			want: &crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits2048,
				Hasher: crypto.RSAHasherSHA256,
			},
		},
		{
			name: "3072, RSA384",
			args: args{&webkey.WebKeyRSAConfig{
				Bits:   webkey.WebKeyRSAConfig_RSA_BITS_3072,
				Hasher: webkey.WebKeyRSAConfig_RSA_HASHER_SHA384,
			}},
			want: &crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits3072,
				Hasher: crypto.RSAHasherSHA384,
			},
		},
		{
			name: "4096, RSA512",
			args: args{&webkey.WebKeyRSAConfig{
				Bits:   webkey.WebKeyRSAConfig_RSA_BITS_4096,
				Hasher: webkey.WebKeyRSAConfig_RSA_HASHER_SHA512,
			}},
			want: &crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits4096,
				Hasher: crypto.RSAHasherSHA512,
			},
		},
		{
			name: "invalid",
			args: args{&webkey.WebKeyRSAConfig{
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
		config *webkey.WebKeyECDSAConfig
	}
	tests := []struct {
		name string
		args args
		want *crypto.WebKeyECDSAConfig
	}{
		{
			name: "unspecified",
			args: args{&webkey.WebKeyECDSAConfig{
				Curve: webkey.WebKeyECDSAConfig_ECDSA_CURVE_UNSPECIFIED,
			}},
			want: &crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP256,
			},
		},
		{
			name: "P256",
			args: args{&webkey.WebKeyECDSAConfig{
				Curve: webkey.WebKeyECDSAConfig_ECDSA_CURVE_P256,
			}},
			want: &crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP256,
			},
		},
		{
			name: "P384",
			args: args{&webkey.WebKeyECDSAConfig{
				Curve: webkey.WebKeyECDSAConfig_ECDSA_CURVE_P384,
			}},
			want: &crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP384,
			},
		},
		{
			name: "P512",
			args: args{&webkey.WebKeyECDSAConfig{
				Curve: webkey.WebKeyECDSAConfig_ECDSA_CURVE_P512,
			}},
			want: &crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP512,
			},
		},
		{
			name: "invalid",
			args: args{&webkey.WebKeyECDSAConfig{
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
	instanceID := "ownerid"
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
	want := []*webkey.GetWebKey{
		{
			Details: &resource_object.Details{
				Id:      "key1",
				Created: &timestamppb.Timestamp{Seconds: 123, Nanos: 456},
				Changed: &timestamppb.Timestamp{Seconds: 789, Nanos: 0},
				Owner:   &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instanceID},
			},
			State: webkey.WebKeyState_STATE_ACTIVE,
			Config: &webkey.WebKey{
				Config: &webkey.WebKey_Rsa{
					Rsa: &webkey.WebKeyRSAConfig{
						Bits:   webkey.WebKeyRSAConfig_RSA_BITS_3072,
						Hasher: webkey.WebKeyRSAConfig_RSA_HASHER_SHA384,
					},
				},
			},
		},
		{
			Details: &resource_object.Details{
				Id:      "key2",
				Created: &timestamppb.Timestamp{Seconds: 123, Nanos: 456},
				Changed: &timestamppb.Timestamp{Seconds: 789, Nanos: 0},
				Owner:   &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instanceID},
			},
			State: webkey.WebKeyState_STATE_ACTIVE,
			Config: &webkey.WebKey{
				Config: &webkey.WebKey_Ed25519{
					Ed25519: &webkey.WebKeyED25519Config{},
				},
			},
		},
	}
	got := webKeyDetailsListToPb(list, instanceID)
	assert.Equal(t, want, got)
}

func Test_webKeyDetailsToPb(t *testing.T) {
	instanceID := "ownerid"
	type args struct {
		details *query.WebKeyDetails
	}
	tests := []struct {
		name string
		args args
		want *webkey.GetWebKey
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
			want: &webkey.GetWebKey{
				Details: &resource_object.Details{
					Id:      "keyID",
					Created: &timestamppb.Timestamp{Seconds: 123, Nanos: 456},
					Changed: &timestamppb.Timestamp{Seconds: 789, Nanos: 0},
					Owner:   &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instanceID},
				},
				State: webkey.WebKeyState_STATE_ACTIVE,
				Config: &webkey.WebKey{
					Config: &webkey.WebKey_Rsa{
						Rsa: &webkey.WebKeyRSAConfig{
							Bits:   webkey.WebKeyRSAConfig_RSA_BITS_3072,
							Hasher: webkey.WebKeyRSAConfig_RSA_HASHER_SHA384,
						},
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
			want: &webkey.GetWebKey{
				Details: &resource_object.Details{
					Id:      "keyID",
					Created: &timestamppb.Timestamp{Seconds: 123, Nanos: 456},
					Changed: &timestamppb.Timestamp{Seconds: 789, Nanos: 0},
					Owner:   &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instanceID},
				},
				State: webkey.WebKeyState_STATE_ACTIVE,
				Config: &webkey.WebKey{
					Config: &webkey.WebKey_Ecdsa{
						Ecdsa: &webkey.WebKeyECDSAConfig{
							Curve: webkey.WebKeyECDSAConfig_ECDSA_CURVE_P384,
						},
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
			want: &webkey.GetWebKey{
				Details: &resource_object.Details{
					Id:      "keyID",
					Created: &timestamppb.Timestamp{Seconds: 123, Nanos: 456},
					Changed: &timestamppb.Timestamp{Seconds: 789, Nanos: 0},
					Owner:   &object.Owner{Type: object.OwnerType_OWNER_TYPE_INSTANCE, Id: instanceID},
				},
				State: webkey.WebKeyState_STATE_ACTIVE,
				Config: &webkey.WebKey{
					Config: &webkey.WebKey_Ed25519{
						Ed25519: &webkey.WebKeyED25519Config{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := webKeyDetailsToPb(tt.args.details, instanceID)
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
		want webkey.WebKeyState
	}{
		{
			name: "unspecified",
			args: args{domain.WebKeyStateUnspecified},
			want: webkey.WebKeyState_STATE_UNSPECIFIED,
		},
		{
			name: "initial",
			args: args{domain.WebKeyStateInitial},
			want: webkey.WebKeyState_STATE_INITIAL,
		},
		{
			name: "active",
			args: args{domain.WebKeyStateActive},
			want: webkey.WebKeyState_STATE_ACTIVE,
		},
		{
			name: "inactive",
			args: args{domain.WebKeyStateInactive},
			want: webkey.WebKeyState_STATE_INACTIVE,
		},
		{
			name: "removed",
			args: args{domain.WebKeyStateRemoved},
			want: webkey.WebKeyState_STATE_REMOVED,
		},
		{
			name: "invalid",
			args: args{99},
			want: webkey.WebKeyState_STATE_UNSPECIFIED,
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
		want *webkey.WebKeyRSAConfig
	}{
		{
			name: "2048, RSA256",
			args: args{&crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits2048,
				Hasher: crypto.RSAHasherSHA256,
			}},
			want: &webkey.WebKeyRSAConfig{
				Bits:   webkey.WebKeyRSAConfig_RSA_BITS_2048,
				Hasher: webkey.WebKeyRSAConfig_RSA_HASHER_SHA256,
			},
		},
		{
			name: "3072, RSA384",
			args: args{&crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits3072,
				Hasher: crypto.RSAHasherSHA384,
			}},
			want: &webkey.WebKeyRSAConfig{
				Bits:   webkey.WebKeyRSAConfig_RSA_BITS_3072,
				Hasher: webkey.WebKeyRSAConfig_RSA_HASHER_SHA384,
			},
		},
		{
			name: "4096, RSA512",
			args: args{&crypto.WebKeyRSAConfig{
				Bits:   crypto.RSABits4096,
				Hasher: crypto.RSAHasherSHA512,
			}},
			want: &webkey.WebKeyRSAConfig{
				Bits:   webkey.WebKeyRSAConfig_RSA_BITS_4096,
				Hasher: webkey.WebKeyRSAConfig_RSA_HASHER_SHA512,
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
		want *webkey.WebKeyECDSAConfig
	}{
		{
			name: "P256",
			args: args{&crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP256,
			}},
			want: &webkey.WebKeyECDSAConfig{
				Curve: webkey.WebKeyECDSAConfig_ECDSA_CURVE_P256,
			},
		},
		{
			name: "P384",
			args: args{&crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP384,
			}},
			want: &webkey.WebKeyECDSAConfig{
				Curve: webkey.WebKeyECDSAConfig_ECDSA_CURVE_P384,
			},
		},
		{
			name: "P512",
			args: args{&crypto.WebKeyECDSAConfig{
				Curve: crypto.EllipticCurveP512,
			}},
			want: &webkey.WebKeyECDSAConfig{
				Curve: webkey.WebKeyECDSAConfig_ECDSA_CURVE_P512,
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
