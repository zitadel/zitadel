package models

import (
	"context"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/errors"
)

type AggregateCreator struct {
	serviceName string
	//ignoreCtxData is needed to ignore ctxData.IsZero() in tests
	ignoreCtxData bool
}

func NewAggregateCreator(serviceName string) *AggregateCreator {
	return &AggregateCreator{serviceName: serviceName}
}

func (c *AggregateCreator) NewAggregate(ctx context.Context, id string, typ AggregateType, version Version, latestSequence uint64) (*Aggregate, error) {
	if id == "" {
		return nil, errors.ThrowInvalidArgument(nil, "MODEL-sPLP8", "no id")
	}
	if string(typ) == "" {
		return nil, errors.ThrowInvalidArgument(nil, "MODEL-yfmjm", "no type")
	}
	if err := version.Validate(); err != nil {
		return nil, err
	}

	ctxData := auth.GetCtxData(ctx)
	if ctxData.IsZero() && !c.ignoreCtxData {
		return nil, errors.ThrowInvalidArgument(nil, "MODEL-lZkk9", "ctxData zero")
	}

	return &Aggregate{
		id:             id,
		typ:            typ,
		latestSequence: latestSequence,
		version:        version,
		Events:         make([]*Event, 0, 2),
		editorOrg:      ctxData.OrgID,
		editorService:  c.serviceName,
		editorUser:     ctxData.UserID,
		resourceOwner:  ctxData.OrgID,
	}, nil
}
