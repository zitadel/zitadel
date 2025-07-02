package setup

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	eventstoreV3 "github.com/zitadel/zitadel/internal/eventstore/v3"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type AddIDUniqueConstraintsForOrgs struct {
	dbClient   *database.DB
	eventstore *eventstore.Eventstore
}

type orgsReadModel struct {
	eventstore.ReadModel
}

func (orm *orgsReadModel) AppendEvents(events ...eventstore.Event) {
	orm.ReadModel.AppendEvents(events...)
}

func (orm *orgsReadModel) Reduce() error {
	return nil
}

type OrgAddEventUpdateIDUniqueConstraint struct {
	*org.OrgAddedEvent
}

func (e *OrgAddEventUpdateIDUniqueConstraint) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{org.NewAddOrgIDUniqueConstraint(e.Aggregate().ID)}
}

type OrgRemoveEventUpdateIDUniqueConstraint struct {
	*org.OrgRemovedEvent
}

func (e *OrgRemoveEventUpdateIDUniqueConstraint) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{org.NewRemoveOrgIDUniqueConstraint(e.Aggregate().ID)}
}

func (mig *AddIDUniqueConstraintsForOrgs) Execute(ctx context.Context, _ eventstore.Event) (err error) {
	orm := orgsReadModel{}
	sqb := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		EventTypes(
			org.OrgAddedEventType,
			org.OrgRemovedEventType,
		).
		Builder()

	err = mig.eventstore.FilterToReducer(ctx, sqb, &orm)
	if err != nil {
		return err
	}

	var tx *sql.Tx
	tx, err = mig.dbClient.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		}
	}()

	for _, event := range orm.Events {
		switch event := event.(type) {
		case *org.OrgAddedEvent:
			orgAddUpdateUniqueIDConstraint := OrgAddEventUpdateIDUniqueConstraint{event}
			err = eventstoreV3.HandleUniqueConstraints(ctx, tx, []eventstore.Command{&orgAddUpdateUniqueIDConstraint})
			if err != nil {
				return err
			}
		case *org.OrgRemovedEvent:
			orgRemoveUpdateUniqueIDConstraint := OrgRemoveEventUpdateIDUniqueConstraint{event}
			err = eventstoreV3.HandleUniqueConstraints(ctx, tx, []eventstore.Command{&orgRemoveUpdateUniqueIDConstraint})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (*AddIDUniqueConstraintsForOrgs) String() string {
	return "61_add_id_unique_constraints_for_orgs"
}

func (f *AddIDUniqueConstraintsForOrgs) Check(lastRun map[string]interface{}) bool {
	return true
}
