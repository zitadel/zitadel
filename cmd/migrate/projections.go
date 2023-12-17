package migrate

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/key"
	admin_handler "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/handler"
	admin_view "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/api/authz"
	auth_handler "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/handler"
	auth_view "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/command"
	crypto_db "github.com/zitadel/zitadel/internal/crypto/database"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_es "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/notification"
	"github.com/zitadel/zitadel/internal/query"
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
	"github.com/zitadel/zitadel/internal/static"
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

	Admin admin_handler.Config
	Auth  auth_handler.Config

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

	staticStorage, err := static.CreateStorage(client, config.Static)
	logging.OnError(err).Fatal("unable create static storage")

	config.Eventstore.Querier = old_es.NewCRDB(client)
	es := eventstore.NewEventstore(config.Eventstore)

	queries, err := query.StartQueries(
		ctx,
		es,
		client,
		config.Projections,
		defaults,
		keys.idpConfigEncryption,
		keys.otpConfigEncryption,
		nil,
		nil,
		config.RoleMapping,
		sessionTokenVerifier,
		permissionCheck,
		0,
		nil,
	)
	logging.OnError(err).Fatal("unable to start queries")

	registerMappers(es)

	projectProjections(ctx, client, es, keys, config)
	projectAdmin(ctx, config.Admin, staticStorage, client)
	projectAuth(ctx, config.Auth, queries, es, client)
	projectNotification(ctx, es, keys, config.Projections)

	logging.WithFields("took", time.Since(start)).Info("projections executed")
}

func projectProjections(ctx context.Context, client *database.DB, es *eventstore.Eventstore, keys *encryptionKeys, config *ProjectionsConfig) {
	err := projection.Create(ctx, client, es, config.Projections, keys.OIDC, keys.SAML, config.SystemAPIUsers)
	logging.OnError(err).Fatal("unable to start projections")

	err = projection.ProjectInstance(ctx)
	logging.OnError(err).Fatal("trigger failed")
}

func projectNotification(ctx context.Context, es *eventstore.Eventstore, queries *query.Queries, commands *command.Commands, keys *encryptionKeys, config *ProjectionsConfig) {
	notification.Register(
		ctx,
		config.Projections.Customizations["notifications"],
		config.Projections.Customizations["notificationsquotas"],
		config.Projections.Customizations["telemetry"],
		*config.Telemetry,
		config.ExternalDomain,
		config.ExternalPort,
		config.ExternalSecure,
		commands,
		queries,
		es,
		config.Login.DefaultOTPEmailURLV2,
		config.SystemDefaults.Notifications.FileSystemPath,
		keys.User,
		keys.SMTP,
		keys.SMS,
	)

	err := notification.ProjectInstance(ctx)
	logging.OnError(err).Fatal("trigger notification failed")
}

func projectAuth(ctx context.Context, config auth_handler.Config, queries *query.Queries, es *eventstore.Eventstore, client *database.DB) {
	view, err := auth_view.StartView(client, oidcEncryption, queries, es)
	logging.OnError(err).Fatal("unable to start auth view")

	auth_handler.Register(ctx, config, view, queries)

	err = auth_handler.ProjectInstance(ctx)
	logging.OnError(err).Fatal("trigger auth handler failed")
}

func projectAdmin(ctx context.Context, config admin_handler.Config, staticStorage static.Storage, client *database.DB) {
	view, err := admin_view.StartView(client)
	logging.OnError(err).Fatal("unable to start admin view")

	admin_handler.Register(ctx, config, view, staticStorage)

	err = admin_handler.ProjectInstance(ctx)
	logging.OnError(err).Fatal("trigger admin handler failed")
}

func registerMappers(es *eventstore.Eventstore) {
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
}
