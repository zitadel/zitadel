package query

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/receiver"
)

type instanceByID struct {
	receiver receiver.InstanceReader
	id       string
}

// InstanceByID returns a new instanceByID query.
func InstanceByID(receiver receiver.InstanceReader, id string) *instanceByID {
	return &instanceByID{
		receiver: receiver,
		id:       id,
	}
}

// Execute implements Query.
func (i *instanceByID) Execute(ctx context.Context) (*receiver.Instance, error) {
	return i.receiver.ByID(ctx, i.id)
}

// Name implements Query.
func (i *instanceByID) Name() string {
	return "instanceByID"
}

var _ Query[*receiver.Instance] = (*instanceByID)(nil)
