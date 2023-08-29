package view

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	eventstore "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/query"
)

type View struct {
	Db           *gorm.DB
	keyAlgorithm crypto.EncryptionAlgorithm
	idGenerator  id.Generator
	query        *query.Queries
	es           eventstore.Eventstore
	client       *database.DB
}

func StartView(sqlClient *database.DB, keyAlgorithm crypto.EncryptionAlgorithm, queries *query.Queries, idGenerator id.Generator, es eventstore.Eventstore) (*View, error) {
	gorm, err := gorm.Open("postgres", sqlClient.DB)
	if err != nil {
		return nil, err
	}
	return &View{
		Db:           gorm,
		keyAlgorithm: keyAlgorithm,
		idGenerator:  idGenerator,
		query:        queries,
		es:           es,
		client:       sqlClient,
	}, nil
}

func (v *View) Health() (err error) {
	return v.Db.DB().Ping()
}

func (v *View) TimeTravel(ctx context.Context, tableName string) string {
	return tableName + v.client.Timetravel(call.Took(ctx))
}
