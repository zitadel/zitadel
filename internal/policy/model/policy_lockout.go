package model

import "github.com/caos/zitadel/internal/eventstore/models"

type PasswordLockoutPolicy struct {
	models.ObjectRoot

	Description         string
	State               int32
	MaxAttempts         uint64
	ShowLockOutFailures bool
}

func (p *PasswordLockoutPolicy) IsValid() bool {
	if p.Description == "" {
		return false
	}
	return true
}
