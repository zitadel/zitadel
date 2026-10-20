package setup

import (
	"context"
	"embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 76/*.sql
	users14LoginEqualityIndexes embed.FS
)

const (
	users14UsernameLowerIdx    = "users14_username_lower_idx"
	users14HumansPhoneLowerIdx = "users14_humans_phone_lower_idx"
)

var users14LoginEqualityIndexNames = []string{
	users14UsernameLowerIdx,
	users14HumansPhoneLowerIdx,
}

type Users14LoginEqualityIndexes struct {
	dbClient *database.DB
}

func (mig *Users14LoginEqualityIndexes) Execute(ctx context.Context, _ eventstore.Event) error {
	return executeConcurrentIndexSQL(ctx, mig.dbClient, mig.String(), users14TableExistsQuery, users14LoginEqualityIndexNames, users14LoginEqualityIndexes, "76")
}

func (mig *Users14LoginEqualityIndexes) String() string {
	return "76_users14_login_equality_indexes"
}
