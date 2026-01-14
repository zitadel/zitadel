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
	var (
		createdAt  any = database.NowInstruction
		updatedAt  any = database.NowInstruction
		verifiedAt any = database.NowInstruction
	)
	if !passkey.CreatedAt.IsZero() {
		createdAt = passkey.CreatedAt
	}
	if !passkey.UpdatedAt.IsZero() {
		updatedAt = passkey.UpdatedAt
	}
	if !passkey.VerifiedAt.IsZero() {
		verifiedAt = passkey.VerifiedAt
	}
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("INSERT INTO ")
			builder.WriteString(u.unqualifiedPasskeysTableName())
			builder.WriteString(" (")
			database.Columns{
				u.passkeyInstanceIDColumn(),
				u.passkeyUserIDColumn(),
				u.passkeyIDColumn(),
				u.passkeyKeyIDColumn(),
				u.passkeyCreatedAtColumn(),
				u.passkeyUpdatedAtColumn(),
				u.passkeyVerifiedAtColumn(),
				// TODO init_verification_id
				u.passkeyTypeColumn(),
				u.passkeyNameColumn(),
				u.passkeySignCountColumn(),
				u.passkeyChallengeColumn(),
				u.passkeyPublicKeyColumn(),
				u.passkeyAttestationTypeColumn(),
				u.passkeyAuthenticatorAttestationGUIDColumn(),
				u.passkeyRelyingPartyIDColumn(),
			}.WriteUnqualified(builder)
			builder.WriteString(") SELECT ")
			database.Columns{
				existingHumanUser.InstanceIDColumn(),
				existingHumanUser.idColumn(),
			}.WriteQualified(builder)
			builder.WriteString(", ")
			builder.WriteArgs(
				passkey.ID,
				passkey.KeyID,
				createdAt,
				updatedAt,
				verifiedAt,
				// TODO: verification_id
				passkey.Type,
				passkey.Name,
				passkey.SignCount,
				passkey.Challenge,
				passkey.PublicKey,
				passkey.AttestationType,
				passkey.AuthenticatorAttestationGUID,
				passkey.RelyingPartyID,
			)

		}, nil,
	)
}

// RemovePasskey implements [domain.HumanUserRepository.RemovePasskey].
func (u userHuman) RemovePasskey(condition database.Condition) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("DELETE FROM ")
			builder.WriteString(u.unqualifiedPasskeysTableName())
			builder.WriteString(" USING ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(existingHumanUser.InstanceIDColumn(), u.passkeyInstanceIDColumn()),
				database.NewColumnCondition(existingHumanUser.idColumn(), u.passkeyUserIDColumn()),
				condition,
			))
		},
		nil,
	)
}

// SetPasskeyAttestationType implements [domain.HumanUserRepository.SetPasskeyAttestationType].
func (u userHuman) SetPasskeyAttestationType(attestationType string) database.Change {
	return database.NewChange(u.passkeyAttestationTypeColumn(), attestationType)
}

// SetPasskeyAuthenticatorAttestationGUID implements [domain.HumanUserRepository.SetPasskeyAuthenticatorAttestationGUID].
func (u userHuman) SetPasskeyAuthenticatorAttestationGUID(aaguid []byte) database.Change {
	return database.NewChange(u.passkeyAuthenticatorAttestationGUIDColumn(), aaguid)
}

// SetPasskeyKeyID implements [domain.HumanUserRepository.SetPasskeyKeyID].
func (u userHuman) SetPasskeyKeyID(keyID []byte) database.Change {
	return database.NewChange(u.passkeyKeyIDColumn(), keyID)
}

// SetPasskeyName implements [domain.HumanUserRepository.SetPasskeyName].
func (u userHuman) SetPasskeyName(name string) database.Change {
	return database.NewChange(u.passkeyNameColumn(), name)
}

// SetPasskeyPublicKey implements [domain.HumanUserRepository.SetPasskeyPublicKey].
func (u userHuman) SetPasskeyPublicKey(publicKey []byte) database.Change {
	return database.NewChange(u.passkeyPublicKeyColumn(), publicKey)
}

// SetPasskeySignCount implements [domain.HumanUserRepository.SetPasskeySignCount].
func (u userHuman) SetPasskeySignCount(signCount uint32) database.Change {
	return database.NewChange(u.passkeySignCountColumn(), signCount)
}

// SetPasskeyUpdatedAt implements [domain.HumanUserRepository.SetPasskeyUpdatedAt].
func (u userHuman) SetPasskeyUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(u.passkeyUpdatedAtColumn(), updatedAt)
}

// SetPasskeyVerifiedAt implements [domain.HumanUserRepository.SetPasskeyVerifiedAt].
func (u userHuman) SetPasskeyVerifiedAt(verifiedAt time.Time) database.Change {
	verifiedAtChange := database.NewChange(u.passkeyVerifiedAtColumn(), database.NowInstruction)
	if !verifiedAt.IsZero() {
		verifiedAtChange = database.NewChange(u.passkeyVerifiedAtColumn(), verifiedAt)
	}
	return database.Changes{
		verifiedAtChange,
		database.NewChange(u.passkeyInitVerificationIDColumn(), database.NullInstruction),
	}
}

// UpdatePasskey implements [domain.HumanUserRepository.UpdatePasskey].
func (u userHuman) UpdatePasskey(condition database.Condition, changes ...database.Change) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("UPDATE ")
			builder.WriteString(u.unqualifiedPasskeysTableName())
			builder.WriteString(" SET ")
			database.Changes(changes).Write(builder)
			builder.WriteString(" FROM ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(existingHumanUser.InstanceIDColumn(), u.passkeyInstanceIDColumn()),
				database.NewColumnCondition(existingHumanUser.idColumn(), u.passkeyUserIDColumn()),
				condition,
			))
		},
		nil,
	)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (u userHuman) ExistsPasskey(condition database.Condition) database.Condition {
	panic("unimplemented")
}

func (u userHuman) PasskeyConditions() domain.HumanPasskeyConditions {
	panic("unimplemented")
}

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

func (u userHuman) passkeyRelyingPartyIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "relying_party_id")
}

func (u userHuman) passkeyAuthenticatorAttestationGUIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "authenticator_attestation_guid")
}

func (u userHuman) passkeyAttestationTypeColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "attestation_type")
}

func (u userHuman) passkeyPublicKeyColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "public_key")
}

func (u userHuman) passkeyChallengeColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "challenge")
}

func (u userHuman) passkeySignCountColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "sign_count")
}

func (u userHuman) passkeyNameColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "name")
}

func (u userHuman) passkeyTypeColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "type")
}

func (u userHuman) passkeyVerifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "verified_at")
}

func (u userHuman) passkeyUpdatedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "updated_at")
}

func (u userHuman) passkeyCreatedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "created_at")
}

func (u userHuman) passkeyInitVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedPasskeysTableName(), "init_verification_id")
}
