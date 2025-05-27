package migration

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed count_trigger.sql
	countTriggerTmpl   string
	countTriggerParsed = template.Must(
		template.New("count_trigger").Parse(countTriggerTmpl),
	)
)

type countTriggerMigration struct {
	db               *database.DB
	Table            string
	ParentType       domain.CountParentType
	Resource         string
	InstanceIDColumn string
	OwnerIDColumn    string
}

func CountTriggerMigration(
	db *database.DB,
	table string,
	parentType domain.CountParentType,
	resource string,
	instanceIDColumn string,
	ownerIDColumn string,
) Migration {
	return &countTriggerMigration{
		db:               db,
		Table:            table,
		ParentType:       parentType,
		Resource:         resource,
		InstanceIDColumn: instanceIDColumn,
		OwnerIDColumn:    ownerIDColumn,
	}
}

func (m *countTriggerMigration) String() string {
	return fmt.Sprintf("init_count_trigger_%s", m.Table)
}

func (m *countTriggerMigration) Execute(ctx context.Context, _ eventstore.Event) error {
	var buf strings.Builder
	err := countTriggerParsed.Execute(&buf, m)
	if err != nil {
		return fmt.Errorf("execute count trigger template: %w", err)
	}
	_, err = m.db.ExecContext(ctx, buf.String())
	if err != nil {
		return fmt.Errorf("exec count trigger: %w", err)
	}
	return nil
}
