package setup

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/owner"
	"github.com/zitadel/zitadel/internal/repository/project"
)

var (
	//go:embed 45.sql
	correctProjectOwnerEvents string
)

type CorrectProjectOwners struct {
	eventstore *eventstore.Eventstore
}

func (mig *CorrectProjectOwners) Execute(ctx context.Context, _ eventstore.Event) error {
	instances, err := mig.eventstore.InstanceIDs(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsInstanceIDs).
			OrderDesc().
			AddQuery().
			AggregateTypes("instance").
			EventTypes(instance.InstanceAddedEventType).
			Builder(),
	)
	if err != nil {
		return err
	}

	ctx = authz.SetCtxData(ctx, authz.CtxData{UserID: "SETUP"})
	for i, instance := range instances {
		ctx = authz.WithInstanceID(ctx, instance)
		logging.WithFields("instance_id", instance, "migration", mig.String(), "progress", fmt.Sprintf("%d/%d", i+1, len(instances))).Info("correct owners of projects")
		didCorrect, err := mig.correctInstanceProjects(ctx, instance)
		if err != nil {
			return err
		}
		if !didCorrect {
			continue
		}
		_, err = projection.ProjectGrantProjection.Trigger(ctx)
		logging.OnError(err).Debug("failed triggering project grant projection to update owners")
	}
	return nil
}

func (mig *CorrectProjectOwners) correctInstanceProjects(ctx context.Context, instance string) (didCorrect bool, err error) {
	var correctedOwners []eventstore.Command

	tx, err := mig.eventstore.Client().BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	rows, err := tx.QueryContext(ctx, correctProjectOwnerEvents, instance)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		aggregate := &eventstore.Aggregate{
			InstanceID: instance,
			Type:       project.AggregateType,
			Version:    project.AggregateVersion,
		}
		var payload json.RawMessage
		err := rows.Scan(
			&aggregate.ID,
			&aggregate.ResourceOwner,
			&payload,
		)
		if err != nil {
			return false, err
		}
		previousOwners := make(map[uint32]string)
		if err := json.Unmarshal(payload, &previousOwners); err != nil {
			return false, err
		}
		correctedOwners = append(correctedOwners, owner.NewCorrected(ctx, aggregate, previousOwners))
	}
	if rows.Err() != nil {
		return false, rows.Err()
	}

	_, err = mig.eventstore.PushWithClient(ctx, tx, correctedOwners...)
	return len(correctedOwners) > 0, err
}

func (*CorrectProjectOwners) String() string {
	return "43_correct_project_owners"
}
