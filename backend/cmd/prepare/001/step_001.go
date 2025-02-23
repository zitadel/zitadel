package step001

import (
	"context"
	"embed"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/zitadel/zitadel/backend/cmd/configure"
	"github.com/zitadel/zitadel/backend/storage/database"
)

var (
	//go:embed sql/*.sql
	migrations embed.FS
)

type Step001 struct {
	Database database.Pool `mapstructure:"-"`

	DatabaseName string
	Username     string
}

// Fields implements configure.StructUpdater.
func (v Step001) Fields() []configure.Updater {
	return []configure.Updater{
		&configure.Field[string]{
			FieldName:   "databaseName",
			Value:       "zitadel",
			Description: "The name of the database Zitadel will store its data in",
			Version:     semver.MustParse("3"),
		},
		&configure.Field[string]{
			FieldName:   "username",
			Value:       "zitadel",
			Description: "The username Zitadel will use to connect to the database",
			Version:     semver.MustParse("3"),
		},
	}
}

// Name implements configure.StructUpdater.
func (v *Step001) Name() string {
	return "step001"
}

// var _ configure.StructUpdater = (*Step001)(nil)

func (v *Step001) Migrate(ctx context.Context) error {
	files, err := migrations.ReadDir("sql")
	if err != nil {
		return err
	}
	for _, file := range files {
		fmt.Println(file.Name())
		fmt.Println(migrations.ReadFile(file.Name()))
	}
	conn, err := v.Database.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release(ctx)

	return nil
}
