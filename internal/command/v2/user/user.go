package user

import "github.com/caos/zitadel/internal/domain"

func isUserStateExists(state domain.UserState) bool {
	return !hasUserState(state, domain.UserStateDeleted, domain.UserStateUnspecified)
}

func hasUserState(check domain.UserState, states ...domain.UserState) bool {
	for _, state := range states {
		if check == state {
			return true
		}
	}
	return false
}
