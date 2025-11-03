package serviceping

type ReportType uint

const (
	ReportTypeBaseInformation ReportType = iota
	ReportTypeResourceCounts
)

type ServicePingReport struct {
	ReportID   string
	ReportType ReportType
}

func (r *ServicePingReport) Kind() string {
	return "service_ping_report"
}

// The following constants define the resource types for which counts are reported.
const (
	ResourceCountInstance                   = "instance"
	ResourceCountOrganization               = "organization"
	ResourceCountProject                    = "project"
	ResourceCountUser                       = "user"
	ResourceCountUserMachine                = "user_machine"
	ResourceCountIAMAdmin                   = "iam_admin"
	ResourceCountIdentityProvider           = "identity_provider"
	ResourceCountIdentityProviderLDAP       = "identity_provider_ldap"
	ResourceCountActionV1                   = "action_v1"
	ResourceCountActionExecution            = "execution"
	ResourceCountActionExecutionTarget      = "execution_target"
	ResourceCountLoginPolicy                = "login_policy"
	ResourceCountPasswordComplexityPolicy   = "password_complexity_policy"
	ResourceCountPasswordExpiryPolicy       = "password_expiry_policy"
	ResourceCountLockoutPolicy              = "lockout_policy"
	ResourceCountEnforceMFA                 = "enforce_mfa"
	ResourceCountPasswordChangeNotification = "password_change_notification"
	ResourceCountScimProvisionedUser        = "scim_provisioned_user"
)
