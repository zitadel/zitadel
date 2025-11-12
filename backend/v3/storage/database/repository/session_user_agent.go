package repository

import (
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type sessionUserAgent struct {
}

func (s sessionUserAgent) qualifiedTableName() string {
	return "zitadel.session_user_agent"
}

func (s sessionUserAgent) unqualifiedTableName() string {
	return "session_user_agent"
}

func (s sessionUserAgent) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		s.InstanceIDColumn(),
		s.FingerprintIDColumn(),
	}
}

func (s sessionUserAgent) InstanceIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "instance_id")
}

func (s sessionUserAgent) FingerprintIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "fingerprint_id")
}

func (s sessionUserAgent) DescriptionColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "description")
}

func (s sessionUserAgent) IPColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "ip")
}

func (s sessionUserAgent) HeadersColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "headers")
}
