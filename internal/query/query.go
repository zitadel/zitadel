package query

import (
	"context"
	"database/sql"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"net/http"
	"sync"

	"github.com/rakyll/statik/fs"
	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	sd "github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/config/types"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/action"
	iam_repo "github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/keypair"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	usr_repo "github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
)

type Queries struct {
	iamID      string
	eventstore *eventstore.Eventstore
	client     *sql.DB

	idpConfigSecretCrypto crypto.EncryptionAlgorithm

	DefaultLanguage                     language.Tag
	LoginDir                            http.FileSystem
	NotificationDir                     http.FileSystem
	mutex                               sync.Mutex
	LoginTranslationFileContents        map[string][]byte
	NotificationTranslationFileContents map[string][]byte
	supportedLangs                      []language.Tag
	zitadelRoles                        []authz.RoleMapping
	multifactors                        domain.MultifactorConfigs
}

type Config struct {
	Eventstore types.SQLUser
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

	repo.idpConfigSecretCrypto, err = crypto.NewAESCrypto(defaults.IDPConfigVerificationKey)
	if err != nil {
		return nil, err
	}

	aesOTPCrypto, err := crypto.NewAESCrypto(defaults.Multifactors.OTP.VerificationKey)
	if err != nil {
		return nil, err
	}
	repo.multifactors = domain.MultifactorConfigs{
		OTP: domain.OTPConfig{
			CryptoMFA: aesOTPCrypto,
			Issuer:    defaults.Multifactors.OTP.Issuer,
		},
	}

	err = projection.Start(ctx, sqlClient, es, projections, defaults, keyChan)
	if err != nil {
		return nil, err
	}

	return repo, nil
}
