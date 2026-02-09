package repository

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userPasskeyRepo struct{}

func (u userPasskeyRepo) unqualifiedTableName() string {
	return "human_passkeys"
}

func (u userPasskeyRepo) qualifiedTableName() string {
	return "zitadel.human_passkeys"
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// AddPasskey implements [domain.HumanUserRepository.AddPasskey].
func (u userPasskeyRepo) AddPasskey(passkey *domain.Passkey) database.Change {
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
			builder.WriteString(u.qualifiedTableName())
			builder.WriteString(" (")
			database.Columns{
				u.instanceIDColumn(),
				u.userIDColumn(),
				u.tokenIDColumn(),
				u.keyIDColumn(),
				u.createdAtColumn(),
				u.updatedAtColumn(),
				u.verifiedAtColumn(),
				u.typeColumn(),
				u.nameColumn(),
				u.signCountColumn(),
				u.challengeColumn(),
				u.publicKeyColumn(),
				u.attestationTypeColumn(),
				u.authenticatorAttestationGUIDColumn(),
				u.relyingPartyIDColumn(),
			}.WriteUnqualified(builder)
			builder.WriteString(") SELECT ")
			database.Columns{
				existingHumanUser.InstanceIDColumn(),
				existingHumanUser.IDColumn(),
			}.WriteQualified(builder)
			builder.WriteString(", ")
			builder.WriteArgs(
				passkey.ID,
				passkey.KeyID,
				createdAt,
				updatedAt,
				verifiedAt,
				passkey.Type,
				passkey.Name,
				passkey.SignCount,
				passkey.Challenge,
				passkey.PublicKey,
				passkey.AttestationType,
				passkey.AuthenticatorAttestationGUID,
				passkey.RelyingPartyID,
			)
			builder.WriteString(" FROM ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
		}, nil,
	)
}

// RemovePasskey implements [domain.HumanUserRepository.RemovePasskey].
func (u userPasskeyRepo) RemovePasskey(condition database.Condition) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("DELETE FROM ")
			builder.WriteString(u.qualifiedTableName())
			builder.WriteString(" USING ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(existingHumanUser.InstanceIDColumn(), u.instanceIDColumn()),
				database.NewColumnCondition(existingHumanUser.IDColumn(), u.userIDColumn()),
				condition,
			))
		},
		nil,
	)
}

// UpdatePasskey implements [domain.HumanUserRepository.UpdatePasskey].
func (u userPasskeyRepo) UpdatePasskey(condition database.Condition, changes ...database.Change) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("UPDATE ")
			builder.WriteString(u.qualifiedTableName())
			builder.WriteString(" SET ")
			database.Changes(changes).Write(builder)
			builder.WriteString(" FROM ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(existingHumanUser.InstanceIDColumn(), u.instanceIDColumn()),
				database.NewColumnCondition(existingHumanUser.IDColumn(), u.userIDColumn()),
				condition,
			))
		},
		nil,
	)
}

// SetPasskeyAttestationType implements [domain.HumanUserRepository.SetPasskeyAttestationType].
func (u userPasskeyRepo) SetPasskeyAttestationType(attestationType string) database.Change {
	return database.NewChange(u.attestationTypeColumn(), attestationType)
}

// SetPasskeyAuthenticatorAttestationGUID implements [domain.HumanUserRepository.SetPasskeyAuthenticatorAttestationGUID].
func (u userPasskeyRepo) SetPasskeyAuthenticatorAttestationGUID(aaguid []byte) database.Change {
	return database.NewChange(u.authenticatorAttestationGUIDColumn(), aaguid)
}

// SetPasskeyKeyID implements [domain.HumanUserRepository.SetPasskeyKeyID].
func (u userPasskeyRepo) SetPasskeyKeyID(keyID []byte) database.Change {
	return database.NewChange(u.keyIDColumn(), keyID)
}

// SetPasskeyName implements [domain.HumanUserRepository.SetPasskeyName].
func (u userPasskeyRepo) SetPasskeyName(name string) database.Change {
	return database.NewChange(u.nameColumn(), name)
}

// SetPasskeyPublicKey implements [domain.HumanUserRepository.SetPasskeyPublicKey].
func (u userPasskeyRepo) SetPasskeyPublicKey(publicKey []byte) database.Change {
	return database.NewChange(u.publicKeyColumn(), publicKey)
}

// SetPasskeySignCount implements [domain.HumanUserRepository.SetPasskeySignCount].
func (u userPasskeyRepo) SetPasskeySignCount(signCount uint32) database.Change {
	return database.NewChange(u.signCountColumn(), signCount)
}

// SetPasskeyUpdatedAt implements [domain.HumanUserRepository.SetPasskeyUpdatedAt].
func (u userPasskeyRepo) SetPasskeyUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(u.updatedAtColumn(), updatedAt)
}

// SetPasskeyVerifiedAt implements [domain.HumanUserRepository.SetPasskeyVerifiedAt].
func (u userPasskeyRepo) SetPasskeyVerifiedAt(verifiedAt time.Time) database.Change {
	verifiedAtChange := database.NewChange(u.verifiedAtColumn(), database.NowInstruction)
	if !verifiedAt.IsZero() {
		verifiedAtChange = database.NewChange(u.verifiedAtColumn(), verifiedAt)
	}
	return verifiedAtChange
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// ExistsPasskey implements [domain.HumanUserRepository.ExistsPasskey].
func (u userPasskeyRepo) ExistsPasskey(condition database.Condition) database.Condition {
	panic("unimplemented")
}

// PasskeyConditions implements [domain.HumanUserRepository.PasskeyConditions].
func (u userPasskeyRepo) PasskeyConditions() domain.HumanPasskeyConditions {
	return u
}

// ChallengeCondition implements [domain.HumanPasskeyConditions.ChallengeCondition].
func (u userPasskeyRepo) ChallengeCondition(challenge string) database.Condition {
	// TODO: implement passkey challenge condition
	panic("unimplemented")
}

// IDCondition implements [domain.HumanPasskeyConditions.IDCondition].
func (u userPasskeyRepo) IDCondition(passkeyID string) database.Condition {
	return database.NewTextCondition(u.tokenIDColumn(), database.TextOperationEqual, passkeyID)
}

// KeyIDCondition implements [domain.HumanPasskeyConditions.KeyIDCondition].
func (u userPasskeyRepo) KeyIDCondition(keyID string) database.Condition {
	return database.NewTextCondition(u.keyIDColumn(), database.TextOperationEqual, keyID)
}

// TypeCondition implements [domain.HumanPasskeyConditions.TypeCondition].
func (u userPasskeyRepo) TypeCondition(passkeyType domain.PasskeyType) database.Condition {
	return database.NewTextCondition(u.typeColumn(), database.TextOperationEqual, string(passkeyType))
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userPasskeyRepo) tokenIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "token_id")
}

func (u userPasskeyRepo) keyIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "key_id")
}

func (u userPasskeyRepo) instanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "instance_id")
}

func (u userPasskeyRepo) userIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "user_id")
}

func (u userPasskeyRepo) relyingPartyIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "relying_party_id")
}

func (u userPasskeyRepo) authenticatorAttestationGUIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "authenticator_attestation_guid")
}

func (u userPasskeyRepo) attestationTypeColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "attestation_type")
}

func (u userPasskeyRepo) publicKeyColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "public_key")
}

func (u userPasskeyRepo) challengeColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "challenge")
}

func (u userPasskeyRepo) signCountColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "sign_count")
}

func (u userPasskeyRepo) nameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "name")
}

func (u userPasskeyRepo) typeColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "type")
}

func (u userPasskeyRepo) verifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "verified_at")
}

func (u userPasskeyRepo) updatedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "updated_at")
}

func (u userPasskeyRepo) createdAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "created_at")
}
