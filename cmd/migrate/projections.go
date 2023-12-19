package migrate

import (
	"bytes"
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/key"
	"github.com/zitadel/zitadel/cmd/tls"
	admin_es "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing"
	admin_handler "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/handler"
	admin_view "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/view"
	internal_authz "github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	auth_es "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing"
	auth_handler "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/handler"
	auth_view "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/authz"
	authz_es "github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/eventstore"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	crypto_db "github.com/zitadel/zitadel/internal/crypto/database"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_es "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
	new_es "github.com/zitadel/zitadel/internal/eventstore/v3"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/notification"
	"github.com/zitadel/zitadel/internal/notification/handlers"
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
	static_config "github.com/zitadel/zitadel/internal/static/config"
	"github.com/zitadel/zitadel/internal/webauthn"
)

func projectionsCmd() *cobra.Command {
	err := viper.MergeConfig(bytes.NewBuffer(defaultConfig))
	logging.OnError(err).Fatal("unable to read setup steps")

	cmd := &cobra.Command{
		Use:   "projections",
		Short: "calls the projections synchronously",
		Run: func(cmd *cobra.Command, args []string) {
			config := mustNewProjectionsConfig(viper.GetViper())

			masterKey, err := key.MasterKey(cmd)
			logging.OnError(err).Fatal("unable to read master key")

			projections(cmd.Context(), config, masterKey)
		},
	}

	migrateProjectionsFlags(cmd)

	return cmd
}

type ProjectionsConfig struct {
	Destination    database.Config
	Projections    projection.Config
	EncryptionKeys *encryptionKeyConfig
	SystemAPIUsers SystemAPIUsers
	Eventstore     *eventstore.Config

	Admin admin_es.Config
	Auth  auth_es.Config

	Log     *logging.Config
	Machine *id.Config

	ExternalPort    uint16
	ExternalDomain  string
	ExternalSecure  bool
	InternalAuthZ   internal_authz.Config
	SystemDefaults  systemdefaults.SystemDefaults
	Telemetry       *handlers.TelemetryPusherConfig
	Login           login.Config
	OIDC            oidc.Config
	WebAuthNName    string
	DefaultInstance command.InstanceSetup
	AssetStorage    static_config.AssetStorageConfig
}

func migrateProjectionsFlags(cmd *cobra.Command) {
	key.AddMasterKeyFlag(cmd)
	tls.AddTLSModeFlag(cmd)
}

func projections(
	ctx context.Context,
	config *ProjectionsConfig,
	masterKey string,
) {
	start := time.Now()

	client, err := database.Connect(config.Destination, false, false)
	logging.OnError(err).Fatal("unable to connect to database")

	keyStorage, err := crypto_db.NewKeyStorage(client, masterKey)
	logging.OnError(err).Fatal("cannot start key storage")

	keys, err := ensureEncryptionKeys(config.EncryptionKeys, keyStorage)
	logging.OnError(err).Fatal("unable to read encryption keys")

	staticStorage, err := config.AssetStorage.NewStorage(client.DB)
	logging.OnError(err).Fatal("unable create static storage")

	config.Eventstore.Querier = old_es.NewCRDB(client)
	esPusherDBClient, err := database.Connect(config.Destination, false, true)
	logging.OnError(err).Fatal("unable to connect eventstore push client")
	config.Eventstore.Pusher = new_es.NewEventstore(esPusherDBClient)
	es := eventstore.NewEventstore(config.Eventstore)

	sessionTokenVerifier := internal_authz.SessionTokenVerifier(keys.OIDC)

	queries, err := query.StartQueries(
		ctx,
		es,
		client,
		config.Projections,
		config.SystemDefaults,
		keys.IDPConfig,
		keys.OTP,
		keys.OIDC,
		keys.SAML,
		config.InternalAuthZ.RolePermissionMappings,
		sessionTokenVerifier,
		func(q *query.Queries) domain.PermissionCheck {
			return func(ctx context.Context, permission, orgID, resourceID string) (err error) {
				return internal_authz.CheckPermission(ctx, &authz_es.UserMembershipRepo{Queries: q}, config.InternalAuthZ.RolePermissionMappings, permission, orgID, resourceID)
			}
		},
		0,
		config.SystemAPIUsers,
	)
	logging.OnError(err).Fatal("unable to start queries")

	authZRepo, err := authz.Start(queries, es, client, keys.OIDC, config.ExternalSecure)
	logging.OnError(err).Fatal("unable to start authz repo")

	webAuthNConfig := &webauthn.Config{
		DisplayName:    config.WebAuthNName,
		ExternalSecure: config.ExternalSecure,
	}
	commands, err := command.StartCommands(
		es,
		config.SystemDefaults,
		config.InternalAuthZ.RolePermissionMappings,
		staticStorage,
		webAuthNConfig,
		config.ExternalDomain,
		config.ExternalSecure,
		config.ExternalPort,
		keys.IDPConfig,
		keys.OTP,
		keys.SMTP,
		keys.SMS,
		keys.User,
		keys.DomainVerification,
		keys.OIDC,
		keys.SAML,
		&http.Client{},
		func(ctx context.Context, permission, orgID, resourceID string) (err error) {
			return internal_authz.CheckPermission(ctx, authZRepo, config.InternalAuthZ.RolePermissionMappings, permission, orgID, resourceID)
		},
		sessionTokenVerifier,
		config.OIDC.DefaultAccessTokenLifetime,
		config.OIDC.DefaultRefreshTokenExpiration,
		config.OIDC.DefaultRefreshTokenIdleExpiration,
		config.DefaultInstance.SecretGenerators,
	)
	logging.OnError(err).Fatal("unable to start commands")

	registerMappers(es)

	err = projection.Create(ctx, client, es, config.Projections, keys.OIDC, keys.SAML, config.SystemAPIUsers)
	logging.OnError(err).Fatal("unable to start projections")

	i18n.MustLoadSupportedLanguagesFromDir()

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

	config.Auth.Spooler.Client = client
	config.Auth.Spooler.Eventstore = es
	authView, err := auth_view.StartView(config.Auth.Spooler.Client, keys.OIDC, queries, config.Auth.Spooler.Eventstore)
	logging.OnError(err).Fatal("unable to start auth view")
	auth_handler.Register(ctx, config.Auth.Spooler, authView, queries)

	config.Admin.Spooler.Client = client
	config.Admin.Spooler.Eventstore = es
	adminView, err := admin_view.StartView(config.Admin.Spooler.Client)
	logging.OnError(err).Fatal("unable to start admin view")

	admin_handler.Register(ctx, config.Admin.Spooler, adminView, staticStorage)

	for _, instance := range queryInstanceIDs(ctx, client) {
		logging.WithFields("instance", instance).Info("projections")
		ctx = internal_authz.WithInstanceID(ctx, instance)

		err = projection.ProjectInstance(ctx)
		logging.OnError(err).Fatal("trigger failed")

		err = admin_handler.ProjectInstance(ctx)
		logging.OnError(err).Fatal("trigger admin handler failed")

		err = auth_handler.ProjectInstance(ctx)
		logging.OnError(err).Fatal("trigger auth handler failed")

		err = notification.ProjectInstance(ctx)
		logging.OnError(err).Fatal("trigger notification failed")
	}

	logging.WithFields("took", time.Since(start)).Info("projections executed")
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

// returns the instance configured by flag
// or all instances which are not removed
func queryInstanceIDs(ctx context.Context, source *database.DB) []string {
	if len(instanceIDs) > 0 {
		return instanceIDs
	}

	instances := []string{}
	err := source.QueryContext(
		ctx,
		func(r *sql.Rows) error {
			for r.Next() {
				var instance string

				if err := r.Scan(&instance); err != nil {
					return err
				}
				instances = append(instances, instance)
			}
			return r.Err()
		},
		"SELECT DISTINCT instance_id FROM eventstore.events2 WHERE instance_id <> '' AND aggregate_type = 'instance' AND event_type = 'instance.added' AND instance_id NOT IN (SELECT instance_id FROM eventstore.events2 WHERE instance_id <> '' AND aggregate_type = 'instance' AND event_type = 'instance.removed')",
	)
	logging.OnError(err).Fatal("unable to query instances")

	return instances
}
