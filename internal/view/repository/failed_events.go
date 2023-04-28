package repository

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	view_model "github.com/zitadel/zitadel/internal/view/model"
)

type FailedEvent struct {
	ViewName       string    `gorm:"column:view_name;primary_key"`
	FailedSequence uint64    `gorm:"column:failed_sequence;primary_key"`
	FailureCount   uint64    `gorm:"column:failure_count"`
	ErrMsg         string    `gorm:"column:err_msg"`
	InstanceID     string    `gorm:"column:instance_id"`
	LastFailed     time.Time `gorm:"column:last_failed"`
}

type failedEventSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn failedEventSearchKey
	Asc           bool
	Queries       []*FailedEventSearchQuery
}

func (f failedEventSearchRequest) GetLimit() uint64 {
	return f.Limit
}

func (f failedEventSearchRequest) GetOffset() uint64 {
	return f.Offset
}

func (f failedEventSearchRequest) GetSortingColumn() ColumnKey {
	if f.SortingColumn == failedEventSearchKey(FailedEventKeyUndefined) {
		return nil
	}
	return f.SortingColumn
}

func (f failedEventSearchRequest) GetAsc() bool {
	return f.Asc
}

func (f failedEventSearchRequest) GetQueries() []SearchQuery {
	result := make([]SearchQuery, len(f.Queries))
	for i, q := range f.Queries {
		result[i] = q
	}
	return result
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
	FailedEventKeyInstanceID
	FailedEventKeyLastFailed
)

type failedEventSearchKey FailedEventSearchKey

func (key failedEventSearchKey) ToColumnName() string {
	switch FailedEventSearchKey(key) {
	case FailedEventKeyViewName:
		return "view_name"
	case FailedEventKeyFailedSequence:
		return "failed_sequence"
	case FailedEventKeyInstanceID:
		return "instance_id"
	case FailedEventKeyLastFailed:
		return "last_failed"
	default:
		return ""
	}
}

func FailedEventFromModel(failedEvent *view_model.FailedEvent) *FailedEvent {
	return &FailedEvent{
		ViewName:       failedEvent.Database + "." + failedEvent.ViewName,
		FailureCount:   failedEvent.FailureCount,
		FailedSequence: failedEvent.FailedSequence,
		InstanceID:     failedEvent.InstanceID,
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
		LastFailed:     failedEvent.LastFailed,
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
		Key{Key: failedEventSearchKey(FailedEventKeyInstanceID), Value: failedEvent.InstanceID},
	)
	return delete(db)
}

func LatestFailedEvent(db *gorm.DB, table, viewName, instanceID string, sequence uint64) (*FailedEvent, error) {
	failedEvent := new(FailedEvent)
	queries := []SearchQuery{
		FailedEventSearchQuery{Key: FailedEventKeyViewName, Method: domain.SearchMethodEqualsIgnoreCase, Value: viewName},
		FailedEventSearchQuery{Key: FailedEventKeyFailedSequence, Method: domain.SearchMethodEquals, Value: sequence},
		FailedEventSearchQuery{Key: FailedEventKeyInstanceID, Method: domain.SearchMethodEquals, Value: instanceID},
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

func AllFailedEvents(db *gorm.DB, table, instanceID string) ([]*FailedEvent, error) {
	queries := make([]*FailedEventSearchQuery, 0, 1)
	if instanceID != "" {
		queries = append(queries, &FailedEventSearchQuery{Key: FailedEventKeyInstanceID, Method: domain.SearchMethodEquals, Value: instanceID})
	}
	failedEvents := make([]*FailedEvent, 0)
	query := PrepareSearchQuery(table, &failedEventSearchRequest{SortingColumn: failedEventSearchKey(FailedEventKeyLastFailed), Queries: queries})
	_, err := query(db, &failedEvents)
	if err != nil {
		return nil, err
	}
	return failedEvents, nil
}
