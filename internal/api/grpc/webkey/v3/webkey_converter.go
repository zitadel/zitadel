package webkey

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/webkey/v3alpha"
)

func generateWebKeyRequestToConfig(req *v3alpha.GenerateWebKeyRequest) crypto.WebKeyConfig {
	switch config := req.GetConfig().(type) {
	case *v3alpha.GenerateWebKeyRequest_Rsa:
		return webKeyRSAConfigToCrypto(config.Rsa)
	case *v3alpha.GenerateWebKeyRequest_Ecdsa:
		return webKeyECDSAConfigToCrypto(config.Ecdsa)
	case *v3alpha.GenerateWebKeyRequest_Ed25519:
		return new(crypto.WebKeyED25519Config)
	default:
		return webKeyRSAConfigToCrypto(nil)
	}
}

func webKeyRSAConfigToCrypto(config *v3alpha.WebKeyRSAConfig) *crypto.WebKeyRSAConfig {
	out := new(crypto.WebKeyRSAConfig)

	switch config.GetBits() {
	case v3alpha.WebKeyRSAConfig_RSA_BITS_UNSPECIFIED:
		out.Bits = crypto.RSABits2048
	case v3alpha.WebKeyRSAConfig_RSA_BITS_2048:
		out.Bits = crypto.RSABits2048
	case v3alpha.WebKeyRSAConfig_RSA_BITS_3072:
		out.Bits = crypto.RSABits3072
	case v3alpha.WebKeyRSAConfig_RSA_BITS_4096:
		out.Bits = crypto.RSABits4096
	default:
		out.Bits = crypto.RSABits2048
	}

	switch config.GetHasher() {
	case v3alpha.WebKeyRSAConfig_RSA_HASHER_UNSPECIFIED:
		out.Hasher = crypto.RSAHasherSHA256
	case v3alpha.WebKeyRSAConfig_RSA_HASHER_SHA256:
		out.Hasher = crypto.RSAHasherSHA256
	case v3alpha.WebKeyRSAConfig_RSA_HASHER_SHA384:
		out.Hasher = crypto.RSAHasherSHA384
	case v3alpha.WebKeyRSAConfig_RSA_HASHER_SHA512:
		out.Hasher = crypto.RSAHasherSHA512
	default:
		out.Hasher = crypto.RSAHasherSHA256
	}

	return out
}

func webKeyECDSAConfigToCrypto(config *v3alpha.WebKeyECDSAConfig) *crypto.WebKeyECDSAConfig {
	out := new(crypto.WebKeyECDSAConfig)

	switch config.GetCurve() {
	case v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_UNSPECIFIED:
		out.Curve = crypto.EllipticCurveP256
	case v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_P256:
		out.Curve = crypto.EllipticCurveP256
	case v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_P384:
		out.Curve = crypto.EllipticCurveP384
	case v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_P512:
		out.Curve = crypto.EllipticCurveP512
	default:
		out.Curve = crypto.EllipticCurveP256
	}

	return out
}

func webKeyDetailsListToPb(list []query.WebKeyDetails) []*v3alpha.WebKeyDetails {
	out := make([]*v3alpha.WebKeyDetails, len(list))
	for i := range list {
		out[i] = webKeyDetailsToPb(&list[i])
	}
	return out
}

func webKeyDetailsToPb(details *query.WebKeyDetails) *v3alpha.WebKeyDetails {
	out := &v3alpha.WebKeyDetails{
		KeyId:       details.KeyID,
		CreatedDate: timestamppb.New(details.CreationDate),
		ChangeDate:  timestamppb.New(details.ChangeDate),
		Sequence:    details.Sequence,
		State:       webKeyStateToPb(details.State),
	}

	switch config := details.Config.(type) {
	case *crypto.WebKeyRSAConfig:
		out.Config = &v3alpha.WebKeyDetails_Rsa{
			Rsa: webKeyRSAConfigToPb(config),
		}
	case *crypto.WebKeyECDSAConfig:
		out.Config = &v3alpha.WebKeyDetails_Ecdsa{
			Ecdsa: webKeyECDSAConfigToPb(config),
		}
	case *crypto.WebKeyED25519Config:
		out.Config = &v3alpha.WebKeyDetails_Ed25519{
			Ed25519: new(v3alpha.WebKeyED25519Config),
		}
	}

	return out
}

func webKeyStateToPb(state domain.WebKeyState) v3alpha.WebKeyState {
	switch state {
	case domain.WebKeyStateUnspecified:
		return v3alpha.WebKeyState_STATE_UNSPECIFIED
	case domain.WebKeyStateInactive:
		return v3alpha.WebKeyState_STATE_INACTIVE
	case domain.WebKeyStateActive:
		return v3alpha.WebKeyState_STATE_ACTIVE
	case domain.WebKeyStateRemoved:
		return v3alpha.WebKeyState_STATE_REMOVED
	default:
		return v3alpha.WebKeyState_STATE_UNSPECIFIED
	}
}

func webKeyRSAConfigToPb(config *crypto.WebKeyRSAConfig) *v3alpha.WebKeyRSAConfig {
	out := new(v3alpha.WebKeyRSAConfig)

	switch config.Bits {
	case crypto.RSABits2048:
		out.Bits = v3alpha.WebKeyRSAConfig_RSA_BITS_2048
	case crypto.RSABits3072:
		out.Bits = v3alpha.WebKeyRSAConfig_RSA_BITS_3072
	case crypto.RSABits4096:
		out.Bits = v3alpha.WebKeyRSAConfig_RSA_BITS_4096
	}

	switch config.Hasher {
	case crypto.RSAHasherSHA256:
		out.Hasher = v3alpha.WebKeyRSAConfig_RSA_HASHER_SHA256
	case crypto.RSAHasherSHA384:
		out.Hasher = v3alpha.WebKeyRSAConfig_RSA_HASHER_SHA384
	case crypto.RSAHasherSHA512:
		out.Hasher = v3alpha.WebKeyRSAConfig_RSA_HASHER_SHA512
	}

	return out
}

func webKeyECDSAConfigToPb(config *crypto.WebKeyECDSAConfig) *v3alpha.WebKeyECDSAConfig {
	out := new(v3alpha.WebKeyECDSAConfig)

	switch config.Curve {
	case crypto.EllipticCurveP256:
		out.Curve = v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_P256
	case crypto.EllipticCurveP384:
		out.Curve = v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_P384
	case crypto.EllipticCurveP512:
		out.Curve = v3alpha.WebKeyECDSAConfig_ECDSA_CURVE_P512
	}

	return out
}
