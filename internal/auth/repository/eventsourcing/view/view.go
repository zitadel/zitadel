package view

import (
	"database/sql"

	"github.com/jinzhu/gorm"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/id"
)

type View struct {
	Db              *gorm.DB
	keyAlgorithm    crypto.EncryptionAlgorithm
	idGenerator     id.Generator
	prefixAvatarURL string
}

func StartView(sqlClient *sql.DB, keyAlgorithm crypto.EncryptionAlgorithm, idGenerator id.Generator, prefixAvatarURL string) (*View, error) {
	gorm, err := gorm.Open("postgres", sqlClient)
	if err != nil {
		return nil, err
	}
	return &View{
		Db:              gorm,
		keyAlgorithm:    keyAlgorithm,
		idGenerator:     idGenerator,
		prefixAvatarURL: prefixAvatarURL,
	}, nil
}

func (v *View) Health() (err error) {
	return v.Db.DB().Ping()
}

func (v *View) PrefixAvatarURL() string {
	return v.prefixAvatarURL
}
