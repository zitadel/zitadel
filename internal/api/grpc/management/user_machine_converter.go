package management

import (
	"github.com/caos/logging"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
)

func machineCreateToModel(machine *management.CreateMachineRequest) *usr_model.Machine {
	return &usr_model.Machine{
		Name:        machine.Name,
		Description: machine.Description,
	}
}

func updateMachineToUserModel(machine *management.UpdateMachineRequest) *usr_model.Machine {
	return &usr_model.Machine{
		Description: machine.Description,
	}
}

func serviceAccountFromUserModel(account *usr_model.User) *management.UserResponse {
	creationDate, err := ptypes.TimestampProto(account.CreationDate)
	logging.Log("MANAG-VwCfF").OnError(err).Debug("unable to parse creation date")

	changeDate, err := ptypes.TimestampProto(account.ChangeDate)
	logging.Log("MANAG-LELvM").OnError(err).Debug("unable to parse chagne date")

	return &management.UserResponse{
		Id:           account.AggregateID,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		State:        userStateFromModel(account.State),
		Sequence:     account.Sequence,
		User: &management.UserResponse_Machine{
			&management.MachineResponse{
				Name:        account.Name,
				Description: account.Description,
			},
		},
	}
}
