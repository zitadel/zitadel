package setup

import (
	"context"

	"github.com/caos/zitadel/internal/command/v2"
)

type DefaultInstance struct {
	cmd           *command.Command
	InstanceSetup command.InstanceSetup
}

func (mig *DefaultInstance) Execute(ctx context.Context) error {
	_, err := mig.cmd.SetUpInstance(ctx, &mig.InstanceSetup)

	return err
}

func (mig *DefaultInstance) String() string {
	return "03_default_instance"
}
