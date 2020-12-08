package handler

import (
	"net/http"
	"time"

	"github.com/caos/logging"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/i18n"
	iam_es "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/view"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
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
}

type EventstoreRepos struct {
	UserEvents *usr_event.UserEventstore
	OrgEvents  *org_event.OrgEventstore
	IAMEvents  *iam_es.IAMEventstore
}

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, eventstore eventstore.Eventstore, repos EventstoreRepos, systemDefaults sd.SystemDefaults, i18n *i18n.Translator, dir http.FileSystem) []query.Handler {
	aesCrypto, err := crypto.NewAESCrypto(systemDefaults.UserVerificationKey)
	if err != nil {
		logging.Log("HANDL-s90ew").WithError(err).Debug("error create new aes crypto")
	}
	return []query.Handler{
		&NotifyUser{
			handler:   handler{view, bulkLimit, configs.cycleDuration("User"), errorCount},
			orgEvents: repos.OrgEvents,
			iamEvents: repos.IAMEvents,
			iamID:     systemDefaults.IamID,
		},
		&Notification{
			handler:        handler{view, bulkLimit, configs.cycleDuration("Notification"), errorCount},
			eventstore:     eventstore,
			userEvents:     repos.UserEvents,
			systemDefaults: systemDefaults,
			AesCrypto:      aesCrypto,
			i18n:           i18n,
			statikDir:      dir,
		},
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 2 * time.Second
	}
	return c.MinimumCycleDuration.Duration
}

func (h *handler) MinimumCycleDuration() time.Duration {
	return h.cycleDuration
}

func (h *handler) QueryLimit() uint64 {
	return h.bulkLimit
}
