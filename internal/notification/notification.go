package notification

import (
	"database/sql"

	"github.com/caos/logging"
	"github.com/rakyll/statik/fs"

	"github.com/caos/zitadel/internal/crypto"

	"github.com/caos/zitadel/internal/command"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing"
	_ "github.com/caos/zitadel/internal/notification/statik"
	"github.com/caos/zitadel/internal/query"
)

type Config struct {
	Repository eventsourcing.Config
}

func Start(config Config,
	systemDefaults sd.SystemDefaults,
	command *command.Commands,
	queries *query.Queries,
	dbClient *sql.DB,
	assetsPrefix string,
	userEncryption crypto.EncryptionAlgorithm,
	smtpEncryption crypto.EncryptionAlgorithm,
) {
	statikFS, err := fs.NewWithNamespace("notification")
	logging.OnError(err).Panic("unable to start listener")

	_, err = eventsourcing.Start(config.Repository, statikFS, systemDefaults, command, queries, dbClient, assetsPrefix, userEncryption, smtpEncryption)
	logging.OnError(err).Panic("unable to start app")
}
