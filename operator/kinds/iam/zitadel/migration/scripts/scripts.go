package scripts

func GetAll() map[string]string {
	return map[string]string{
		"V1.0__databases.sql":                V10Databases,
		"V1.1__eventstore.sql":               V11Eventstore,
		"V1.2__views.sql":                    V12Views,
		"V1.3__usermembership.sql":           V13Usermembership,
		"V1.4__compliance.sql":               V14Compliance,
		"V1.5__orgdomain_validationtype.sql": V15OrgDomainValidationType,
		"V1.6__origin_allow_list.sql":        V16OriginAllowList,
		"V1.7__idps.sql":                     V17IDPs,
		"V1.8__username_change.sql":          V18UsernameChange,
		"V1.9__token.sql":                    V19Token,
		"V1.10__user_machine_keys.sql":       V110UserMachineKeys,
		"V1.11__usermembership.sql":          V111UserMembership,
		"V1.12__machine_keys.sql":            V112MachineKeys,
		"V1.13__machine_keys_public.sql":     V113MachineKeysPublic,
		"V1.14__auth_loginpolicy.sql":        V114AuthLoginPolicy,
		"V1.15__idp_providers.sql":           V115IdpProviders,
	}
}
