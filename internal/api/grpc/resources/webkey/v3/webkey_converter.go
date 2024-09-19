package webkey

import (
	resource_object "github.com/zitadel/zitadel/v2/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/v2/internal/crypto"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/query"
	object "github.com/zitadel/zitadel/v2/pkg/grpc/object/v3alpha"
	webkey "github.com/zitadel/zitadel/v2/pkg/grpc/resources/webkey/v3alpha"
)

func createWebKeyRequestToConfig(req *webkey.CreateWebKeyRequest) crypto.WebKeyConfig {
	switch config := req.GetKey().GetConfig().(type) {
	case *webkey.WebKey_Rsa:
		return webKeyRSAConfigToCrypto(config.Rsa)
	case *webkey.WebKey_Ecdsa:
		return webKeyECDSAConfigToCrypto(config.Ecdsa)
	case *webkey.WebKey_Ed25519:
		return new(crypto.WebKeyED25519Config)
	default:
		return webKeyRSAConfigToCrypto(nil)
	}
}

func webKeyRSAConfigToCrypto(config *webkey.WebKeyRSAConfig) *crypto.WebKeyRSAConfig {
	out := new(crypto.WebKeyRSAConfig)

	switch config.GetBits() {
	case webkey.WebKeyRSAConfig_RSA_BITS_UNSPECIFIED:
		out.Bits = crypto.RSABits2048
	case webkey.WebKeyRSAConfig_RSA_BITS_2048:
		out.Bits = crypto.RSABits2048
	case webkey.WebKeyRSAConfig_RSA_BITS_3072:
		out.Bits = crypto.RSABits3072
	case webkey.WebKeyRSAConfig_RSA_BITS_4096:
		out.Bits = crypto.RSABits4096
	default:
		out.Bits = crypto.RSABits2048
	}

	switch config.GetHasher() {
	case webkey.WebKeyRSAConfig_RSA_HASHER_UNSPECIFIED:
		out.Hasher = crypto.RSAHasherSHA256
	case webkey.WebKeyRSAConfig_RSA_HASHER_SHA256:
		out.Hasher = crypto.RSAHasherSHA256
	case webkey.WebKeyRSAConfig_RSA_HASHER_SHA384:
		out.Hasher = crypto.RSAHasherSHA384
	case webkey.WebKeyRSAConfig_RSA_HASHER_SHA512:
		out.Hasher = crypto.RSAHasherSHA512
	default:
		out.Hasher = crypto.RSAHasherSHA256
	}

	return out
}

func webKeyECDSAConfigToCrypto(config *webkey.WebKeyECDSAConfig) *crypto.WebKeyECDSAConfig {
	out := new(crypto.WebKeyECDSAConfig)

	switch config.GetCurve() {
	case webkey.WebKeyECDSAConfig_ECDSA_CURVE_UNSPECIFIED:
		out.Curve = crypto.EllipticCurveP256
	case webkey.WebKeyECDSAConfig_ECDSA_CURVE_P256:
		out.Curve = crypto.EllipticCurveP256
	case webkey.WebKeyECDSAConfig_ECDSA_CURVE_P384:
		out.Curve = crypto.EllipticCurveP384
	case webkey.WebKeyECDSAConfig_ECDSA_CURVE_P512:
		out.Curve = crypto.EllipticCurveP512
	default:
		out.Curve = crypto.EllipticCurveP256
	}

	return out
}

func webKeyDetailsListToPb(list []query.WebKeyDetails, instanceID string) []*webkey.GetWebKey {
	out := make([]*webkey.GetWebKey, len(list))
	for i := range list {
		out[i] = webKeyDetailsToPb(&list[i], instanceID)
	}
	return out
}

func webKeyDetailsToPb(details *query.WebKeyDetails, instanceID string) *webkey.GetWebKey {
	out := &webkey.GetWebKey{
		Details: resource_object.DomainToDetailsPb(&domain.ObjectDetails{
			ID:           details.KeyID,
			CreationDate: details.CreationDate,
			EventDate:    details.ChangeDate,
		}, object.OwnerType_OWNER_TYPE_INSTANCE, instanceID),
		State:  webKeyStateToPb(details.State),
		Config: &webkey.WebKey{},
	}

	switch config := details.Config.(type) {
	case *crypto.WebKeyRSAConfig:
		out.Config.Config = &webkey.WebKey_Rsa{
			Rsa: webKeyRSAConfigToPb(config),
		}
	case *crypto.WebKeyECDSAConfig:
		out.Config.Config = &webkey.WebKey_Ecdsa{
			Ecdsa: webKeyECDSAConfigToPb(config),
		}
	case *crypto.WebKeyED25519Config:
		out.Config.Config = &webkey.WebKey_Ed25519{
			Ed25519: new(webkey.WebKeyED25519Config),
		}
	}

	return out
}

func webKeyStateToPb(state domain.WebKeyState) webkey.WebKeyState {
	switch state {
	case domain.WebKeyStateUnspecified:
		return webkey.WebKeyState_STATE_UNSPECIFIED
	case domain.WebKeyStateInitial:
		return webkey.WebKeyState_STATE_INITIAL
	case domain.WebKeyStateActive:
		return webkey.WebKeyState_STATE_ACTIVE
	case domain.WebKeyStateInactive:
		return webkey.WebKeyState_STATE_INACTIVE
	case domain.WebKeyStateRemoved:
		return webkey.WebKeyState_STATE_REMOVED
	default:
		return webkey.WebKeyState_STATE_UNSPECIFIED
	}
}

func webKeyRSAConfigToPb(config *crypto.WebKeyRSAConfig) *webkey.WebKeyRSAConfig {
	out := new(webkey.WebKeyRSAConfig)

	switch config.Bits {
	case crypto.RSABitsUnspecified:
		out.Bits = webkey.WebKeyRSAConfig_RSA_BITS_UNSPECIFIED
	case crypto.RSABits2048:
		out.Bits = webkey.WebKeyRSAConfig_RSA_BITS_2048
	case crypto.RSABits3072:
		out.Bits = webkey.WebKeyRSAConfig_RSA_BITS_3072
	case crypto.RSABits4096:
		out.Bits = webkey.WebKeyRSAConfig_RSA_BITS_4096
	}

	switch config.Hasher {
	case crypto.RSAHasherUnspecified:
		out.Hasher = webkey.WebKeyRSAConfig_RSA_HASHER_UNSPECIFIED
	case crypto.RSAHasherSHA256:
		out.Hasher = webkey.WebKeyRSAConfig_RSA_HASHER_SHA256
	case crypto.RSAHasherSHA384:
		out.Hasher = webkey.WebKeyRSAConfig_RSA_HASHER_SHA384
	case crypto.RSAHasherSHA512:
		out.Hasher = webkey.WebKeyRSAConfig_RSA_HASHER_SHA512
	}

	return out
}

func webKeyECDSAConfigToPb(config *crypto.WebKeyECDSAConfig) *webkey.WebKeyECDSAConfig {
	out := new(webkey.WebKeyECDSAConfig)

	switch config.Curve {
	case crypto.EllipticCurveUnspecified:
		out.Curve = webkey.WebKeyECDSAConfig_ECDSA_CURVE_UNSPECIFIED
	case crypto.EllipticCurveP256:
		out.Curve = webkey.WebKeyECDSAConfig_ECDSA_CURVE_P256
	case crypto.EllipticCurveP384:
		out.Curve = webkey.WebKeyECDSAConfig_ECDSA_CURVE_P384
	case crypto.EllipticCurveP512:
		out.Curve = webkey.WebKeyECDSAConfig_ECDSA_CURVE_P512
	}

	return out
}
