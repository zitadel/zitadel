package handlers

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/query"
)

type NotificationQueries struct {
	*query.Queries
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
	baseQueries *query.Queries,
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
