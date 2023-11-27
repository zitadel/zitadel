package migrate

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/key"
	"github.com/zitadel/zitadel/internal/api/authz"
	crypto_db "github.com/zitadel/zitadel/internal/crypto/database"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_es "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/action"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	iam_repo "github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/keypair"
	"github.com/zitadel/zitadel/internal/repository/limits"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/repository/restrictions"
	"github.com/zitadel/zitadel/internal/repository/session"
	usr_repo "github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
)

func projectionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "projections",
		Short: "calls the projections synchronously",
		Run: func(cmd *cobra.Command, args []string) {
			config := mustNewProjectionsConfig(viper.GetViper())
			ctx := authz.WithInstanceID(cmd.Context(), instanceID)

			masterKey, err := key.MasterKey(cmd)
			logging.OnError(err).Fatal("unable to read master key")

			projections(ctx, config, masterKey)
		},
	}

	migrateProjectionsFlags(cmd)

	return cmd
}

type ProjectionsConfig struct {
	Database       database.Config
	Projections    projection.Config
	EncryptionKeys *encryptionKeyConfig
	SystemAPIUsers SystemAPIUsers
	Eventstore     *eventstore.Config

	Log     *logging.Config
	Machine *id.Config
}

func migrateProjectionsFlags(cmd *cobra.Command) {
	key.AddMasterKeyFlag(cmd)
	cmd.Flags().StringArrayVar(&configPaths, "config", nil, "paths to config files")
}

func projections(
	ctx context.Context,
	config *ProjectionsConfig,
	masterKey string,
) {
	start := time.Now()

	client, err := database.Connect(config.Database, false, false)
	logging.OnError(err).Fatal("unable to connect to database")

	keyStorage, err := crypto_db.NewKeyStorage(client, masterKey)
	logging.OnError(err).Fatal("cannot start key storage")

	keys, err := ensureEncryptionKeys(config.EncryptionKeys, keyStorage)
	logging.OnError(err).Fatal("unable to read encryption keys")

	config.Eventstore.Querier = old_es.NewCRDB(client)
	es := eventstore.NewEventstore(config.Eventstore)

	iam_repo.RegisterEventMappers(es)
	usr_repo.RegisterEventMappers(es)
	org.RegisterEventMappers(es)
	project.RegisterEventMappers(es)
	action.RegisterEventMappers(es)
	keypair.RegisterEventMappers(es)
	usergrant.RegisterEventMappers(es)
	session.RegisterEventMappers(es)
	idpintent.RegisterEventMappers(es)
	authrequest.RegisterEventMappers(es)
	oidcsession.RegisterEventMappers(es)
	quota.RegisterEventMappers(es)
	limits.RegisterEventMappers(es)
	restrictions.RegisterEventMappers(es)

	err = projection.Create(ctx, client, es, config.Projections, keys.OIDC, keys.SAML, config.SystemAPIUsers)
	logging.OnError(err).Fatal("unable to start projections")

	err = projection.ProjectInstance(ctx)
	logging.OnError(err).Fatal("trigger failed")

	logging.WithFields("took", time.Since(start)).Info("projections executed")
}
