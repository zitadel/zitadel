package migration

import (
	"context"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/mitchellh/mapstructure"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	countTriggerTmpl       = "count_trigger"
	deleteParentCountsTmpl = "delete_parent_counts_trigger"
)

var (
	//go:embed *.sql
	templateFS embed.FS
	templates  = template.Must(template.ParseFS(templateFS, "*.sql"))
)

// CountTrigger registers the existing projections.count_trigger function.
// The trigger than takes care of keeping count of existing
// rows in the source table.
// It also pre-populates the projections.resource_counts table with
// the counts for the given table.
//
// During the population of the resource_counts table,
// the source table is share-locked to prevent concurrent modifications.
// Projection handlers will be halted until the lock is released.
// SELECT statements are not blocked by the lock.
//
// This migration repeats when any of the arguments are changed,
// such as renaming of a projection table.
func CountTrigger(
	db *database.DB,
	table string,
	parentType domain.CountParentType,
	instanceIDColumn string,
	parentIDColumn string,
	resource string,
) RepeatableMigration {
	return &triggerMigration{
		triggerConfig: triggerConfig{
			Table:            table,
			ParentType:       parentType.String(),
			InstanceIDColumn: instanceIDColumn,
			ParentIDColumn:   parentIDColumn,
			Resource:         resource,
		},
		db:           db,
		templateName: countTriggerTmpl,
	}
}

// DeleteParentCountsTrigger
//
// This migration repeats when any of the arguments are changed,
// such as renaming of a projection table.
func DeleteParentCountsTrigger(
	db *database.DB,
	table string,
	parentType domain.CountParentType,
	instanceIDColumn string,
	parentIDColumn string,
	resource string,
) RepeatableMigration {
	return &triggerMigration{
		triggerConfig: triggerConfig{
			Table:            table,
			ParentType:       parentType.String(),
			InstanceIDColumn: instanceIDColumn,
			ParentIDColumn:   parentIDColumn,
			Resource:         resource,
		},
		db:           db,
		templateName: deleteParentCountsTmpl,
	}
}

type triggerMigration struct {
	triggerConfig
	db           *database.DB
	templateName string
}

// String implements [Migration] and [fmt.Stringer].
func (m *triggerMigration) String() string {
	return fmt.Sprintf("repeatable_%s_%s", m.Resource, m.templateName)
}

// Execute implements [Migration]
func (m *triggerMigration) Execute(ctx context.Context, _ eventstore.Event) error {
	var query strings.Builder
	err := templates.ExecuteTemplate(&query, m.templateName, m.triggerConfig)
	if err != nil {
		return fmt.Errorf("%s: execute trigger template: %w", m, err)
	}
	_, err = m.db.ExecContext(ctx, query.String())
	if err != nil {
		return fmt.Errorf("%s: exec trigger query: %w", m, err)
	}
	return nil
}

type triggerConfig struct {
	Table            string `json:"table,omitempty" mapstructure:"table"`
	ParentType       string `json:"parent_type,omitempty" mapstructure:"parent_type"`
	InstanceIDColumn string `json:"instance_id_column,omitempty" mapstructure:"instance_id_column"`
	ParentIDColumn   string `json:"parent_id_column,omitempty" mapstructure:"parent_id_column"`
	Resource         string `json:"resource,omitempty" mapstructure:"resource"`
}

// Check implements [RepeatableMigration].
func (c *triggerConfig) Check(lastRun map[string]any) bool {
	var dst triggerConfig
	if err := mapstructure.Decode(lastRun, &dst); err != nil {
		panic(err)
	}
	return dst != *c
}
