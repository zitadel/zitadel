package management

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
)

func createServiceAccountToUserModel(account *management.CreateServiceAccountRequest) *usr_model.User {
	return &usr_model.User{
		Machine: &usr_model.Machine{
			Name:        account.Name,
			Description: account.Description,
		},
	}
}

func updateServiceAccountToUserModel(account *management.UpdateServiceAccountRequest) *usr_model.User {
	return &usr_model.User{
		ObjectRoot: models.ObjectRoot{
			AggregateID: account.Id,
		},
		Machine: &usr_model.Machine{
			Description: account.Description,
		},
	}
}

func serviceAccountFromUserModel(account *usr_model.User) *management.ServiceAccountResponse {
	creationDate, err := ptypes.TimestampProto(account.CreationDate)
	logging.Log("MANAG-VwCfF").OnError(err).Debug("unable to parse creation date")

	changeDate, err := ptypes.TimestampProto(account.ChangeDate)
	logging.Log("MANAG-LELvM").OnError(err).Debug("unable to parse chagne date")

	return &management.ServiceAccountResponse{
		Id:            account.AggregateID,
		CreationDate:  creationDate,
		ChangeDate:    changeDate,
		ResourceOwner: account.ResourceOwner,
		Sequence:      account.Sequence,
		Name:          account.Name,
		Description:   account.Description,
	}
}
