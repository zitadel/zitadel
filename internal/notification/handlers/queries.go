package handlers

import (
	"context"
	"time"

	"github.com/go-jose/go-jose/v4"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
)

type Queries interface {
	ActiveLabelPolicyByOrg(ctx context.Context, orgID string, withOwnerRemoved bool) (*query.LabelPolicy, error)
	MailTemplateByOrg(ctx context.Context, orgID string, withOwnerRemoved bool) (*query.MailTemplate, error)
	GetNotifyUserByID(ctx context.Context, shouldTriggered bool, userID string) (*query.NotifyUser, error)
	CustomTextListByTemplate(ctx context.Context, aggregateID, template string, withOwnerRemoved bool) (*query.CustomTexts, error)
	SearchInstanceDomains(ctx context.Context, queries *query.InstanceDomainSearchQueries) (*query.InstanceDomains, error)
	SessionByID(ctx context.Context, shouldTriggerBulk bool, id, sessionToken string) (*query.Session, error)
	NotificationPolicyByOrg(ctx context.Context, shouldTriggerBulk bool, orgID string, withOwnerRemoved bool) (*query.NotificationPolicy, error)
	SearchMilestones(ctx context.Context, instanceIDs []string, queries *query.MilestonesSearchQueries) (*query.Milestones, error)
	NotificationProviderByIDAndType(ctx context.Context, aggID string, providerType domain.NotificationProviderType) (*query.DebugNotificationProvider, error)
	SMSProviderConfigActive(ctx context.Context, resourceOwner string) (config *query.SMSConfig, err error)
	SMTPConfigActive(ctx context.Context, resourceOwner string) (*query.SMTPConfig, error)
	GetDefaultLanguage(ctx context.Context) language.Tag
	GetInstanceRestrictions(ctx context.Context) (restrictions query.Restrictions, err error)
	InstanceByID(ctx context.Context, id string) (instance authz.Instance, err error)
	GetActiveSigningWebKey(ctx context.Context) (*jose.JSONWebKey, error)
	ActivePrivateSigningKey(ctx context.Context, t time.Time) (keys *query.PrivateKeys, err error)

	ActiveInstances() []string
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
	}
}
