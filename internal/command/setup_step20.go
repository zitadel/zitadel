package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
)

type Step20 struct{}

func (s *Step20) Step() domain.Step {
	return domain.Step20
}

func (s *Step20) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep20(ctx, s)
}

func (c *Commands) SetupStep20(ctx context.Context, step *Step20) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		err := c.eventstore.Step20(ctx, iam.Events[len(iam.Events)-1].Sequence())
		return nil, err
	}
	return c.setup(ctx, step, fn)
}
