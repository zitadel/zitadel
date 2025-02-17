package step001

import (
	"context"
	"embed"
	"fmt"

	"github.com/zitadel/zitadel/backend/cmd/config"
	"github.com/zitadel/zitadel/backend/cmd/configure"
	"github.com/zitadel/zitadel/backend/storage/database"
)

var (
	//go:embed sql/*.sql
	migrations embed.FS
)

type Step001 struct {
	Database database.Pool `mapstructure:"-"`

	DatabaseName string `configure:"added:"v3",default:"zitadel"`
	Username     string `configure:"added:"v3",default:"zitadel"`
}

// Fields implements configure.StructUpdater.
func (v Step001) Fields() []configure.Updater {
	return []configure.Updater{
		configure.Field[string]{
			FieldName:   "databaseName",
			Default:     "zitadel",
			Value:       &v.DatabaseName,
			Description: "The name of the database Zitadel will store its data in",
			Version:     config.V3,
		},
		configure.Field[string]{
			FieldName:   "username",
			Default:     "zitadel",
			Value:       &v.Username,
			Description: "The username Zitadel will use to connect to the database",
			Version:     config.V3,
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
