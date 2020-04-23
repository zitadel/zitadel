package systemdefaults

import "github.com/caos/zitadel/internal/crypto"

type SystemDefaults struct {
	SecretGenerator SecretGenerator
}

type SecretGenerator struct {
	PasswordSaltCost      int
	ClientSecretGenerator crypto.GeneratorConfig
}
