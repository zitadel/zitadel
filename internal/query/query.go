package query

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sync"

	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/action"
	iam_repo "github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/keypair"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	usr_repo "github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
)

type Queries struct {
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

func StartQueries(ctx context.Context, es *eventstore.Eventstore, sqlClient *sql.DB, projections projection.Config, keyEncryptionAlgorithm crypto.EncryptionAlgorithm, zitadelRoles []authz.RoleMapping) (repo *Queries, err error) {
	statikLoginFS, err := fs.NewWithNamespace("login")
	if err != nil {
		return nil, fmt.Errorf("unable to start login statik dir")
	}

	statikNotificationFS, err := fs.NewWithNamespace("notification")
	if err != nil {
		return nil, fmt.Errorf("unable to start notification statik dir")
	}

	repo = &Queries{
		eventstore:                          es,
		client:                              sqlClient,
		DefaultLanguage:                     language.Und,
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

	err = projection.Start(ctx, sqlClient, es, projections, keyEncryptionAlgorithm)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (q *Queries) Health(ctx context.Context) error {
	return q.client.Ping()
}
