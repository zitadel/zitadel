package authn

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/authn"
)

func KeysToPb(keys []*query.AuthNKey) []*authn.Key {
	k := make([]*authn.Key, len(keys))
	for i, key := range keys {
		k[i] = KeyToPb(key)
	}
	return k
}

func KeyToPb(key *query.AuthNKey) *authn.Key {
	return &authn.Key{
		Id:             key.ID,
		Type:           KeyTypeToPb(key.Type),
		ExpirationDate: timestamppb.New(key.Expiration),
		Details: object.ToViewDetailsPb(
			key.Sequence,
			key.CreationDate,
			key.CreationDate,
			key.ResourceOwner,
		),
	}
}

func KeyTypeToPb(typ domain.AuthNKeyType) authn.KeyType {
	switch typ {
	case domain.AuthNKeyTypeJSON:
		return authn.KeyType_KEY_TYPE_JSON
	default:
		return authn.KeyType_KEY_TYPE_UNSPECIFIED
	}
}

func KeyTypeToDomain(typ authn.KeyType) domain.AuthNKeyType {
	switch typ {
	case authn.KeyType_KEY_TYPE_JSON:
		return domain.AuthNKeyTypeJSON
	default:
		return domain.AuthNKeyTypeNONE
	}
}
