package setup

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 39/cockroach/39_init_push_func.sql
	initPushFuncCRDB string
	//go:embed 39/postgres/39_init_push_func.sql
	initPushFuncPG string
)

type InitPushFunc struct {
	dbClient *database.DB
}

func (mig *InitPushFunc) Execute(ctx context.Context, _ eventstore.Event) (err error) {
	switch mig.dbClient.Type() {
	case "cockroach":
		_, err = mig.dbClient.ExecContext(ctx, initPushFuncCRDB)
	case "postgres":
		_, err = mig.dbClient.ExecContext(ctx, initPushFuncPG)
	default:
		err = fmt.Errorf("add cache schema: unsupported db type %q", mig.dbClient.Type())
	}
	return err
}

func (mig *InitPushFunc) String() string {
	return "39_init_push_func"
}
