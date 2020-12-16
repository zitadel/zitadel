package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_lockout"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(SetupStartedEventType, SetupStepMapper).
		RegisterFilterEventMapper(SetupDoneEventType, SetupStepMapper).
		RegisterFilterEventMapper(GlobalOrgSetEventType, GlobalOrgSetMapper).
		RegisterFilterEventMapper(ProjectSetEventType, ProjectSetMapper).
		RegisterFilterEventMapper(LabelPolicyAddedEventType, LabelPolicyAddedEventMapper).
		RegisterFilterEventMapper(LabelPolicyChangedEventType, LabelPolicyChangedEventMapper).
		RegisterFilterEventMapper(LoginPolicyAddedEventType, LoginPolicyAddedEventMapper).
		RegisterFilterEventMapper(LoginPolicyChangedEventType, LoginPolicyChangedEventMapper).
		RegisterFilterEventMapper(OrgIAMPolicyAddedEventType, OrgIAMPolicyAddedEventMapper).
		RegisterFilterEventMapper(PasswordAgePolicyAddedEventType, PasswordAgePolicyAddedEventMapper).
		RegisterFilterEventMapper(PasswordAgePolicyChangedEventType, PasswordAgePolicyChangedEventMapper).
		RegisterFilterEventMapper(PasswordComplexityPolicyAddedEventType, PasswordComplexityPolicyAddedEventMapper).
		RegisterFilterEventMapper(PasswordComplexityPolicyChangedEventType, PasswordComplexityPolicyChangedEventMapper).
		RegisterFilterEventMapper(password_lockout.PasswordLockoutPolicyAddedEventType, password_lockout.AddedEventMapper).
		RegisterFilterEventMapper(password_lockout.PasswordLockoutPolicyChangedEventType, password_lockout.ChangedEventMapper).
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
