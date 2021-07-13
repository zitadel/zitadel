package handler

import (
	"net/http"
	"time"

	"github.com/caos/zitadel/internal/command"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/view"
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

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, es v1.Eventstore, command *command.Commands, systemDefaults sd.SystemDefaults, i18n *i18n.Translator, dir http.FileSystem, apiDomain string) []query.Handler {
	// aesCrypto, err := crypto.NewAESCrypto(systemDefaults.UserVerificationKey)
	// if err != nil {
	// 	logging.Log("HANDL-s90ew").WithError(err).Debug("error create new aes crypto")
	// }
	return []query.Handler{
		// newNotifyUser(
		// 	handler{view, bulkLimit, configs.cycleDuration("User"), errorCount, es},
		// 	systemDefaults.IamID,
		// ),
		// newNotification(
		// 	handler{view, bulkLimit, configs.cycleDuration("Notification"), errorCount, es},
		// 	command,
		// 	systemDefaults,
		// 	aesCrypto,
		// 	i18n,
		// 	dir,
		// 	apiDomain,
		// ),
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 3 * time.Minute
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
