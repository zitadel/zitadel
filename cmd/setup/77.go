package setup

import (
	"context"
	_ "embed"
	"errors"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/migration"
)

var (
	//go:embed 77.sql
	stampEventPositionAtInsert string
)

// StampEventPositionAtInsert is historical setup 77. It is no longer executed.
// Existing Done events with this name remain the flag that filter_offset already
// stores events2.in_tx_order. New installs never record this name.
type StampEventPositionAtInsert struct{}

func (*StampEventPositionAtInsert) Execute(context.Context, eventstore.Event) error {
	return errors.New("setup 77 is historical and must not run")
}

func (*StampEventPositionAtInsert) String() string {
	return migration.EventstorePositionClockTimestampStep
}
