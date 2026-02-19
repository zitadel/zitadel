package mirror

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	old_logging "github.com/zitadel/logging" //nolint:staticcheck

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/cmd/encryption"
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
	"github.com/zitadel/zitadel/internal/cache/connector"
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
	static_config "github.com/zitadel/zitadel/internal/static/config"
	es_v4 "github.com/zitadel/zitadel/internal/v2/eventstore"
	es_v4_pg "github.com/zitadel/zitadel/internal/v2/eventstore/postgres"
	"github.com/zitadel/zitadel/internal/webauthn"
)

func projectionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "projections",
		Short: "calls the projections synchronously",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				logging.OnError(cmd.Context(), err).Error("zitadel mirror projections command failed")
			}()
			config, shutdown, err := newProjectionsConfig(cmd, viper.GetViper())
			if err != nil {
				return fmt.Errorf("unable to create projections config: %w", err)
			}
			defer func() {
				err = errors.Join(err, shutdown(cmd.Context()))
			}()

			masterKey, err := key.MasterKey(cmd)
			if err != nil {
				return fmt.Errorf("unable to read master key: %w", err)
			}
			projections(cmd.Context(), config, masterKey)
			return nil
		},
	}

	migrateProjectionsFlags(cmd)

	return cmd
}

type ProjectionsConfig struct {
	Instrumentation instrumentation.Config
	Destination     database.Config
	Projections     projection.Config
	Notifications   handlers.WorkerConfig
	EncryptionKeys  *encryption.EncryptionKeyConfig
	SystemAPIUsers  map[string]*internal_authz.SystemAPIUser
	Eventstore      *eventstore.Config
	Caches          *connector.CachesConfig

	Admin admin_es.Config
	Auth  auth_es.Config

	Log     *old_logging.Config
	Machine *id.Config

	ExternalPort    uint16
	ExternalDomain  string
	ExternalSecure  bool
	InternalAuthZ   internal_authz.Config
	SystemAuthZ     internal_authz.Config
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
) error {
	logging.Info(ctx, "starting to fill projections")
	start := time.Now()

	client, err := database.Connect(config.Destination, false)
	logging.OnError(ctx, err).Fatal("unable to connect to database")

	keyStorage, err := crypto_db.NewKeyStorage(client, masterKey)
	logging.OnError(ctx, err).Fatal("cannot start key storage")

	keys, err := encryption.EnsureEncryptionKeys(ctx, config.EncryptionKeys, keyStorage)
	logging.OnError(ctx, err).Fatal("unable to read encryption keys")

	staticStorage, err := config.AssetStorage.NewStorage(client.DB)
	logging.OnError(ctx, err).Fatal("unable create static storage")

	newEventstore := new_es.NewEventstore(client)
	config.Eventstore.Querier = old_es.NewPostgres(client)
	config.Eventstore.Pusher = newEventstore
	config.Eventstore.Searcher = newEventstore

	es := eventstore.NewEventstore(config.Eventstore)
	esV4 := es_v4.NewEventstoreFromOne(es_v4_pg.New(client, &es_v4_pg.Config{
		MaxRetries: config.Eventstore.MaxRetries,
	}))

	sessionTokenVerifier := internal_authz.SessionTokenVerifier(keys.OIDC)

	cacheConnectors, err := connector.StartConnectors(config.Caches, client)
	logging.OnError(ctx, err).Fatal("unable to start caches")

	queries, err := query.StartQueries(
		ctx,
		es,
		esV4.Querier,
		client,
		client,
		cacheConnectors,
		config.Projections,
		config.SystemDefaults,
		keys.IDPConfig,
		keys.OTP,
		keys.OIDC,
		keys.SAML,
		keys.Target,
		keys.SMS,
		keys.SMTP,
		config.InternalAuthZ.RolePermissionMappings,
		sessionTokenVerifier,
		func(q *query.Queries) domain.PermissionCheck {
			return func(ctx context.Context, permission, orgID, resourceID string) (err error) {
				return internal_authz.CheckPermission(ctx, &authz_es.UserMembershipRepo{Queries: q}, config.SystemAuthZ.RolePermissionMappings, config.InternalAuthZ.RolePermissionMappings, permission, orgID, resourceID)
			}
		},
		0,
		config.SystemAPIUsers,
		false,
	)
	logging.OnError(ctx, err).Fatal("unable to start queries")

	authZRepo, err := authz.Start(queries, es, client, keys.OIDC, config.ExternalSecure)
	logging.OnError(ctx, err).Fatal("unable to start authz repo")

	webAuthNConfig := &webauthn.Config{
		DisplayName:    config.WebAuthNName,
		ExternalSecure: config.ExternalSecure,
	}
	commands, err := command.StartCommands(ctx,
		es,
		cacheConnectors,
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
		keys.Target,
		&http.Client{},
		func(ctx context.Context, permission, orgID, resourceID string) (err error) {
			return internal_authz.CheckPermission(ctx, authZRepo, config.SystemAuthZ.RolePermissionMappings, config.InternalAuthZ.RolePermissionMappings, permission, orgID, resourceID)
		},
		sessionTokenVerifier,
		config.OIDC.DefaultAccessTokenLifetime,
		config.OIDC.DefaultRefreshTokenExpiration,
		config.OIDC.DefaultRefreshTokenIdleExpiration,
		config.DefaultInstance.SecretGenerators,
		config.Login.DefaultPaths,
	)
	logging.OnError(ctx, err).Fatal("unable to start commands")

	projection.CreateAll(ctx, client, es, config.Projections, keys.OIDC, keys.SAML)

	i18n.MustLoadSupportedLanguagesFromDir()

	notification.Register(
		ctx,
		config.Projections.Customizations["notifications"],
		config.Projections.Customizations["notificationsquotas"],
		config.Projections.Customizations["backchannel"],
		config.Projections.Customizations["telemetry"],
		config.Notifications,
		config.OIDC.BackChannelLogoutConfig(),
		*config.Telemetry,
		config.ExternalDomain,
		config.ExternalPort,
		config.ExternalSecure,
		commands,
		queries,
		es,
		config.Login.DefaultPaths.DefaultOTPEmailURLTemplate,
		config.SystemDefaults.Notifications.FileSystemPath,
		keys.User,
		keys.SMTP,
		keys.SMS,
		nil,
	)

	config.Auth.Spooler.Client = client
	config.Auth.Spooler.Eventstore = es
	authView, err := auth_view.StartView(config.Auth.Spooler.Client, keys.OIDC, queries, config.Auth.Spooler.Eventstore)
	logging.OnError(ctx, err).Fatal("unable to start auth view")
	auth_handler.Register(ctx, config.Auth.Spooler, authView, queries)

	config.Admin.Spooler.Client = client
	config.Admin.Spooler.Eventstore = es
	adminView, err := admin_view.StartView(config.Admin.Spooler.Client)
	logging.OnError(ctx, err).Fatal("unable to start admin view")

	admin_handler.Register(ctx, config.Admin.Spooler, adminView, staticStorage)

	instances := make(chan string, config.Projections.ConcurrentInstances)
	failedInstances := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(int(config.Projections.ConcurrentInstances))

	go func() {
		for instance := range failedInstances {
			logging.WithError(ctx, errors.New("projection failed for instance")).Error("projection failed", "instance", instance)
		}
	}()

	for range int(config.Projections.ConcurrentInstances) {
		go execProjections(ctx, instances, failedInstances, &wg)
	}

	existingInstances := queryInstanceIDs(ctx, client)
	for i, instance := range existingInstances {
		instances <- instance
		logging.Info(ctx, "instance queued for projection", "instance", instance, "index", fmt.Sprintf("%d/%d", i, len(existingInstances)))
	}
	close(instances)
	wg.Wait()

	close(failedInstances)

	logging.Info(ctx, "projections executed", "took", time.Since(start))
	return nil
}

func execProjections(ctx context.Context, instances <-chan string, failedInstances chan<- string, wg *sync.WaitGroup) {
	for instance := range instances {
		ctx = internal_authz.WithInstanceID(ctx, instance)
		logging.Info(ctx, "starting projections")

		err := projection.ProjectInstance(ctx)
		if err != nil {
			logging.WithError(ctx, err).Error("trigger failed")
			failedInstances <- instance
			continue
		}

		err = projection.ProjectInstanceFields(ctx)
		if err != nil {
			logging.WithError(ctx, err).Error("trigger fields failed")
			failedInstances <- instance
			continue
		}

		err = admin_handler.ProjectInstance(ctx)
		if err != nil {
			logging.WithError(ctx, err).Error("trigger admin handler failed")
			failedInstances <- instance
			continue
		}

		err = projection.ProjectInstanceFields(ctx)
		if err != nil {
			logging.WithError(ctx, err).Error("trigger fields failed")
			failedInstances <- instance
			continue
		}

		err = auth_handler.ProjectInstance(ctx)
		if err != nil {
			logging.WithError(ctx, err).Error("trigger auth handler failed")
			failedInstances <- instance
			continue
		}

		err = notification.ProjectInstance(ctx)
		if err != nil {
			logging.WithError(ctx, err).Error("trigger notification failed")
			failedInstances <- instance
			continue
		}

		logging.Info(ctx, "projections done")
	}
	wg.Done()
}

// queryInstanceIDs returns the instance configured by flag
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
	logging.OnError(ctx, err).Fatal("unable to query instances")
	return instances
}
