package view

import (
	"database/sql"

	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/query"
)

type View struct {
	Db           *gorm.DB
	keyAlgorithm crypto.EncryptionAlgorithm
	idGenerator  id.Generator
	query        *query.Queries
}

func StartView(sqlClient *sql.DB, keyAlgorithm crypto.EncryptionAlgorithm, queries *query.Queries, idGenerator id.Generator) (*View, error) {
	gorm, err := gorm.Open("postgres", sqlClient)
	if err != nil {
		return nil, err
	}
	return &View{
		Db:           gorm,
		keyAlgorithm: keyAlgorithm,
		idGenerator:  idGenerator,
		query:        queries,
	}, nil
}

func (v *View) Health() (err error) {
	return v.Db.DB().Ping()
}
