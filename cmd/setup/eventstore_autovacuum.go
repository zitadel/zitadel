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

	Enabled          bool   `json:"enabled"`
	VacuumThreshold  uint32 `json:"vacuumThreshold"`
	AnalyzeThreshold uint32 `json:"analyzeThreshold"`
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
	if mig.VacuumThreshold <= minAutovacuumThreshold || mig.AnalyzeThreshold <= minAutovacuumThreshold {
		return fmt.Errorf("%s: VacuumThreshold and AnalyzeThreshold must be greater than %d", mig, minAutovacuumThreshold)
	}

	// autovacuum_vacuum_threshold covers rows changed by updates and deletes, while
	// autovacuum_vacuum_insert_threshold only covers inserts. events2 is append-only, but
	// VacuumThreshold is applied to both, so a manual UPDATE/DELETE does not fall back to
	// the effectively disabled default scale factor.
	stmt := fmt.Sprintf(`ALTER TABLE eventstore.events2 SET (
	autovacuum_vacuum_scale_factor = 0.0,
	autovacuum_analyze_scale_factor = 0.0,
	autovacuum_vacuum_insert_scale_factor = 0.0,
	autovacuum_vacuum_insert_threshold = %d,
	autovacuum_analyze_threshold = %d,
	autovacuum_vacuum_threshold = %d
)`, mig.VacuumThreshold, mig.AnalyzeThreshold, mig.VacuumThreshold)
	_, err := mig.dbClient.ExecContext(ctx, stmt)
	return err
}

func (mig *eventstoreAutovacuum) String() string {
	return "repeatable_eventstore_autovacuum"
}

// Check implements [migration.RepeatableMigration].
func (mig *eventstoreAutovacuum) Check(lastRun map[string]interface{}) bool {
	enabled, _ := lastRun["enabled"].(bool)
	vacuumThreshold, _ := lastRun["vacuumThreshold"].(float64)
	analyzeThreshold, _ := lastRun["analyzeThreshold"].(float64)

	return enabled != mig.Enabled ||
		uint32(vacuumThreshold) != mig.VacuumThreshold ||
		uint32(analyzeThreshold) != mig.AnalyzeThreshold
}
