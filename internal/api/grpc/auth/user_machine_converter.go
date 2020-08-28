package auth

import (
	"github.com/caos/logging"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/auth"
	"github.com/golang/protobuf/ptypes"
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

func machineKeyViewsFromModel(keys ...*usr_model.MachineKeyView) []*auth.MachineKeyView {
	keyViews := make([]*auth.MachineKeyView, len(keys))
	for i, key := range keys {
		keyViews[i] = machineKeyViewFromModel(key)
	}
	return keyViews
}

func machineKeyViewFromModel(key *usr_model.MachineKeyView) *auth.MachineKeyView {
	creationDate, err := ptypes.TimestampProto(key.CreationDate)
	logging.Log("MANAG-gluk7").OnError(err).Debug("unable to parse timestamp")

	expirationDate, err := ptypes.TimestampProto(key.CreationDate)
	logging.Log("MANAG-gluk7").OnError(err).Debug("unable to parse timestamp")

	return &auth.MachineKeyView{
		Id:             key.ID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
		Sequence:       key.Sequence,
		Type:           auth.MachineKeyType_MACHINEKEY_UNSPECIFIED, //TODO: type mapping
	}
}
