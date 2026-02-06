package repository

import (
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type sessionUserAgent struct {
}

func (s sessionUserAgent) qualifiedTableName() string {
	return "zitadel.session_user_agents"
}

func (s sessionUserAgent) unqualifiedTableName() string {
	return "session_user_agents"
}

func (s sessionUserAgent) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		s.instanceIDColumn(),
		s.fingerprintIDColumn(),
	}
}

func (s sessionUserAgent) instanceIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "instance_id")
}

func (s sessionUserAgent) fingerprintIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "fingerprint_id")
}

func (s sessionUserAgent) descriptionColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "description")
}

func (s sessionUserAgent) ipColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "ip")
}

func (s sessionUserAgent) headersColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "headers")
}
