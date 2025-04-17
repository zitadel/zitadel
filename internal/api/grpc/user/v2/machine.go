package user

import (
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func patchMachineUserToCommand(userId string, userName *string, machine *user.UpdateUserRequest_Machine) *command.ChangeMachine {
	return &command.ChangeMachine{
		ID:          userId,
		Username:    userName,
		Name:        machine.Name,
		Description: machine.Description,
	}
}
