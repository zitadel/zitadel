package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/command"
)

func TestFirstInstance_validatePassword(t *testing.T) {
	policy := func(minLen uint64, lower, upper, number, symbol bool) command.InstanceSetup {
		return command.InstanceSetup{
			PasswordComplexityPolicy: struct {
				MinLength    uint64
				HasLowercase bool
				HasUppercase bool
				HasNumber    bool
				HasSymbol    bool
			}{
				MinLength:    minLen,
				HasLowercase: lower,
				HasUppercase: upper,
				HasNumber:    number,
				HasSymbol:    symbol,
			},
		}
	}

	tests := []struct {
		name     string
		password string
		setup    command.InstanceSetup
		wantErr  string
	}{
		{
			name:     "valid password passes all checks",
			password: "Password1!",
			setup:    policy(8, true, true, true, true),
		},
		{
			name:     "empty password is allowed (handled elsewhere)",
			password: "",
			setup:    policy(8, true, true, true, true),
		},
		{
			name:     "too short",
			password: "Pa1!",
			setup:    policy(8, true, true, true, true),
			wantErr:  "at least 8 characters",
		},
		{
			name:     "missing lowercase",
			password: "ALLUPPERCASE1!",
			setup:    policy(8, true, true, true, true),
			wantErr:  "lowercase letter",
		},
		{
			name:     "missing uppercase",
			password: "alllowercase1!",
			setup:    policy(8, true, true, true, true),
			wantErr:  "uppercase letter",
		},
		{
			name:     "missing number",
			password: "NoNumbers!!",
			setup:    policy(8, true, true, true, true),
			wantErr:  "one number",
		},
		{
			name:     "missing symbol",
			password: "NoSymbols1a",
			setup:    policy(8, true, true, true, true),
			wantErr:  "one symbol",
		},
		{
			name:     "borawak from issue 11651",
			password: "borawak",
			setup:    policy(8, true, true, true, true),
			wantErr:  "at least 8 characters",
		},
		{
			name:     "policy with no requirements accepts anything",
			password: "x",
			setup:    policy(0, false, false, false, false),
		},
		{
			name:     "error references env var name",
			password: "short",
			setup:    policy(8, false, false, false, false),
			wantErr:  FirstInstancePasswordEnvVar,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mig := &FirstInstance{
				Org: command.InstanceOrgSetup{
					Human: &command.AddHuman{
						Password: tt.password,
					},
				},
				instanceSetup: tt.setup,
			}

			err := mig.validatePassword()

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestFirstInstance_validatePassword_nilHuman(t *testing.T) {
	mig := &FirstInstance{
		Org: command.InstanceOrgSetup{
			Human: nil,
		},
	}
	assert.NoError(t, mig.validatePassword())
}
