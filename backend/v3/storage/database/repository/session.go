package repository

import "github.com/zitadel/zitadel/backend/v3/domain"

type session struct{}

func SessionRepository() domain.SessionRepository {
	return nil //new(session)
}
