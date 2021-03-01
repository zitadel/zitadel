package repository

import (
	"github.com/caos/zitadel/internal/domain"
	"strings"

	"github.com/caos/zitadel/internal/errors"
	view_model "github.com/caos/zitadel/internal/view/model"
	"github.com/jinzhu/gorm"
)

const (
	errViewNameKey  = "view_name"
	errFailedSeqKey = "failed_sequence"
)

type FailedEvent struct {
	ViewName       string `gorm:"column:view_name;primary_key"`
	FailedSequence uint64 `gorm:"column:failed_sequence;primary_key"`
	FailureCount   uint64 `gorm:"column:failure_count"`
	ErrMsg         string `gorm:"column:err_msg"`
}

type FailedEventSearchQuery struct {
	Key    FailedEventSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

func (req FailedEventSearchQuery) GetKey() ColumnKey {
	return failedEventSearchKey(req.Key)
}

func (req FailedEventSearchQuery) GetMethod() domain.SearchMethod {
	return req.Method
}

func (req FailedEventSearchQuery) GetValue() interface{} {
	return req.Value
}

type FailedEventSearchKey int32

const (
	FailedEventKeyUndefined FailedEventSearchKey = iota
	FailedEventKeyViewName
	FailedEventKeyFailedSequence
)

type failedEventSearchKey FailedEventSearchKey

func (key failedEventSearchKey) ToColumnName() string {
	switch FailedEventSearchKey(key) {
	case FailedEventKeyViewName:
		return "view_name"
	case FailedEventKeyFailedSequence:
		return "failed_sequence"
	default:
		return ""
	}
}

func FailedEventFromModel(failedEvent *view_model.FailedEvent) *FailedEvent {
	return &FailedEvent{
		ViewName:       failedEvent.Database + "." + failedEvent.ViewName,
		FailureCount:   failedEvent.FailureCount,
		FailedSequence: failedEvent.FailedSequence,
		ErrMsg:         failedEvent.ErrMsg,
	}
}
func FailedEventToModel(failedEvent *FailedEvent) *view_model.FailedEvent {
	dbView := strings.Split(failedEvent.ViewName, ".")
	return &view_model.FailedEvent{
		Database:       dbView[0],
		ViewName:       dbView[1],
		FailureCount:   failedEvent.FailureCount,
		FailedSequence: failedEvent.FailedSequence,
		ErrMsg:         failedEvent.ErrMsg,
	}
}

func SaveFailedEvent(db *gorm.DB, table string, failedEvent *FailedEvent) error {
	save := PrepareSave(table)
	err := save(db, failedEvent)

	if err != nil {
		return errors.ThrowInternal(err, "VIEW-4F8us", "unable to updated failed events")
	}
	return nil
}

func RemoveFailedEvent(db *gorm.DB, table string, failedEvent *FailedEvent) error {
	delete := PrepareDeleteByKeys(table,
		Key{Key: failedEventSearchKey(FailedEventKeyViewName), Value: failedEvent.ViewName},
		Key{Key: failedEventSearchKey(FailedEventKeyFailedSequence), Value: failedEvent.FailedSequence},
	)
	return delete(db)
}

func LatestFailedEvent(db *gorm.DB, table, viewName string, sequence uint64) (*FailedEvent, error) {
	failedEvent := new(FailedEvent)
	queries := []SearchQuery{
		FailedEventSearchQuery{Key: FailedEventKeyViewName, Method: domain.SearchMethodEqualsIgnoreCase, Value: viewName},
		FailedEventSearchQuery{Key: FailedEventKeyFailedSequence, Method: domain.SearchMethodEquals, Value: sequence},
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

func AllFailedEvents(db *gorm.DB, table string) ([]*FailedEvent, error) {
	failedEvents := make([]*FailedEvent, 0)
	query := PrepareSearchQuery(table, GeneralSearchRequest{})
	_, err := query(db, &failedEvents)
	if err != nil {
		return nil, err
	}
	return failedEvents, nil
}
