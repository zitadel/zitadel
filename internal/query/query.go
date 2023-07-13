package query

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	sd "github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/action"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	iam_repo "github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/keypair"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/session"
	usr_repo "github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
)

type Queries struct {
	eventstore *eventstore.Eventstore
	client     *database.DB

	idpConfigEncryption  crypto.EncryptionAlgorithm
	sessionTokenVerifier func(ctx context.Context, sessionToken string, sessionID string, tokenID string) (err error)
	checkPermission      domain.PermissionCheck

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

func StartQueries(
	ctx context.Context,
	es *eventstore.Eventstore,
	sqlClient *database.DB,
	projections projection.Config,
	defaults sd.SystemDefaults,
	idpConfigEncryption, otpEncryption, keyEncryptionAlgorithm, certEncryptionAlgorithm crypto.EncryptionAlgorithm,
	zitadelRoles []authz.RoleMapping,
	sessionTokenVerifier func(ctx context.Context, sessionToken string, sessionID string, tokenID string) (err error),
	permissionCheck func(q *Queries) domain.PermissionCheck,
) (repo *Queries, err error) {
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
		sessionTokenVerifier:                sessionTokenVerifier,
	}
	iam_repo.RegisterEventMappers(repo.eventstore)
	usr_repo.RegisterEventMappers(repo.eventstore)
	org.RegisterEventMappers(repo.eventstore)
	project.RegisterEventMappers(repo.eventstore)
	action.RegisterEventMappers(repo.eventstore)
	keypair.RegisterEventMappers(repo.eventstore)
	usergrant.RegisterEventMappers(repo.eventstore)
	session.RegisterEventMappers(repo.eventstore)
	idpintent.RegisterEventMappers(repo.eventstore)
	authrequest.RegisterEventMappers(repo.eventstore)
	oidcsession.RegisterEventMappers(repo.eventstore)

	repo.idpConfigEncryption = idpConfigEncryption
	repo.multifactors = domain.MultifactorConfigs{
		OTP: domain.OTPConfig{
			CryptoMFA: otpEncryption,
			Issuer:    defaults.Multifactors.OTP.Issuer,
		},
	}

	repo.checkPermission = permissionCheck(repo)

	err = projection.Create(ctx, sqlClient, es, projections, keyEncryptionAlgorithm, certEncryptionAlgorithm)
	if err != nil {
		return nil, err
	}
	projection.Start()

	return repo, nil
}

func (q *Queries) Health(ctx context.Context) error {
	return q.client.Ping()
}

type prepareDatabase interface {
	Timetravel(d time.Duration) string
}

// cleanStaticQueries removes whitespaces,
// such as ` `, \t, \n, from queries to improve
// readability in logs and errors.
func cleanStaticQueries(qs ...*string) {
	regex := regexp.MustCompile(`\s+`)
	for _, q := range qs {
		*q = regex.ReplaceAllString(*q, " ")
	}
}

func init() {
	cleanStaticQueries(
		&authRequestByIDQuery,
	)
}
