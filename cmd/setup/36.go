package setup

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/milestone"
)

var (
	//go:embed 36.sql
	getProjectedMilestones string
)

type FillV3Milestones struct {
	dbClient   *database.DB
	eventstore *eventstore.Eventstore
}

type instanceMilestone struct {
	Type    milestone.Type
	Reached time.Time
	Pushed  *time.Time
}

func (mig *FillV3Milestones) Execute(ctx context.Context, _ eventstore.Event) error {
	im, err := mig.getProjectedMilestones(ctx)
	if err != nil {
		return err
	}
	return mig.pushEventsByInstance(ctx, im)
}

func (mig *FillV3Milestones) getProjectedMilestones(ctx context.Context) (map[string][]instanceMilestone, error) {
	type row struct {
		InstanceID string
		Type       milestone.Type
		Reached    time.Time
		Pushed     *time.Time
	}

	rows, _ := mig.dbClient.Pool.Query(ctx, getProjectedMilestones)
	scanned, err := pgx.CollectRows(rows, pgx.RowToStructByPos[row])
	var pgError *pgconn.PgError
	// catch ERROR:  relation "projections.milestones" does not exist
	if errors.As(err, &pgError) && pgError.SQLState() == "42P01" {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("milestones get: %w", err)
	}
	milestoneMap := make(map[string][]instanceMilestone)
	for _, s := range scanned {
		milestoneMap[s.InstanceID] = append(milestoneMap[s.InstanceID], instanceMilestone{
			Type:    s.Type,
			Reached: s.Reached,
			Pushed:  s.Pushed,
		})
	}
	return milestoneMap, nil
}

// pushEventsByInstance creates the v2 milestone events by instance.
// This prevents we will try to push 6*N(instance) events in one push.
func (mig *FillV3Milestones) pushEventsByInstance(ctx context.Context, milestoneMap map[string][]instanceMilestone) error {
	// keep a deterministic order by instance ID.
	order := make([]string, 0, len(milestoneMap))
	for k := range milestoneMap {
		order = append(order, k)
	}
	slices.Sort(order)

	for i, instanceID := range order {
		logging.WithFields("instance_id", instanceID, "migration", mig.String(), "progress", fmt.Sprintf("%d/%d", i+1, len(order))).Info("filter existing milestone events")

		// because each Push runs in a separate TX, we need to make sure that events
		// from a partially executed migration are pushed again.
		model := command.NewMilestonesReachedWriteModel(instanceID)
		if err := mig.eventstore.FilterToQueryReducer(ctx, model); err != nil {
			return fmt.Errorf("milestones filter: %w", err)
		}
		if model.InstanceCreated {
			logging.WithFields("instance_id", instanceID, "migration", mig.String()).Info("milestone events already migrated")
			continue // This instance was migrated, skip
		}
		logging.WithFields("instance_id", instanceID, "migration", mig.String()).Info("push milestone events")

		aggregate := milestone.NewInstanceAggregate(instanceID)

		cmds := make([]eventstore.Command, 0, len(milestoneMap[instanceID])*2)
		for _, m := range milestoneMap[instanceID] {
			cmds = append(cmds, milestone.NewReachedEventWithDate(ctx, aggregate, m.Type, &m.Reached))
			if m.Pushed != nil {
				cmds = append(cmds, milestone.NewPushedEventWithDate(ctx, aggregate, m.Type, nil, "", m.Pushed))
			}
		}

		if _, err := mig.eventstore.Push(ctx, cmds...); err != nil {
			return fmt.Errorf("milestones push: %w", err)
		}
	}
	return nil
}

func (mig *FillV3Milestones) String() string {
	return "36_fill_v3_milestones"
}
