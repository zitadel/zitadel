package query

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/cache/connector"
	sd "github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	es_v4 "github.com/zitadel/zitadel/internal/v2/eventstore"
)

type Queries struct {
	eventstore   *eventstore.Eventstore
	eventStoreV4 es_v4.Querier
	client       *database.DB
	caches       *Caches

	keyEncryptionAlgorithm    crypto.EncryptionAlgorithm
	idpConfigEncryption       crypto.EncryptionAlgorithm
	targetEncryptionAlgorithm crypto.EncryptionAlgorithm
	smtpEncryptionAlgorithm   crypto.EncryptionAlgorithm
	smsEncryptionAlgorithm    crypto.EncryptionAlgorithm
	sessionTokenVerifier      func(ctx context.Context, sessionToken string, sessionID string, tokenID string) (err error)
	checkPermission           domain.PermissionCheck

	DefaultLanguage                     language.Tag
	mutex                               sync.Mutex
	LoginTranslationFileContents        map[string][]byte
	NotificationTranslationFileContents map[string][]byte
	supportedLangs                      []language.Tag
	zitadelRoles                        []authz.RoleMapping
	multifactors                        domain.MultifactorConfigs
	defaultAuditLogRetention            time.Duration
}

func StartQueries(
	ctx context.Context,
	es *eventstore.Eventstore,
	esV4 es_v4.Querier,
	querySqlClient, projectionSqlClient *database.DB,
	cacheConnectors connector.Connectors,
	projections projection.Config,
	defaults sd.SystemDefaults,
	idpConfigEncryption, otpEncryption, keyEncryptionAlgorithm, certEncryptionAlgorithm, targetEncryptionAlgorithm, smsEncryptionAlgorithm, smtpEncryptionAlgorithm crypto.EncryptionAlgorithm,
	zitadelRoles []authz.RoleMapping,
	sessionTokenVerifier func(ctx context.Context, sessionToken string, sessionID string, tokenID string) (err error),
	permissionCheck func(q *Queries) domain.PermissionCheck,
	defaultAuditLogRetention time.Duration,
	systemAPIUsers map[string]*authz.SystemAPIUser,
	startProjections bool,
) (repo *Queries, err error) {
	repo = &Queries{
		eventstore:                          es,
		eventStoreV4:                        esV4,
		client:                              querySqlClient,
		DefaultLanguage:                     language.Und,
		LoginTranslationFileContents:        make(map[string][]byte),
		NotificationTranslationFileContents: make(map[string][]byte),
		zitadelRoles:                        zitadelRoles,
		keyEncryptionAlgorithm:              keyEncryptionAlgorithm,
		idpConfigEncryption:                 idpConfigEncryption,
		targetEncryptionAlgorithm:           targetEncryptionAlgorithm,
		smsEncryptionAlgorithm:              smsEncryptionAlgorithm,
		smtpEncryptionAlgorithm:             smtpEncryptionAlgorithm,
		sessionTokenVerifier:                sessionTokenVerifier,
		multifactors: domain.MultifactorConfigs{
			OTP: domain.OTPConfig{
				CryptoMFA: otpEncryption,
				Issuer:    defaults.Multifactors.OTP.Issuer,
			},
		},
		defaultAuditLogRetention: defaultAuditLogRetention,
	}

	repo.checkPermission = permissionCheck(repo)

	projections.ActiveInstancer = repo
	err = projection.Create(ctx, projectionSqlClient, es, projections, keyEncryptionAlgorithm, certEncryptionAlgorithm, systemAPIUsers)
	if err != nil {
		return nil, err
	}
	if startProjections {
		err = projection.Start(ctx)
		if err != nil {
			return nil, err
		}
	}

	repo.caches, err = startCaches(
		ctx,
		cacheConnectors,
		ActiveInstanceConfig{
			MaxEntries: int(projections.MaxActiveInstances),
			TTL:        projections.HandleActiveInstances,
		},
	)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (q *Queries) Health(ctx context.Context) error {
	return q.client.Ping()
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

// triggerBatch calls Trigger on every handler in a separate Go routine.
// The returned context is the context returned by the Trigger that finishes last.
func triggerBatch(ctx context.Context, handlers ...*handler.Handler) {
	var wg sync.WaitGroup
	wg.Add(len(handlers))

	for _, h := range handlers {
		go func(ctx context.Context, h *handler.Handler) {
			name := h.ProjectionName()
			_, traceSpan := tracing.NewNamedSpan(ctx, fmt.Sprintf("Trigger%s", name))
			_, err := h.Trigger(ctx, handler.WithAwaitRunning())
			logging.OnError(err).WithField("projection", name).Debug("trigger failed")
			traceSpan.EndWithError(err)

			wg.Done()
		}(ctx, h)
	}

	wg.Wait()
}

func findTextEqualsQuery(column Column, queries []SearchQuery) (text string, ok bool) {
	for _, query := range queries {
		if query.Col() != column {
			continue
		}
		tq, ok := query.(*textQuery)
		if ok && tq.Compare == TextEquals {
			return tq.Text, true
		}
	}
	return "", false
}
