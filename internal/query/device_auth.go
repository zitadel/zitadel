package query

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
)

func (q *Queries) DeviceAuthByDeviceCode(ctx context.Context, clientID, deviceCode string) (_ *domain.DeviceAuth, err error) {
	return nil, nil
}

func (q *Queries) DeviceAuthByUserCode(ctx context.Context, userCode string) (_ *domain.DeviceAuth, err error) {
	return nil, nil
}
