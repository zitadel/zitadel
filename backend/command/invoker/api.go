package invoker

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/command"
	"github.com/zitadel/zitadel/backend/command/query"
	"github.com/zitadel/zitadel/backend/command/receiver"
	"github.com/zitadel/zitadel/backend/command/receiver/db"
	"github.com/zitadel/zitadel/backend/storage/database"
)

type api struct {
	db database.Pool

	manipulator receiver.InstanceManipulator
	reader      receiver.InstanceReader
}

func (a *api) CreateInstance(ctx context.Context) error {
	cmd := command.CreateInstance(db.NewInstance(a.db), &receiver.Instance{
		ID:   "123",
		Name: "test",
	})
	return cmd.Execute(ctx)
}

func (a *api) DeleteInstance(ctx context.Context) error {
	cmd := command.DeleteInstance(db.NewInstance(a.db), &receiver.Instance{
		ID: "123",
	})
	return cmd.Execute(ctx)
}

func (a *api) InstanceByID(ctx context.Context) (*receiver.Instance, error) {
	q := query.InstanceByID(a.reader, "123")
	return q.Execute(ctx)
}
