package policy

type PasswordAgePolicyDefault struct {
	Description    string
	MaxAgeDays     uint64
	ExpireWarnDays uint64
}

type PasswordComplexityPolicyDefault struct {
	Description  string
	MinLength    uint64
	HasLowercase bool
	HasUppercase bool
	HasNumber    bool
	HasSymbol    bool
}

type PasswordLockoutPolicyDefault struct {
	Description         string
	MaxAttempts         uint64
	ShowLockOutFailures bool
}
