package repository

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userPasskey struct{}

func (u userPasskey) unqualifiedTableName() string {
	return "user_passkeys"
}

func (u userPasskey) qualifiedTableName() string {
	return "zitadel.user_passkeys"
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// AddPasskey implements [domain.HumanUserRepository].
func (u userPasskey) AddPasskey(passkey *domain.Passkey) database.Change {
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

// RemovePasskey implements [domain.HumanUserRepository].
func (u userPasskey) RemovePasskey(condition database.Condition) database.Change {
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

// UpdatePasskey implements [domain.HumanUserRepository].
func (u userPasskey) UpdatePasskey(condition database.Condition, changes ...database.Change) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("UPDATE ")
			builder.WriteString(u.qualifiedTableName())
			builder.WriteString(" SET ")
			err := database.Changes(changes).Write(builder)
			logging.New(logging.StreamRuntime).Debug("write changes in cte failed", "error", err)
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

// SetPasskeyAttestationType implements [domain.HumanUserRepository].
func (u userPasskey) SetPasskeyAttestationType(attestationType string) database.Change {
	return database.NewChange(u.attestationTypeColumn(), attestationType)
}

// SetPasskeyAuthenticatorAttestationGUID implements [domain.HumanUserRepository].
func (u userPasskey) SetPasskeyAuthenticatorAttestationGUID(aaguid []byte) database.Change {
	return database.NewChange(u.authenticatorAttestationGUIDColumn(), aaguid)
}

// SetPasskeyKeyID implements [domain.HumanUserRepository].
func (u userPasskey) SetPasskeyKeyID(keyID []byte) database.Change {
	return database.NewChange(u.keyIDColumn(), keyID)
}

// SetPasskeyName implements [domain.HumanUserRepository].
func (u userPasskey) SetPasskeyName(name string) database.Change {
	return database.NewChange(u.nameColumn(), name)
}

// SetPasskeyPublicKey implements [domain.HumanUserRepository].
func (u userPasskey) SetPasskeyPublicKey(publicKey []byte) database.Change {
	return database.NewChange(u.publicKeyColumn(), publicKey)
}

// SetPasskeySignCount implements [domain.HumanUserRepository].
func (u userPasskey) SetPasskeySignCount(signCount uint32) database.Change {
	return database.NewChange(u.signCountColumn(), signCount)
}

// SetPasskeyUpdatedAt implements [domain.HumanUserRepository].
func (u userPasskey) SetPasskeyUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(u.updatedAtColumn(), updatedAt)
}

// SetPasskeyVerifiedAt implements [domain.HumanUserRepository].
func (u userPasskey) SetPasskeyVerifiedAt(verifiedAt time.Time) database.Change {
	verifiedAtChange := database.NewChange(u.verifiedAtColumn(), database.NowInstruction)
	if !verifiedAt.IsZero() {
		verifiedAtChange = database.NewChange(u.verifiedAtColumn(), verifiedAt)
	}
	return verifiedAtChange
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// PasskeyConditions implements [domain.HumanUserRepository].
func (u userPasskey) PasskeyConditions() domain.HumanPasskeyConditions {
	return u
}

// ChallengeCondition implements [domain.HumanPasskeyConditions].
func (u userPasskey) ChallengeCondition(challenge []byte) database.Condition {
	return database.NewBytesCondition[[]byte](database.SHA256Column(u.challengeColumn()), database.BytesOperationEqual, database.SHA256Value(challenge))
}

// IDCondition implements [domain.HumanPasskeyConditions].
func (u userPasskey) IDCondition(passkeyID string) database.Condition {
	return database.NewTextCondition(u.tokenIDColumn(), database.TextOperationEqual, passkeyID)
}

// KeyIDCondition implements [domain.HumanPasskeyConditions].
func (u userPasskey) KeyIDCondition(keyID string) database.Condition {
	return database.NewTextCondition(u.keyIDColumn(), database.TextOperationEqual, keyID)
}

// TypeCondition implements [domain.HumanPasskeyConditions].
func (u userPasskey) TypeCondition(passkeyType domain.PasskeyType) database.Condition {
	return database.NewTextCondition(u.typeColumn(), database.TextOperationEqual, passkeyType.String())
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userPasskey) tokenIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "token_id")
}

func (u userPasskey) keyIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "key_id")
}

func (u userPasskey) instanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "instance_id")
}

func (u userPasskey) userIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "user_id")
}

func (u userPasskey) relyingPartyIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "relying_party_id")
}

func (u userPasskey) authenticatorAttestationGUIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "authenticator_attestation_guid")
}

func (u userPasskey) attestationTypeColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "attestation_type")
}

func (u userPasskey) publicKeyColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "public_key")
}

func (u userPasskey) challengeColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "challenge")
}

func (u userPasskey) signCountColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "sign_count")
}

func (u userPasskey) nameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "name")
}

func (u userPasskey) typeColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "type")
}

func (u userPasskey) verifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "verified_at")
}

func (u userPasskey) updatedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "updated_at")
}

func (u userPasskey) createdAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "created_at")
}
