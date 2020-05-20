package view

import (
	"database/sql"

	"github.com/jinzhu/gorm"

	"github.com/caos/zitadel/internal/crypto"
)

type View struct {
	Db           *gorm.DB
	keyAlgorithm crypto.EncryptionAlgorithm
}

func StartView(sqlClient *sql.DB, keyAlgorithm crypto.EncryptionAlgorithm) (*View, error) {
	gorm, err := gorm.Open("postgres", sqlClient)
	if err != nil {
		return nil, err
	}
	return &View{
		Db:           gorm,
		keyAlgorithm: keyAlgorithm,
	}, nil
}

func (v *View) Health() (err error) {
	return v.Db.DB().Ping()
}
