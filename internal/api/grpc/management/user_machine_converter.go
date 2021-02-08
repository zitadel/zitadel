package management

import (
	"encoding/json"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/v2/domain"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func machineCreateToDomain(machine *management.CreateMachineRequest) *domain.Machine {
	return &domain.Machine{
		Name:        machine.Name,
		Description: machine.Description,
	}
}

func updateMachineToDomain(ctxData authz.CtxData, machine *management.UpdateMachineRequest) *domain.Machine {
	return &domain.Machine{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   machine.Id,
			ResourceOwner: ctxData.ResourceOwner,
		},
		Name:        machine.Name,
		Description: machine.Description,
	}
}

func machineFromDomain(account *domain.Machine) *management.MachineResponse {
	return &management.MachineResponse{
		Name:        account.Name,
		Description: account.Description,
	}
}

func machineViewFromModel(machine *usr_model.MachineView) *management.MachineView {
	lastKeyAdded, err := ptypes.TimestampProto(machine.LastKeyAdded)
	logging.Log("MANAG-wGcAQ").OnError(err).Debug("unable to parse date")
	return &management.MachineView{
		Description:  machine.Description,
		Name:         machine.Name,
		LastKeyAdded: lastKeyAdded,
	}
}

func machineKeyViewsFromModel(keys ...*usr_model.MachineKeyView) []*management.MachineKeyView {
	keyViews := make([]*management.MachineKeyView, len(keys))
	for i, key := range keys {
		keyViews[i] = machineKeyViewFromModel(key)
	}
	return keyViews
}

func machineKeyViewFromModel(key *usr_model.MachineKeyView) *management.MachineKeyView {
	creationDate, err := ptypes.TimestampProto(key.CreationDate)
	logging.Log("MANAG-gluk7").OnError(err).Debug("unable to parse timestamp")

	expirationDate, err := ptypes.TimestampProto(key.ExpirationDate)
	logging.Log("MANAG-gluk7").OnError(err).Debug("unable to parse timestamp")

	return &management.MachineKeyView{
		Id:             key.ID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
		Sequence:       key.Sequence,
		Type:           machineKeyTypeFromModel(key.Type),
	}
}

func addMachineKeyToDomain(key *management.AddMachineKeyRequest) *domain.MachineKey {
	expirationDate := time.Time{}
	if key.ExpirationDate != nil {
		var err error
		expirationDate, err = ptypes.Timestamp(key.ExpirationDate)
		logging.Log("MANAG-iNshR").OnError(err).Debug("unable to parse expiration date")
	}

	return &domain.MachineKey{
		ExpirationDate: expirationDate,
		Type:           machineKeyTypeToDomain(key.Type),
		ObjectRoot:     models.ObjectRoot{AggregateID: key.UserId},
	}
}

func addMachineKeyFromDomain(key *domain.MachineKey) *management.AddMachineKeyResponse {
	detail, err := json.Marshal(struct {
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
	logging.Log("MANAG-lFQ2g").OnError(err).Warn("unable to marshall key")

	return &management.AddMachineKeyResponse{
		Id:             key.KeyID,
		CreationDate:   timestamppb.New(key.CreationDate),
		ExpirationDate: timestamppb.New(key.ExpirationDate),
		Sequence:       key.Sequence,
		KeyDetails:     detail,
		Type:           machineKeyTypeFromDomain(key.Type),
	}
}

func machineKeyTypeToDomain(typ management.MachineKeyType) domain.MachineKeyType {
	switch typ {
	case management.MachineKeyType_MACHINEKEY_JSON:
		return domain.MachineKeyTypeJSON
	default:
		return domain.MachineKeyTypeNONE
	}
}

func machineKeyTypeFromDomain(typ domain.MachineKeyType) management.MachineKeyType {
	switch typ {
	case domain.MachineKeyTypeJSON:
		return management.MachineKeyType_MACHINEKEY_JSON
	default:
		return management.MachineKeyType_MACHINEKEY_UNSPECIFIED
	}
}

func machineKeyTypeFromModel(typ usr_model.MachineKeyType) management.MachineKeyType {
	switch typ {
	case usr_model.MachineKeyTypeJSON:
		return management.MachineKeyType_MACHINEKEY_JSON
	default:
		return management.MachineKeyType_MACHINEKEY_UNSPECIFIED
	}
}

func machineKeySearchRequestToModel(req *management.MachineKeySearchRequest) *usr_model.MachineKeySearchRequest {
	return &usr_model.MachineKeySearchRequest{
		Offset: req.Offset,
		Limit:  req.Limit,
		Asc:    req.Asc,
		Queries: []*usr_model.MachineKeySearchQuery{
			{
				Key:    usr_model.MachineKeyKeyUserID,
				Method: model.SearchMethodEquals,
				Value:  req.UserId,
			},
		},
	}
}

func machineKeySearchResponseFromModel(req *usr_model.MachineKeySearchResponse) *management.MachineKeySearchResponse {
	viewTimestamp, err := ptypes.TimestampProto(req.Timestamp)
	logging.Log("MANAG-Sk9ds").OnError(err).Debug("unable to parse cretaion date")

	return &management.MachineKeySearchResponse{
		Offset:            req.Offset,
		Limit:             req.Limit,
		TotalResult:       req.TotalResult,
		ProcessedSequence: req.Sequence,
		ViewTimestamp:     viewTimestamp,
		Result:            machineKeyViewsFromModel(req.Result...),
	}
}
