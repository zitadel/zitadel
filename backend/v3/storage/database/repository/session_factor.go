package repository

import (
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
		s.FactorTypeColumn(),
	}
}

func (s sessionFactor) InstanceIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "instance_id")
}

func (s sessionFactor) SessionIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "session_id")
}

func (s sessionFactor) FactorTypeColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "factor_type")
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

func (s sessionFactor) PayloadColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "payload")
}
