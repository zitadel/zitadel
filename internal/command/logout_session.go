package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/repository/sessionlogout"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) BackChannelLogoutSent(ctx context.Context, id, oidcSessionID, instanceID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	sessionWriteModel := NewSessionLogoutWriteModel(id, instanceID, oidcSessionID)
	if err = c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel); err != nil {
		return err
	}

	return c.pushAppendAndReduce(
		ctx,
		sessionWriteModel,
		sessionlogout.NewBackChannelLogoutSentEvent(ctx, sessionWriteModel.aggregate, oidcSessionID),
	)
}
