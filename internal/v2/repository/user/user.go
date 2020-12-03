package user

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	UserLocked      models.EventType = "user.locked"
	UserUnlocked    models.EventType = "user.unlocked"
	UserDeactivated models.EventType = "user.deactivated"
	UserReactivated models.EventType = "user.reactivated"
	UserRemoved     models.EventType = "user.removed"

	UserTokenAdded models.EventType = "user.token.added"
)
