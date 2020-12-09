package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/label"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/login"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/org_iam"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_age"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_complexity"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_lockout"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(SetupStartedEventType, SetupStepMapper).
		RegisterFilterEventMapper(SetupDoneEventType, SetupStepMapper).
		RegisterFilterEventMapper(GlobalOrgSetEventType, GlobalOrgSetMapper).
		RegisterFilterEventMapper(ProjectSetEventType, ProjectSetMapper).
		RegisterFilterEventMapper(label.LabelPolicyAddedEventType, label.LabelPolicyAddedEventMapper).
		RegisterFilterEventMapper(label.LabelPolicyChangedEventType, label.LabelPolicyChangedEventMapper).
		RegisterFilterEventMapper(login.LoginPolicyAddedEventType, login.LoginPolicyAddedEventMapper).
		RegisterFilterEventMapper(login.LoginPolicyChangedEventType, login.LoginPolicyChangedEventMapper).
		RegisterFilterEventMapper(org_iam.OrgIAMPolicyAddedEventType, org_iam.OrgIAMPolicyAddedEventMapper).
		RegisterFilterEventMapper(password_age.PasswordAgePolicyAddedEventType, password_age.PasswordAgePolicyAddedEventMapper).
		RegisterFilterEventMapper(password_age.PasswordAgePolicyChangedEventType, password_age.PasswordAgePolicyChangedEventMapper).
		RegisterFilterEventMapper(password_complexity.PasswordComplexityPolicyAddedEventType, password_complexity.PasswordComplexityPolicyAddedEventMapper).
		RegisterFilterEventMapper(password_complexity.PasswordComplexityPolicyChangedEventType, password_complexity.PasswordComplexityPolicyChangedEventMapper).
		RegisterFilterEventMapper(password_lockout.PasswordLockoutPolicyAddedEventType, password_lockout.PasswordLockoutPolicyAddedEventMapper).
		RegisterFilterEventMapper(password_lockout.PasswordLockoutPolicyChangedEventType, password_lockout.PasswordLockoutPolicyChangedEventMapper).
		RegisterFilterEventMapper(MemberAddedEventType, MemberAddedEventMapper).
		RegisterFilterEventMapper(MemberChangedEventType, MemberChangedEventMapper).
		RegisterFilterEventMapper(MemberRemovedEventType, MemberRemovedEventMapper).
		RegisterFilterEventMapper(IDPConfigAddedEventType, IDPConfigAddedEventMapper).
		RegisterFilterEventMapper(IDPConfigChangedEventType, IDPConfigChangedEventMapper).
		RegisterFilterEventMapper(IDPConfigRemovedEventType, IDPConfigRemovedEventMapper).
		RegisterFilterEventMapper(IDPConfigDeactivatedEventType, IDPConfigDeactivatedEventMapper).
		RegisterFilterEventMapper(IDPConfigReactivatedEventType, IDPConfigReactivatedEventMapper).
		RegisterFilterEventMapper(IDPOIDCConfigAddedEventType, IDPOIDCConfigAddedEventMapper).
		RegisterFilterEventMapper(IDPOIDCConfigChangedEventType, IDPOIDCConfigChangedEventMapper)
}
