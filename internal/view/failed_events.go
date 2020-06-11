package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/model"
	"github.com/jinzhu/gorm"
)

const (
	errViewNameKey  = "view_name"
	errFailedSeqKey = "failed_sequence"
)

type FailedEvent struct {
	ViewName       string `gorm:"column:view_name;primary_key"`
	FailedSequence uint64 `gorm:"column:failed_sequence;primary_key`
	FailureCount   uint64 `gorm:"column:failure_count`
	ErrMsg         string `gorm:"column:err_msg`
}

type FailedEventSearchQuery struct {
	Key    FailedEventSearchKey
	Method model.SearchMethod
	Value  interface{}
}

func (req FailedEventSearchQuery) GetKey() ColumnKey {
	return failedEventSearchKey(req.Key)
}

func (req FailedEventSearchQuery) GetMethod() model.SearchMethod {
	return req.Method
}

func (req FailedEventSearchQuery) GetValue() interface{} {
	return req.Value
}

type FailedEventSearchKey int32

const (
	FAILEDEVENTKEY_UNDEFINED FailedEventSearchKey = iota
	FAILEDEVENTKEY_VIEW_NAME
	FAILEDEVENTKEY_FAILED_SEQUENCE
)

type failedEventSearchKey FailedEventSearchKey

func (key failedEventSearchKey) ToColumnName() string {
	switch FailedEventSearchKey(key) {
	case FAILEDEVENTKEY_VIEW_NAME:
		return "view_name"
	case FAILEDEVENTKEY_FAILED_SEQUENCE:
		return "failed_sequence"
	default:
		return ""
	}
}

func SaveFailedEvent(db *gorm.DB, table string, failedEvent *FailedEvent) error {
	save := PrepareSave(table)
	err := save(db, failedEvent)

	if err != nil {
		return errors.ThrowInternal(err, "VIEW-5kOhP", "unable to updated failed events")
	}
	return nil
}

func LatestFailedEvent(db *gorm.DB, table, viewName string, sequence uint64) (*FailedEvent, error) {
	failedEvent := new(FailedEvent)
	queries := []SearchQuery{
		FailedEventSearchQuery{Key: FAILEDEVENTKEY_VIEW_NAME, Method: model.SEARCHMETHOD_EQUALS_IGNORE_CASE, Value: viewName},
		FailedEventSearchQuery{Key: FAILEDEVENTKEY_FAILED_SEQUENCE, Method: model.SEARCHMETHOD_EQUALS, Value: sequence},
	}
	query := PrepareGetByQuery(table, queries...)
	err := query(db, failedEvent)

	if err == nil && failedEvent.ViewName != "" {
		return failedEvent, nil
	}

	if errors.IsNotFound(err) {
		return &FailedEvent{
			ViewName:       viewName,
			FailedSequence: sequence,
			FailureCount:   0,
		}, nil
	}
	return nil, errors.ThrowInternalf(err, "VIEW-9LyCB", "unable to get failed events of %s", viewName)

}
