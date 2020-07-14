package repository

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/view/model"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type actualSequece struct {
	ActualSequence   uint64    `gorm:"column:current_sequence"`
	CurrentTimestamp time.Time `gorm:"column:timestamp"`
}

type CurrentSequence struct {
	ViewName         string    `gorm:"column:view_name;primary_key"`
	CurrentSequence  uint64    `gorm:"column:current_sequence"`
	CurrentTimestamp time.Time `gorm:"column:timestamp"`
}

type SequenceSearchKey int32

const (
	SequenceSearchKeyUndefined SequenceSearchKey = iota
	SequenceSearchKeyViewName
)

type sequenceSearchKey SequenceSearchKey

func (key sequenceSearchKey) ToColumnName() string {
	switch SequenceSearchKey(key) {
	case SequenceSearchKeyViewName:
		return "view_name"
	default:
		return ""
	}
}

func CurrentSequenceToModel(sequence *CurrentSequence) *model.View {
	dbView := strings.Split(sequence.ViewName, ".")
	return &model.View{
		Database:         dbView[0],
		ViewName:         dbView[1],
		CurrentSequence:  sequence.CurrentSequence,
		CurrentTimestamp: sequence.CurrentTimestamp,
	}
}

func SaveCurrentSequence(db *gorm.DB, table, viewName string, sequence uint64) error {
	save := PrepareSave(table)
	err := save(db, &CurrentSequence{viewName, sequence, time.Now()})

	if err != nil {
		return caos_errs.ThrowInternal(err, "VIEW-5kOhP", "unable to updated processed sequence")
	}
	return nil
}

func LatestSequence(db *gorm.DB, table, viewName string) (uint64, time.Time, error) {
	sequence := new(actualSequece)
	query := PrepareGetByKey(table, sequenceSearchKey(SequenceSearchKeyViewName), viewName)
	err := query(db, sequence)

	if err == nil {
		return sequence.ActualSequence, sequence.CurrentTimestamp, nil
	}

	if caos_errs.IsNotFound(err) {
		return 0, time.Time{}, nil
	}
	return 0, time.Time{}, caos_errs.ThrowInternalf(err, "VIEW-9LyCB", "unable to get latest sequence of %s", viewName)
}

func AllCurrentSequences(db *gorm.DB, table string) ([]*CurrentSequence, error) {
	sequences := make([]*CurrentSequence, 0)
	query := PrepareSearchQuery(table, GeneralSearchRequest{})
	_, err := query(db, &sequences)
	if err != nil {
		return nil, err
	}
	return sequences, nil
}

func ClearView(db *gorm.DB, truncateView, sequenceTable string) error {
	truncate := PrepareTruncate(truncateView)
	err := truncate(db)
	if err != nil {
		return err
	}
	return SaveCurrentSequence(db, sequenceTable, truncateView, 0)
}
