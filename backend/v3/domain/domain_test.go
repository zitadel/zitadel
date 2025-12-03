package domain_test

import (
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		return m.Run()
	}())
}

func getUserIDCondition(repo *domainmock.MockUserRepository, userID string) database.Condition {
	idCondition := getTextCondition("zitadel.users", "id", userID)

	repo.EXPECT().
		IDCondition(userID).
		AnyTimes().
		Return(idCondition)
	return idCondition
}

func getSessionIDCondition(repo *domainmock.MockSessionRepository, sessionID string) database.Condition {
	idCondition := getTextCondition("zitadel.sessions", "id", sessionID)

	repo.EXPECT().
		IDCondition(sessionID).
		AnyTimes().
		Return(idCondition)
	return idCondition
}

func getHumanUserIDCondition(repo *domainmock.MockHumanUserRepository, userID string) database.Condition {
	idCondition := getTextCondition("zitadel.users", "id", userID)

	repo.EXPECT().
		IDCondition(userID).
		AnyTimes().
		Return(idCondition)
	return idCondition
}

func getPasskeyTypeCondition(repo *domainmock.MockHumanUserRepository, passkeyType domain.PasskeyType) database.Condition {
	typeCondition := getNumberCondition("zitadel.human_users", "passkeys", int(passkeyType))

	repo.EXPECT().
		PasskeyTypeCondition(database.NumberOperationEqual, passkeyType).
		AnyTimes().
		Return(typeCondition)
	return typeCondition
}

func getTextCondition(tableName, column, value string) database.Condition {
	return database.NewTextCondition(
		database.NewColumn(tableName, column),
		database.TextOperationEqual,
		value,
	)
}

func getNumberCondition(tableName, column string, value int) database.Condition {
	return database.NewNumberCondition(
		database.NewColumn(tableName, column),
		database.NumberOperationEqual,
		value,
	)
}

func getSessionPasswordFactorChange(repo *domainmock.MockSessionRepository, lastVerified, lastFailed *time.Time) database.Change {
	sessionPasswordFactor := &domain.SessionFactorPassword{}
	var factorChange database.Change
	if lastVerified != nil {
		sessionPasswordFactor.LastVerifiedAt = *lastVerified
		factorChange = database.NewChange(
			database.NewColumn("zitadel.sessions", "password_factor_last_verified"),
			*lastVerified,
		)
	}
	if lastFailed != nil {
		sessionPasswordFactor.LastFailedAt = *lastFailed
		factorChange = database.NewChange(
			database.NewColumn("zitadel.sessions", "password_factor_last_failed"),
			*lastFailed,
		)
	}

	repo.EXPECT().
		SetFactor(sessionPasswordFactor).
		AnyTimes().
		Return(factorChange)
	return factorChange
}
