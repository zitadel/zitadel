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

type StampEventPositionAtInsert struct{}

func (*StampEventPositionAtInsert) Execute(context.Context, eventstore.Event) error {
	return errors.New("setup 77 is historical and must not run")
}

func (*StampEventPositionAtInsert) String() string {
	return migration.EventstorePositionClockTimestampStep
}
