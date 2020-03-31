package view

import (
	"fmt"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/jinzhu/gorm"
)

const (
	viewNameKey = "view_name"
)

type actualSequece struct {
	ActualSequence uint64 `gorm:"column:current_sequence"`
}

type currentSequence struct {
	ViewName        string `gorm:"column:view_name;primary_key"`
	CurrentSequence uint64 `gorm:"column:current_sequence`
}

func SaveCurrentSequence(db *gorm.DB, table, viewName string, sequence uint64) error {
	err := db.Table(table).
		Save(&currentSequence{viewName, sequence}).
		Error

	if err != nil {
		return caos_errs.ThrowInternal(err, "VIEW-5kOhP", "unable to updated processed sequence")
	}
	return nil
}

func LatestSequence(db *gorm.DB, table, viewName string) (uint64, error) {
	sequence := actualSequece{}
	err := db.Table(table).
		Where(fmt.Sprintf("%s = ?", viewNameKey), viewName).
		Scan(&sequence).
		Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return 0, nil
		}
		return 0, caos_errs.ThrowInternalf(err, "VIEW-9LyCB", "unable to get latest sequence of %s", viewName)
	}

	return sequence.ActualSequence, nil
}
