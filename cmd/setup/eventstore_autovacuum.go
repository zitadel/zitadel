package setup

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

// eventstoreAutovacuum is a repeatable migration step that tunes PostgreSQL's
// autovacuum and autoanalyze behavior of the eventstore.events2 table,
// as configured by Eventstore.Autovacuum.
//
// It repeats whenever the configuration changed since the last run.
type eventstoreAutovacuum struct {
	dbClient *database.DB

	Enabled               bool   `json:"enabled"`
	VacuumInsertThreshold uint32 `json:"vacuumInsertThreshold"`
	AnalyzeThreshold      uint32 `json:"analyzeThreshold"`
}

const eventstoreAutovacuumResetStmt = `ALTER TABLE eventstore.events2 RESET (
	autovacuum_vacuum_scale_factor,
	autovacuum_analyze_scale_factor,
	autovacuum_vacuum_insert_scale_factor,
	autovacuum_vacuum_insert_threshold,
	autovacuum_analyze_threshold,
	autovacuum_vacuum_threshold
)`

// minAutovacuumThreshold guards against configuring vacuum/analyze runs so
// frequently that they thrash the database with constant autovacuum activity.
const minAutovacuumThreshold = 10000

func (mig *eventstoreAutovacuum) Execute(ctx context.Context, _ eventstore.Event) error {
	if !mig.Enabled {
		_, err := mig.dbClient.ExecContext(ctx, eventstoreAutovacuumResetStmt)
		return err
	}
	if mig.VacuumInsertThreshold <= minAutovacuumThreshold || mig.AnalyzeThreshold <= minAutovacuumThreshold {
		return fmt.Errorf("%s: VacuumInsertThreshold and AnalyzeThreshold must be greater than %d", mig, minAutovacuumThreshold)
	}

	// autovacuum_vacuum_threshold covers rows changed by updates and deletes. events2 is
	// append-only, but it is set to VacuumInsertThreshold as well, so that a manual
	// UPDATE/DELETE does not fall back to the effectively disabled default scale factor.
	stmt := fmt.Sprintf(`ALTER TABLE eventstore.events2 SET (
	autovacuum_vacuum_scale_factor = 0.0,
	autovacuum_analyze_scale_factor = 0.0,
	autovacuum_vacuum_insert_scale_factor = 0.0,
	autovacuum_vacuum_insert_threshold = %d,
	autovacuum_analyze_threshold = %d,
	autovacuum_vacuum_threshold = %d
)`, mig.VacuumInsertThreshold, mig.AnalyzeThreshold, mig.VacuumInsertThreshold)
	_, err := mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *eventstoreAutovacuum) String() string {
	return "repeatable_eventstore_autovacuum"
}

// Check implements [migration.RepeatableMigration].
func (mig *eventstoreAutovacuum) Check(lastRun map[string]interface{}) bool {
	enabled, _ := lastRun["enabled"].(bool)
	vacuumInsertThreshold, _ := lastRun["vacuumInsertThreshold"].(float64)
	analyzeThreshold, _ := lastRun["analyzeThreshold"].(float64)

	return enabled != mig.Enabled ||
		uint32(vacuumInsertThreshold) != mig.VacuumInsertThreshold ||
		uint32(analyzeThreshold) != mig.AnalyzeThreshold
}
