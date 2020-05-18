package model

import (
	"database/sql"
	"time"

	"github.com/caos/zitadel/internal/crypto"
)

const (
	KeyKeyID     = "id"
	KeyPrivate   = "private"
	KeyUsage     = "usage"
	KeyAlgorithm = "algorithm"
	KeyExpiry    = "expiry"
)

type KeyView struct {
	ID              string              `json:"-" gorm:"column:key_id;primary_key"`
	Private         sql.NullBool        `json:"-" gorm:"column:private;primary_key"`
	Expiry          time.Time           `json:"-" gorm:"column:expiry"`
	Algorithm       string              `json:"-" gorm:"column:algorithm"`
	Usage           string              `json:"-" gorm:"column:usage"`
	Key             *crypto.CryptoValue `json:"key" gorm:"column:key"`
	CurrentSequence uint64              `gorm:"column:current_sequence"`
}
