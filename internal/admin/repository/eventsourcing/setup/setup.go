package setup

import (
	"context"
	"github.com/caos/zitadel/internal/config/types"
)

type Setup struct {
	iamID       string
	setUpConfig types.IAMSetUp
}

func SetUp(ctx context.Context) error {
	return nil
}
