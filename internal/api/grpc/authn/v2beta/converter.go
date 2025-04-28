package authn

import (
	"github.com/zitadel/zitadel/internal/domain"
	authn "github.com/zitadel/zitadel/pkg/grpc/authn/v2beta"
)

func KeyTypeToDomain(t authn.KeyType) domain.AuthNKeyType {
	switch t {
	case authn.KeyType_KEY_TYPE_JSON:
		return domain.AuthNKeyTypeJSON
	default:
		return domain.AuthNKeyTypeNONE
	}
}
