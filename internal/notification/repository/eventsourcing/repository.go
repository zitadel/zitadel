package eventsourcing

import (
	"database/sql"
	"net/http"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	es_spol "github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
	"github.com/zitadel/zitadel/internal/notification/repository/eventsourcing/spooler"
	noti_view "github.com/zitadel/zitadel/internal/notification/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/query"
)

type Config struct {
	Spooler spooler.SpoolerConfig
}

type EsRepository struct {
	spooler *es_spol.Spooler
}

func Start(conf Config,
	dir http.FileSystem,
	externalPort uint16,
	externalSecure bool,
	command *command.Commands,
	queries *query.Queries,
	dbClient *sql.DB,
	assetsPrefix,
	fileSystemPath string,
	userEncryption crypto.EncryptionAlgorithm,
	smtpEncryption crypto.EncryptionAlgorithm,
	smsEncryption crypto.EncryptionAlgorithm,
) (*EsRepository, error) {
	es, err := v1.Start(dbClient)
	if err != nil {
		return nil, err
	}

	view, err := noti_view.StartView(dbClient)
	if err != nil {
		return nil, err
	}

	spool := spooler.StartSpooler(conf.Spooler, es, view, dbClient, command, queries, externalPort, externalSecure, dir, assetsPrefix, fileSystemPath, userEncryption, smtpEncryption, smsEncryption)

	return &EsRepository{
		spool,
	}, nil
}

func (repo *EsRepository) Health() error {
	return nil
}
