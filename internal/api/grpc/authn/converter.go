package authn

import (
	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/pkg/grpc/authn"
)

func KeyViewsToPb(keys []*key_model.AuthNKeyView) []*authn.Key {
	k := make([]*authn.Key, len(keys))
	for i, key := range keys {
		k[i] = KeyViewToPb(key)
	}
	return k
}

func KeyViewToPb(key *key_model.AuthNKeyView) *authn.Key {
	expDate, err := ptypes.TimestampProto(key.ExpirationDate)
	logging.Log("AUTHN-uhYmM").OnError(err).Debug("unable to parse expiry")

	return &authn.Key{
		Id:             key.ID,
		Type:           authn.KeyType_KEY_TYPE_JSON,
		ExpirationDate: expDate,
		Details: object.ToViewDetailsPb(
			key.Sequence,
			key.CreationDate,
			key.CreationDate,
			"", //TODO: details
		),
	}
}

func KeyToPb(key *key_model.AuthNKeyView) *authn.Key {
	expDate, err := ptypes.TimestampProto(key.ExpirationDate)
	logging.Log("AUTHN-4n12g").OnError(err).Debug("unable to parse expiration date")

	return &authn.Key{
		Id:             key.ID,
		Type:           KeyTypeToPb(key.Type),
		ExpirationDate: expDate,
		Details: object.ToViewDetailsPb(
			key.Sequence,
			key.CreationDate,
			key.CreationDate,
			"", //TODO: details
		),
	}
}

func KeyTypeToPb(typ key_model.AuthNKeyType) authn.KeyType {
	switch typ {
	case key_model.AuthNKeyTypeJSON:
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
