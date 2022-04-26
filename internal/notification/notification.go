package notification

import (
	"database/sql"

	"github.com/caos/logging"
	"github.com/rakyll/statik/fs"

	"github.com/zitadel/zitadel/internal/crypto"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/notification/repository/eventsourcing"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/query"
)

type Config struct {
	Repository eventsourcing.Config
}

func Start(config Config,
	externalPort uint16,
	externalSecure bool,
	command *command.Commands,
	queries *query.Queries,
	dbClient *sql.DB,
	assetsPrefix string,
	userEncryption crypto.EncryptionAlgorithm,
	smtpEncryption crypto.EncryptionAlgorithm,
	smsEncryption crypto.EncryptionAlgorithm,
) {
	statikFS, err := fs.NewWithNamespace("notification")
	logging.OnError(err).Panic("unable to start listener")

	_, err = eventsourcing.Start(config.Repository, statikFS, externalPort, externalSecure, command, queries, dbClient, assetsPrefix, userEncryption, smtpEncryption, smsEncryption)
	logging.OnError(err).Panic("unable to start app")
}
