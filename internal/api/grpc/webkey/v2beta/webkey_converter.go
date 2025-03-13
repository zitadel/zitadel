package webkey

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	webkey "github.com/zitadel/zitadel/pkg/grpc/webkey/v2beta"
)

func createWebKeyRequestToConfig(req *webkey.CreateWebKeyRequest) crypto.WebKeyConfig {
	switch config := req.GetKey().(type) {
	case *webkey.CreateWebKeyRequest_Rsa:
		return rsaToCrypto(config.Rsa)
	case *webkey.CreateWebKeyRequest_Ecdsa:
		return ecdsaToCrypto(config.Ecdsa)
	case *webkey.CreateWebKeyRequest_Ed25519:
		return new(crypto.WebKeyED25519Config)
	default:
		return rsaToCrypto(nil)
	}
}

func rsaToCrypto(config *webkey.RSA) *crypto.WebKeyRSAConfig {
	out := new(crypto.WebKeyRSAConfig)

	switch config.GetBits() {
	case webkey.RSABits_RSA_BITS_UNSPECIFIED:
		out.Bits = crypto.RSABits2048
	case webkey.RSABits_RSA_BITS_2048:
		out.Bits = crypto.RSABits2048
	case webkey.RSABits_RSA_BITS_3072:
		out.Bits = crypto.RSABits3072
	case webkey.RSABits_RSA_BITS_4096:
		out.Bits = crypto.RSABits4096
	default:
		out.Bits = crypto.RSABits2048
	}

	switch config.GetHasher() {
	case webkey.RSAHasher_RSA_HASHER_UNSPECIFIED:
		out.Hasher = crypto.RSAHasherSHA256
	case webkey.RSAHasher_RSA_HASHER_SHA256:
		out.Hasher = crypto.RSAHasherSHA256
	case webkey.RSAHasher_RSA_HASHER_SHA384:
		out.Hasher = crypto.RSAHasherSHA384
	case webkey.RSAHasher_RSA_HASHER_SHA512:
		out.Hasher = crypto.RSAHasherSHA512
	default:
		out.Hasher = crypto.RSAHasherSHA256
	}

	return out
}

func ecdsaToCrypto(config *webkey.ECDSA) *crypto.WebKeyECDSAConfig {
	out := new(crypto.WebKeyECDSAConfig)

	switch config.GetCurve() {
	case webkey.ECDSACurve_ECDSA_CURVE_UNSPECIFIED:
		out.Curve = crypto.EllipticCurveP256
	case webkey.ECDSACurve_ECDSA_CURVE_P256:
		out.Curve = crypto.EllipticCurveP256
	case webkey.ECDSACurve_ECDSA_CURVE_P384:
		out.Curve = crypto.EllipticCurveP384
	case webkey.ECDSACurve_ECDSA_CURVE_P512:
		out.Curve = crypto.EllipticCurveP512
	default:
		out.Curve = crypto.EllipticCurveP256
	}

	return out
}

func webKeyDetailsListToPb(list []query.WebKeyDetails) []*webkey.WebKey {
	out := make([]*webkey.WebKey, len(list))
	for i := range list {
		out[i] = webKeyDetailsToPb(&list[i])
	}
	return out
}

func webKeyDetailsToPb(details *query.WebKeyDetails) *webkey.WebKey {
	out := &webkey.WebKey{
		Id:           details.KeyID,
		CreationDate: timestamppb.New(details.CreationDate),
		ChangeDate:   timestamppb.New(details.ChangeDate),
		State:        webKeyStateToPb(details.State),
	}

	switch config := details.Config.(type) {
	case *crypto.WebKeyRSAConfig:
		out.Key = &webkey.WebKey_Rsa{
			Rsa: webKeyRSAConfigToPb(config),
		}
	case *crypto.WebKeyECDSAConfig:
		out.Key = &webkey.WebKey_Ecdsa{
			Ecdsa: webKeyECDSAConfigToPb(config),
		}
	case *crypto.WebKeyED25519Config:
		out.Key = &webkey.WebKey_Ed25519{
			Ed25519: new(webkey.ED25519),
		}
	}

	return out
}

func webKeyStateToPb(state domain.WebKeyState) webkey.State {
	switch state {
	case domain.WebKeyStateUnspecified:
		return webkey.State_STATE_UNSPECIFIED
	case domain.WebKeyStateInitial:
		return webkey.State_STATE_INITIAL
	case domain.WebKeyStateActive:
		return webkey.State_STATE_ACTIVE
	case domain.WebKeyStateInactive:
		return webkey.State_STATE_INACTIVE
	case domain.WebKeyStateRemoved:
		return webkey.State_STATE_REMOVED
	default:
		return webkey.State_STATE_UNSPECIFIED
	}
}

func webKeyRSAConfigToPb(config *crypto.WebKeyRSAConfig) *webkey.RSA {
	out := new(webkey.RSA)

	switch config.Bits {
	case crypto.RSABitsUnspecified:
		out.Bits = webkey.RSABits_RSA_BITS_UNSPECIFIED
	case crypto.RSABits2048:
		out.Bits = webkey.RSABits_RSA_BITS_2048
	case crypto.RSABits3072:
		out.Bits = webkey.RSABits_RSA_BITS_3072
	case crypto.RSABits4096:
		out.Bits = webkey.RSABits_RSA_BITS_4096
	}

	switch config.Hasher {
	case crypto.RSAHasherUnspecified:
		out.Hasher = webkey.RSAHasher_RSA_HASHER_UNSPECIFIED
	case crypto.RSAHasherSHA256:
		out.Hasher = webkey.RSAHasher_RSA_HASHER_SHA256
	case crypto.RSAHasherSHA384:
		out.Hasher = webkey.RSAHasher_RSA_HASHER_SHA384
	case crypto.RSAHasherSHA512:
		out.Hasher = webkey.RSAHasher_RSA_HASHER_SHA512
	}

	return out
}

func webKeyECDSAConfigToPb(config *crypto.WebKeyECDSAConfig) *webkey.ECDSA {
	out := new(webkey.ECDSA)

	switch config.Curve {
	case crypto.EllipticCurveUnspecified:
		out.Curve = webkey.ECDSACurve_ECDSA_CURVE_UNSPECIFIED
	case crypto.EllipticCurveP256:
		out.Curve = webkey.ECDSACurve_ECDSA_CURVE_P256
	case crypto.EllipticCurveP384:
		out.Curve = webkey.ECDSACurve_ECDSA_CURVE_P384
	case crypto.EllipticCurveP512:
		out.Curve = webkey.ECDSACurve_ECDSA_CURVE_P512
	}

	return out
}
