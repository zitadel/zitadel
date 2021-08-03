package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
)

type Step4 struct {
	DefaultPasswordLockoutPolicy domain.LockoutPolicy
}

func (s *Step4) Step() domain.Step {
	return domain.Step4
}

func (s *Step4) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep4(ctx, s)
}

func (c *Commands) SetupStep4(ctx context.Context, step *Step4) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		return nil, nil
	}
	return c.setup(ctx, step, fn)
}
