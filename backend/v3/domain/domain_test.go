package domain_test

import (
	"os"
	"testing"
	"time"

	"github.com/zitadel/passwap"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		domain.SetPasswordHasher(&crypto.Hasher{Swapper: &passwap.Swapper{}})
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

func getPasskeyCondition(repo *domainmock.MockHumanUserRepository, pkeyID string) database.Condition {
	idCondition := getTextCondition("zitadel.users", "passkey_id", pkeyID)

	repo.EXPECT().
		PasskeyIDCondition(pkeyID).
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

func getSessionPasskeyFactorChange(repo *domainmock.MockSessionRepository, lastVerified time.Time, userVerified bool) database.Change {
	sessionPasskeyFactor := &domain.SessionFactorPasskey{
		LastVerifiedAt: lastVerified,
		UserVerified:   userVerified,
	}
	factorChange := make(database.Changes, 2)
	factorChange[0] = database.NewChange(
		database.NewColumn("zitadel.sessions", "password_factor_last_verified"),
		lastVerified,
	)
	factorChange[1] = database.NewChange(
		database.NewColumn("zitadel.sessions", "passkey_factor_user_verified"),
		userVerified,
	)

	repo.EXPECT().
		SetFactor(sessionPasskeyFactor).
		AnyTimes().
		Return(factorChange)
	return factorChange
}

func getHumanPasskeySignCount(repo *domainmock.MockHumanUserRepository, signCount uint32) database.Change {
	signCountChange := database.NewChange(
		database.NewColumn("zitadel.humans", "pkey_sign_count"),
		signCount,
	)
	repo.EXPECT().
		SetPasskeySignCount(signCount).
		AnyTimes().
		Return(signCountChange)

	return signCountChange
}
