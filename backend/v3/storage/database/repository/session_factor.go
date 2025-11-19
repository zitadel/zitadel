package repository

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type sessionFactor struct {
}

func (s sessionFactor) qualifiedTableName() string {
	return "zitadel.session_factors"
}

func (s sessionFactor) unqualifiedTableName() string {
	return "session_factors"
}

func (s sessionFactor) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		s.InstanceIDColumn(),
		s.SessionIDColumn(),
		s.TypeColumn(),
	}
}

func (s sessionFactor) InstanceIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "instance_id")
}

func (s sessionFactor) SessionIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "session_id")
}

func (s sessionFactor) TypeColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "type")
}

func (s sessionFactor) LastChallengedAtColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "last_challenged_at")
}

func (s sessionFactor) LastFailedAtColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "last_failed_at")
}

func (s sessionFactor) LastVerifiedAtColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "last_verified_at")
}

func (s sessionFactor) ChallengedPayloadColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "challenged_payload")
}

func (s sessionFactor) VerifiedPayloadColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "verified_payload")
}

//TODO: below conditions need to be fixed

func (s sessionFactor) FactorTypeCondition(factorType domain.SessionFactorType) database.Condition {
	return database.NewTextCondition(s.TypeColumn(), database.TextOperationEqual, "factorType")
}

func (s sessionFactor) LastVerifiedBeforeCondition(lastVerifiedAt time.Time) database.Condition {
	return database.NewTextCondition(s.LastVerifiedAtColumn(), database.TextOperationContains, "")
}
