package handlers

import (
	"context"
	"net/http"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
)

type Queries interface {
	ActiveLabelPolicyByOrg(ctx context.Context, orgID string, withOwnerRemoved bool) (*query.LabelPolicy, error)
	MailTemplateByOrg(ctx context.Context, orgID string, withOwnerRemoved bool) (*query.MailTemplate, error)
	GetNotifyUserByID(ctx context.Context, shouldTriggered bool, userID string, queries ...query.SearchQuery) (*query.NotifyUser, error)
	CustomTextListByTemplate(ctx context.Context, aggregateID, template string, withOwnerRemoved bool) (*query.CustomTexts, error)
	SearchInstanceDomains(ctx context.Context, queries *query.InstanceDomainSearchQueries) (*query.InstanceDomains, error)
	SessionByID(ctx context.Context, shouldTriggerBulk bool, id, sessionToken string) (*query.Session, error)
	NotificationPolicyByOrg(ctx context.Context, shouldTriggerBulk bool, orgID string, withOwnerRemoved bool) (*query.NotificationPolicy, error)
	SearchMilestones(ctx context.Context, instanceIDs []string, queries *query.MilestonesSearchQueries) (*query.Milestones, error)
	NotificationProviderByIDAndType(ctx context.Context, aggID string, providerType domain.NotificationProviderType) (*query.DebugNotificationProvider, error)
	SMSProviderConfig(ctx context.Context, queries ...query.SearchQuery) (*query.SMSConfig, error)
	SMTPConfigByAggregateID(ctx context.Context, aggregateID string) (*query.SMTPConfig, error)
	GetDefaultLanguage(ctx context.Context) language.Tag
}

type NotificationQueries struct {
	Queries
	es                 *eventstore.Eventstore
	externalDomain     string
	externalPort       uint16
	externalSecure     bool
	fileSystemPath     string
	UserDataCrypto     crypto.EncryptionAlgorithm
	SMTPPasswordCrypto crypto.EncryptionAlgorithm
	SMSTokenCrypto     crypto.EncryptionAlgorithm
	statikDir          http.FileSystem
}

func NewNotificationQueries(
	baseQueries Queries,
	es *eventstore.Eventstore,
	externalDomain string,
	externalPort uint16,
	externalSecure bool,
	fileSystemPath string,
	userDataCrypto crypto.EncryptionAlgorithm,
	smtpPasswordCrypto crypto.EncryptionAlgorithm,
	smsTokenCrypto crypto.EncryptionAlgorithm,
	statikDir http.FileSystem,
) *NotificationQueries {
	return &NotificationQueries{
		Queries:            baseQueries,
		es:                 es,
		externalDomain:     externalDomain,
		externalPort:       externalPort,
		externalSecure:     externalSecure,
		fileSystemPath:     fileSystemPath,
		UserDataCrypto:     userDataCrypto,
		SMTPPasswordCrypto: smtpPasswordCrypto,
		SMSTokenCrypto:     smsTokenCrypto,
		statikDir:          statikDir,
	}
}
