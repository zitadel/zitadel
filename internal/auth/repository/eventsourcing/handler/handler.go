package handler

import (
	"github.com/caos/zitadel/internal/eventstore/v1"
	"time"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	key_model "github.com/caos/zitadel/internal/key/model"
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

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, es v1.Eventstore, systemDefaults sd.SystemDefaults, keyChan chan<- *key_model.KeyView) []query.Handler {
	return []query.Handler{
		newUser(
			handler{view, bulkLimit, configs.cycleDuration("User"), errorCount, es},
			systemDefaults.IamID),
		newUserSession(
			handler{view, bulkLimit, configs.cycleDuration("UserSession"), errorCount, es}),
		newUserMembership(
			handler{view, bulkLimit, configs.cycleDuration("UserMembership"), errorCount, es}),
		newToken(
			handler{view, bulkLimit, configs.cycleDuration("Token"), errorCount, es}),
		newKey(
			handler{view, bulkLimit, configs.cycleDuration("Key"), errorCount, es},
			keyChan),
		newApplication(handler{view, bulkLimit, configs.cycleDuration("Application"), errorCount, es}),
		newOrg(
			handler{view, bulkLimit, configs.cycleDuration("Org"), errorCount, es}),
		newUserGrant(
			handler{view, bulkLimit, configs.cycleDuration("UserGrant"), errorCount, es},
			systemDefaults.IamID),
		newAuthNKeys(
			handler{view, bulkLimit, configs.cycleDuration("MachineKey"), errorCount, es}),
		newLoginPolicy(
			handler{view, bulkLimit, configs.cycleDuration("LoginPolicy"), errorCount, es}),
		newIDPConfig(
			handler{view, bulkLimit, configs.cycleDuration("IDPConfig"), errorCount, es}),
		newIDPProvider(
			handler{view, bulkLimit, configs.cycleDuration("IDPProvider"), errorCount, es},
			systemDefaults),
		newExternalIDP(
			handler{view, bulkLimit, configs.cycleDuration("ExternalIDP"), errorCount, es},
			systemDefaults),
		newPasswordComplexityPolicy(
			handler{view, bulkLimit, configs.cycleDuration("PasswordComplexityPolicy"), errorCount, es}),
		newOrgIAMPolicy(
			handler{view, bulkLimit, configs.cycleDuration("OrgIAMPolicy"), errorCount, es}),
		newProjectRole(handler{view, bulkLimit, configs.cycleDuration("ProjectRole"), errorCount, es}),
		newLabelPolicy(handler{view, bulkLimit, configs.cycleDuration("LabelPolicy"), errorCount, es}),
		newFeatures(handler{view, bulkLimit, configs.cycleDuration("Features"), errorCount, es}),
		newRefreshToken(handler{view, bulkLimit, configs.cycleDuration("RefreshToken"), errorCount, es}),
		newPrivacyPolicy(handler{view, bulkLimit, configs.cycleDuration("PrivacyPolicy"), errorCount, es}),
		newCustomText(handler{view, bulkLimit, configs.cycleDuration("CustomTexts"), errorCount, es}),
		newPasswordLockoutPolicy(handler{view, bulkLimit, configs.cycleDuration("LockoutPolicy"), errorCount, es}),
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
