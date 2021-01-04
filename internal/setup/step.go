package setup

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/business/command"
)

type step interface {
	step() iam_model.Step
	execute(context.Context) (*iam_model.IAM, error)
	init(*Setup)
	isNil() bool
}

type stepV2 interface {
	step() iam_model.Step
	execute(context.Context, string, command.CommandSide) error
	isNil() bool
}
