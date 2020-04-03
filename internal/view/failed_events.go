package view

import (
	"fmt"
	"github.com/caos/zitadel/internal/errors"
	"github.com/jinzhu/gorm"
)

const (
	errViewNameKey  = "view_name"
	errFailedSeqKey = "failed_sequence"
)

type FailedEvent struct {
	ViewName      string `gorm:"column:view_name;primary_key"`
	FailedSequnce uint64 `gorm:"column:failed_sequence;primary_key`
	FailureCount  uint64 `gorm:"column:failure_count`
}

func SaveFailedEvent(db *gorm.DB, table string, failedEvent *FailedEvent) error {
	err := db.Table(table).
		Save(failedEvent).
		Error

	if err != nil {
		return errors.ThrowInternal(err, "VIEW-5kOhP", "unable to updated failed events")
	}
	return nil
}

func LatestFailedEvent(db *gorm.DB, table, viewName string, sequence uint64) (*FailedEvent, error) {
	failedEvent := &FailedEvent{}
	err := db.Table(table).
		Where(fmt.Sprintf("%s = ?", errViewNameKey), viewName).
		Where(fmt.Sprintf("%s = ?", errFailedSeqKey), sequence).
		Scan(&sequence).
		Error

	if err == nil {
		return failedEvent, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		failedEvent.ViewName = viewName
		failedEvent.FailedSequnce = sequence
		failedEvent.FailureCount = 0
		return failedEvent, nil
	}
	return nil, errors.ThrowInternalf(err, "VIEW-9LyCB", "unable to get failed events of %s", viewName)

}
