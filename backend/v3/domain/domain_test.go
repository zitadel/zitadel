package domain_test

import (
	"os"
	"testing"

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

func getTextCondition(tableName, column, value string) database.Condition {
	return database.NewTextCondition(
		database.NewColumn(tableName, column),
		database.TextOperationEqual,
		value,
	)
}
