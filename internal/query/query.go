package query

import (
	"context"
	"database/sql"
	"net/http"
	"sync"

	"github.com/caos/logging"
	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/api/authz"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/query/projection"
	"github.com/caos/zitadel/internal/repository/action"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/keypair"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/project"
	usr_repo "github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/repository/usergrant"
)

type Queries struct {
	iamID      string
	eventstore *eventstore.Eventstore
	client     *sql.DB

	DefaultLanguage                     language.Tag
	LoginDir                            http.FileSystem
	NotificationDir                     http.FileSystem
	mutex                               sync.Mutex
	LoginTranslationFileContents        map[string][]byte
	NotificationTranslationFileContents map[string][]byte
	supportedLangs                      []language.Tag
	zitadelRoles                        []authz.RoleMapping
}

type Config struct {
	Eventstore types.SQLUser
}
type ConfigV2 struct {
	Eventstore types.SQLUser2
}

//
//var DefaultConfig = ConfigV2{
//	Eventstore: types.SQLUser2{
//		User:            queryUser,
//		Password:        os.Getenv(keyQueriesPassword),
//		MaxOpenConns:    2,
//		MaxConnLifetime: types.Duration{Duration: 30 * time.Minute},
//		MaxConnIdleTime: types.Duration{Duration: 30 * time.Minute},
//		SSL: types.SSLUser{
//			Cert: os.Getenv(keyQueriesCert),
//			Key:  os.Getenv(keyQueriesKey),
//		},
//	},
//}

func StartQueries2(ctx context.Context, es *eventstore.Eventstore, sqlClient *sql.DB, projections projection.Config, defaults sd.SystemDefaults, keyChan chan<- interface{}, zitadelRoles []authz.RoleMapping) (repo *Queries, err error) {
	//statikLoginFS, err := fs.NewWithNamespace("login")
	//logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start login statik dir")
	//
	//statikNotificationFS, err := fs.NewWithNamespace("notification")
	//logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start notification statik dir")

	repo = &Queries{
		iamID:           defaults.IamID,
		eventstore:      es,
		client:          sqlClient,
		DefaultLanguage: defaults.DefaultLanguage,
		//LoginDir:                            statikLoginFS,
		//NotificationDir:                     statikNotificationFS,
		LoginTranslationFileContents:        make(map[string][]byte),
		NotificationTranslationFileContents: make(map[string][]byte),
		zitadelRoles:                        zitadelRoles,
	}
	iam_repo.RegisterEventMappers(repo.eventstore)
	usr_repo.RegisterEventMappers(repo.eventstore)
	org.RegisterEventMappers(repo.eventstore)
	project.RegisterEventMappers(repo.eventstore)
	action.RegisterEventMappers(repo.eventstore)
	keypair.RegisterEventMappers(repo.eventstore)
	usergrant.RegisterEventMappers(repo.eventstore)

	err = projection.Start(ctx, sqlClient, es, projections, defaults, keyChan)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func StartQueries(ctx context.Context, es *eventstore.Eventstore, projections projection.Config, defaults sd.SystemDefaults, keyChan chan<- interface{}, zitadelRoles []authz.RoleMapping) (repo *Queries, err error) {
	sqlClient, err := projections.CRDB.Start()
	if err != nil {
		return nil, err
	}

	statikLoginFS, err := fs.NewWithNamespace("login")
	logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start login statik dir")

	statikNotificationFS, err := fs.NewWithNamespace("notification")
	logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start notification statik dir")

	repo = &Queries{
		iamID:                               defaults.IamID,
		eventstore:                          es,
		client:                              sqlClient,
		DefaultLanguage:                     defaults.DefaultLanguage,
		LoginDir:                            statikLoginFS,
		NotificationDir:                     statikNotificationFS,
		LoginTranslationFileContents:        make(map[string][]byte),
		NotificationTranslationFileContents: make(map[string][]byte),
		zitadelRoles:                        zitadelRoles,
	}
	iam_repo.RegisterEventMappers(repo.eventstore)
	usr_repo.RegisterEventMappers(repo.eventstore)
	org.RegisterEventMappers(repo.eventstore)
	project.RegisterEventMappers(repo.eventstore)
	action.RegisterEventMappers(repo.eventstore)
	keypair.RegisterEventMappers(repo.eventstore)
	usergrant.RegisterEventMappers(repo.eventstore)

	err = projection.Start(ctx, sqlClient, es, projections, defaults, keyChan)
	if err != nil {
		return nil, err
	}

	return repo, nil
}
