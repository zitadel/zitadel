package query

import (
	"context"
	"database/sql"
	"net/http"
	"sync"

	"github.com/caos/logging"
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
	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"
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
}

type Config struct {
	Eventstore types.SQLUser
}

func StartQueries(ctx context.Context, es *eventstore.Eventstore, projections projection.Config, defaults sd.SystemDefaults, keyChan chan<- interface{}) (repo *Queries, err error) {
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

//
//func (r *Queries) IAMByID(ctx context.Context, id string) (_ *iam_model.IAM, err error) {
//	readModel, err := r.iamByID(ctx, id)
//	if err != nil {
//		return nil, err
//	}
//
//	return readModelToIAM(readModel), nil
//}
//
//func (r *Queries) iamByID(ctx context.Context, id string) (_ *ReadModel, err error) {
//	ctx, span := tracing.NewSpan(ctx)
//	defer func() { span.EndWithError(err) }()
//
//	readModel := NewReadModel(id)
//	err = r.eventstore.FilterToQueryReducer(ctx, readModel)
//	if err != nil {
//		return nil, err
//	}
//
//	return readModel, nil
//}
