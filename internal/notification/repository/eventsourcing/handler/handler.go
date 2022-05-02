package handler

import (
	"net/http"
	"time"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	queryv1 "github.com/zitadel/zitadel/internal/eventstore/v1/query"
	"github.com/zitadel/zitadel/internal/notification/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/query"
)

type Configs map[string]*Config

type Config struct {
	MinimumCycleDuration time.Duration
}

type handler struct {
	view                *view.View
	bulkLimit           uint64
	cycleDuration       time.Duration
	errorCountUntilSkip uint64

	es v1.Eventstore
}

func (h *handler) Eventstore() v1.Eventstore {
	return h.es
}

func Register(configs Configs,
	bulkLimit,
	errorCount uint64,
	view *view.View,
	es v1.Eventstore,
	command *command.Commands,
	queries *query.Queries,
	externalPort uint16,
	externalSecure bool,
	dir http.FileSystem,
	assetsPrefix,
	fileSystemPath string,
	userEncryption crypto.EncryptionAlgorithm,
	smtpEncryption crypto.EncryptionAlgorithm,
	smsEncryption crypto.EncryptionAlgorithm,
) []queryv1.Handler {
	return []queryv1.Handler{
		newNotifyUser(
			handler{view, bulkLimit, configs.cycleDuration("User"), errorCount, es},
			queries,
		),
		newNotification(
			handler{view, bulkLimit, configs.cycleDuration("Notification"), errorCount, es},
			command,
			queries,
			externalPort,
			externalSecure,
			dir,
			assetsPrefix,
			fileSystemPath,
			userEncryption,
			smtpEncryption,
			smsEncryption,
		),
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 1 * time.Minute
	}
	return c.MinimumCycleDuration
}

func (h *handler) MinimumCycleDuration() time.Duration {
	return h.cycleDuration
}

func (h *handler) LockDuration() time.Duration {
	return h.cycleDuration / 3
}

func (h *handler) QueryLimit() uint64 {
	return h.bulkLimit
}
