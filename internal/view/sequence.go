package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/jinzhu/gorm"
)

type actualSequece struct {
	ActualSequence uint64 `gorm:"column:current_sequence"`
}

type currentSequence struct {
	ViewName        string `gorm:"column:view_name;primary_key"`
	CurrentSequence uint64 `gorm:"column:current_sequence`
}

type SequenceSearchKey int32

const (
	SEQUENCESEARCHKEY_UNDEFINED SequenceSearchKey = iota
	SEQUENCESEARCHKEY_VIEW_NAME
)

type sequenceSearchKey SequenceSearchKey

func (key sequenceSearchKey) ToColumnName() string {
	switch SequenceSearchKey(key) {
	case SEQUENCESEARCHKEY_VIEW_NAME:
		return "view_name"
	default:
		return ""
	}
}

func SaveCurrentSequence(db *gorm.DB, table, viewName string, sequence uint64) error {
	save := PrepareSave(table)
	err := save(db, &currentSequence{viewName, sequence})

	if err != nil {
		return caos_errs.ThrowInternal(err, "VIEW-5kOhP", "unable to updated processed sequence")
	}
	return nil
}

func LatestSequence(db *gorm.DB, table, viewName string) (uint64, error) {
	sequence := actualSequece{}
	query := PrepareGetByKey(table, sequenceSearchKey(SEQUENCESEARCHKEY_VIEW_NAME), viewName)
	err := query(db, sequence)

	if err == nil {
		return sequence.ActualSequence, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		return 0, nil
	}
	return 0, caos_errs.ThrowInternalf(err, "VIEW-9LyCB", "unable to get latest sequence of %s", viewName)
}
