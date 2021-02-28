package authn

import (
	"encoding/json"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/pkg/grpc/authn"
	"github.com/golang/protobuf/ptypes"
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
		Details: object.ToDetailsPb(
			key.Sequence,
			key.CreationDate,
			key.CreationDate,    //TODO: details
			"key.ResourceOwner", //TODO: details
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
		Details: object.ToDetailsPb(
			key.Sequence,
			key.CreationDate,
			key.CreationDate,    //TODO: not very pretty
			"key.ResourceOwner", //TODO: details
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

func KeyDetailsToPb(key *domain.MachineKey) []byte {
	details, err := json.Marshal(struct {
		Type   string `json:"type"`
		KeyID  string `json:"keyId"`
		Key    string `json:"key"`
		UserID string `json:"userId"`
	}{
		Type:   "serviceaccount",
		KeyID:  key.KeyID,
		Key:    string(key.PrivateKey),
		UserID: key.AggregateID,
	})
	logging.Log("AUTHN-sAiH5").OnError(err).Warn("unable to marshall key")

	return details
}
