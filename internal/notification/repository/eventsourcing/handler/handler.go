package handler

import (
	"net/http"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/query"

	"github.com/zitadel/zitadel/internal/command"
	sd "github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/config/types"
	"github.com/zitadel/zitadel/internal/crypto"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	queryv1 "github.com/zitadel/zitadel/internal/eventstore/v1/query"
	"github.com/zitadel/zitadel/internal/notification/repository/eventsourcing/view"
)

type Configs map[string]*Config

type Config struct {
	MinimumCycleDuration types.Duration
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

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, es v1.Eventstore, command *command.Commands, queries *query.Queries, systemDefaults sd.SystemDefaults, dir http.FileSystem, apiDomain string) []queryv1.Handler {
	aesCrypto, err := crypto.NewAESCrypto(systemDefaults.UserVerificationKey)
	if err != nil {
		logging.Log("HANDL-s90ew").WithError(err).Debug("error create new aes crypto")
	}
	return []queryv1.Handler{
		newNotifyUser(
			handler{view, bulkLimit, configs.cycleDuration("User"), errorCount, es},
			systemDefaults.IamID,
			queries,
		),
		newNotification(
			handler{view, bulkLimit, configs.cycleDuration("Notification"), errorCount, es},
			command,
			queries,
			systemDefaults,
			aesCrypto,
			dir,
			apiDomain,
		),
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 1 * time.Minute
	}
	return c.MinimumCycleDuration.Duration
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
