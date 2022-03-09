package handler

import "context"

type Check struct {
	Execute func(ex Executer, projectionName string) error
}

func (c *Check) IsNoop() bool {
	return c.Execute == nil
}

func (h *ProjectionHandler) Initialize(ctx context.Context, lock Lock, unlock Unlock, init Init, checks []*Check) error {
	inits := func(ctx context.Context) error {
		for _, check := range checks {
			return init(ctx, check)
		}
		return nil
	}
	return h.executeWithLock(ctx, lock, unlock, inits)
}
