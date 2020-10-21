package setup

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type step interface {
	step() iam_model.Step
	execute(context.Context) (*iam_model.IAM, error)
	init(*Setup)
	isNil() bool
}
