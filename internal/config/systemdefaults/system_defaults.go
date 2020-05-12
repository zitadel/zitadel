package systemdefaults

import (
	"github.com/caos/zitadel/internal/crypto"
	pol "github.com/caos/zitadel/internal/policy"
)

type SystemDefaults struct {
	SecretGenerator SecretGenerator
	DefaultPolicies DefaultPolicies
}

type SecretGenerator struct {
	PasswordSaltCost      int
	ClientSecretGenerator crypto.GeneratorConfig
}

type DefaultPolicies struct {
	Age        pol.PasswordAgePolicyDefault
	Complexity pol.PasswordComplexityPolicyDefault
	Lockout    pol.PasswordLockoutPolicyDefault
}
