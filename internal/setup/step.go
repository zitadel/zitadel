package setup

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type step interface {
	step() iam_model.Step
	execute(context.Context) error
	init(*Setup)
	isNil() bool
}
