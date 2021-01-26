package auth

import (
	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/auth"
)

func machineViewFromModel(machine *usr_model.MachineView) *auth.MachineView {
	lastKeyAdded, err := ptypes.TimestampProto(machine.LastKeyAdded)
	logging.Log("MANAG-wGcAQ").OnError(err).Debug("unable to parse date")
	return &auth.MachineView{
		Description:  machine.Description,
		Name:         machine.Name,
		LastKeyAdded: lastKeyAdded,
	}
}

//
//func machineKeyViewsFromModel(keys ...*key_model.AuthNKeyView) []*auth.MachineKeyView {
//	keyViews := make([]*auth.MachineKeyView, len(keys))
//	for i, key := range keys {
//		keyViews[i] = machineKeyViewFromModel(key)
//	}
//	return keyViews
//}
//
//func machineKeyViewFromModel(key *key_model.AuthNKeyView) *auth.MachineKeyView {
//	creationDate, err := ptypes.TimestampProto(key.CreationDate)
//	logging.Log("MANAG-gluk7").OnError(err).Debug("unable to parse timestamp")
//
//	expirationDate, err := ptypes.TimestampProto(key.CreationDate)
//	logging.Log("MANAG-gluk7").OnError(err).Debug("unable to parse timestamp")
//
//	return &auth.MachineKeyView{
//		Id:             key.ID,
//		CreationDate:   creationDate,
//		ExpirationDate: expirationDate,
//		Sequence:       key.Sequence,
//		Type:           machineKeyTypeFromModel(key.Type),
//	}
//}
//
//func machineKeyTypeFromModel(typ key_model.AuthNKeyType) auth.MachineKeyType {
//	switch typ {
//	case key_model.AuthNKeyTypeJSON:
//		return auth.MachineKeyType_MACHINEKEY_JSON
//	default:
//		return auth.MachineKeyType_MACHINEKEY_UNSPECIFIED
//	}
//}
