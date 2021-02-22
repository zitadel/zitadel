package handler

import (
	"time"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
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

	es eventstore.Eventstore
}

func (h *handler) Eventstore() eventstore.Eventstore {
	return h.es
}

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, es eventstore.Eventstore, defaults systemdefaults.SystemDefaults) []query.Handler {
	return []query.Handler{
		newProject(
			handler{view, bulkLimit, configs.cycleDuration("Project"), errorCount, es}),
		newProjectGrant(
			handler{view, bulkLimit, configs.cycleDuration("ProjectGrant"), errorCount, es}),
		newProjectRole(handler{view, bulkLimit, configs.cycleDuration("ProjectRole"), errorCount, es}),
		newProjectMember(handler{view, bulkLimit, configs.cycleDuration("ProjectMember"), errorCount, es}),
		newProjectGrantMember(handler{view, bulkLimit, configs.cycleDuration("ProjectGrantMember"), errorCount, es}),
		newApplication(handler{view, bulkLimit, configs.cycleDuration("Application"), errorCount, es}),
		newUser(handler{view, bulkLimit, configs.cycleDuration("User"), errorCount, es},
			defaults.IamID),
		newUserGrant(handler{view, bulkLimit, configs.cycleDuration("UserGrant"), errorCount, es}),
		newOrg(
			handler{view, bulkLimit, configs.cycleDuration("Org"), errorCount, es}),
		newOrgMember(
			handler{view, bulkLimit, configs.cycleDuration("OrgMember"), errorCount, es}),
		newOrgDomain(
			handler{view, bulkLimit, configs.cycleDuration("OrgDomain"), errorCount, es}),
		newUserMembership(
			handler{view, bulkLimit, configs.cycleDuration("UserMembership"), errorCount, es}),
		newAuthNKeys(
			handler{view, bulkLimit, configs.cycleDuration("MachineKeys"), errorCount, es}),
		newIDPConfig(
			handler{view, bulkLimit, configs.cycleDuration("IDPConfig"), errorCount, es}),
		newLoginPolicy(
			handler{view, bulkLimit, configs.cycleDuration("LoginPolicy"), errorCount, es}),
		newLabelPolicy(
			handler{view, bulkLimit, configs.cycleDuration("LabelPolicy"), errorCount, es}),
		newIDPProvider(
			handler{view, bulkLimit, configs.cycleDuration("IDPProvider"), errorCount, es},
			defaults),
		newExternalIDP(
			handler{view, bulkLimit, configs.cycleDuration("ExternalIDP"), errorCount, es},
			defaults),
		newPasswordComplexityPolicy(
			handler{view, bulkLimit, configs.cycleDuration("PasswordComplexityPolicy"), errorCount, es}),
		newPasswordAgePolicy(
			handler{view, bulkLimit, configs.cycleDuration("PasswordAgePolicy"), errorCount, es}),
		newPasswordLockoutPolicy(
			handler{view, bulkLimit, configs.cycleDuration("PasswordLockoutPolicy"), errorCount, es}),
		newOrgIAMPolicy(
			handler{view, bulkLimit, configs.cycleDuration("OrgIAMPolicy"), errorCount, es}),
		newMailTemplate(
			handler{view, bulkLimit, configs.cycleDuration("MailTemplate"), errorCount, es}),
		newMailText(
			handler{view, bulkLimit, configs.cycleDuration("MailText"), errorCount, es}),
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
