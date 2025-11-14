package repository

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func (u userHuman) unqualifiedPasskeysTableName() string {
	return "human_passkeys"
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// AddPasskey implements [domain.HumanUserRepository.AddPasskey].
func (u userHuman) AddPasskey(passkey *domain.Passkey) database.Change {
	panic("unimplemented")
}

// RemovePasskey implements [domain.HumanUserRepository.RemovePasskey].
func (u userHuman) RemovePasskey(condition database.Condition) database.Change {
	return database.NewCTEChange(
		nil,
		func(builder *database.StatementBuilder) {
			builder.WriteString("DELETE FROM ")
			builder.WriteString(u.unqualifiedPasskeysTableName())
			builder.WriteString(" USING ")
			builder.WriteString(u.existingUserCTEName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(u.existingUserInstanceIDColumn(), u.passkeyInstanceIDColumn()),
				database.NewColumnCondition(u.existingUserIDColumn(), u.passkeyUserIDColumn()),
				condition,
			))
		},
	)
}

// SetPasskeyAttestationType implements [domain.HumanUserRepository.SetPasskeyAttestationType].
func (u userHuman) SetPasskeyAttestationType(attestationType string) database.Change {
	panic("unimplemented")
}

// SetPasskeyAuthenticatorAttestationGUID implements [domain.HumanUserRepository.SetPasskeyAuthenticatorAttestationGUID].
func (u userHuman) SetPasskeyAuthenticatorAttestationGUID(aaguid []byte) database.Change {
	panic("unimplemented")
}

// SetPasskeyKeyID implements [domain.HumanUserRepository.SetPasskeyKeyID].
func (u userHuman) SetPasskeyKeyID(keyID []byte) database.Change {
	panic("unimplemented")
}

// SetPasskeyName implements [domain.HumanUserRepository.SetPasskeyName].
func (u userHuman) SetPasskeyName(name string) database.Change {
	panic("unimplemented")
}

// SetPasskeyPublicKey implements [domain.HumanUserRepository.SetPasskeyPublicKey].
func (u userHuman) SetPasskeyPublicKey(publicKey []byte) database.Change {
	panic("unimplemented")
}

// SetPasskeySignCount implements [domain.HumanUserRepository.SetPasskeySignCount].
func (u userHuman) SetPasskeySignCount(signCount uint32) database.Change {
	panic("unimplemented")
}

// SetPasskeyUpdatedAt implements [domain.HumanUserRepository.SetPasskeyUpdatedAt].
func (u userHuman) SetPasskeyUpdatedAt(updatedAt time.Time) database.Change {
	panic("unimplemented")
}

// SetPasskeyVerifiedAt implements [domain.HumanUserRepository.SetPasskeyVerifiedAt].
func (u userHuman) SetPasskeyVerifiedAt(verifiedAt time.Time) database.Change {
	panic("unimplemented")
}

// UpdatePasskey implements [domain.HumanUserRepository.UpdatePasskey].
func (u userHuman) UpdatePasskey(condition database.Condition, changes ...database.Change) database.Change {
	panic("unimplemented")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// PasskeyChallengeCondition implements [domain.HumanUserRepository.PasskeyChallengeCondition].
func (u userHuman) PasskeyChallengeCondition(challenge string) database.Condition {
	// TODO: implement passkey challenge condition
	panic("unimplemented")
}

// PasskeyIDCondition implements [domain.HumanUserRepository.PasskeyIDCondition].
func (u userHuman) PasskeyIDCondition(passkeyID string) database.Condition {
	return database.NewTextCondition(u.passkeyIDColumn(), database.TextOperationEqual, passkeyID)
}

// PasskeyKeyIDCondition implements [domain.HumanUserRepository.PasskeyKeyIDCondition].
func (u userHuman) PasskeyKeyIDCondition(keyID string) database.Condition {
	return database.NewTextCondition(u.passkeyKeyIDColumn(), database.TextOperationEqual, keyID)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userHuman) passkeyIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "id")
}

func (u userHuman) passkeyKeyIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "key_id")
}

func (u userHuman) passkeyInstanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "instance_id")
}

func (u userHuman) passkeyUserIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "user_id")
}
