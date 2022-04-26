package repository

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/view/model"
)

type CurrentSequence struct {
	ViewName                 string    `gorm:"column:view_name;primary_key"`
	CurrentSequence          uint64    `gorm:"column:current_sequence"`
	EventTimestamp           time.Time `gorm:"column:event_timestamp"`
	LastSuccessfulSpoolerRun time.Time `gorm:"column:last_successful_spooler_run"`
	InstanceID               string    `gorm:"column:instance_id;primary_key"`
}

type currentSequenceViewWithSequence struct {
	ViewName                 string    `gorm:"column:view_name;primary_key"`
	CurrentSequence          uint64    `gorm:"column:current_sequence"`
	LastSuccessfulSpoolerRun time.Time `gorm:"column:last_successful_spooler_run"`
}

type currentSequenceView struct {
	ViewName                 string    `gorm:"column:view_name;primary_key"`
	LastSuccessfulSpoolerRun time.Time `gorm:"column:last_successful_spooler_run"`
}

type SequenceSearchKey int32

const (
	SequenceSearchKeyUndefined SequenceSearchKey = iota
	SequenceSearchKeyViewName
	SequenceSearchKeyAggregateType
	SequenceSearchKeyInstanceID
)

type sequenceSearchKey SequenceSearchKey

func (key sequenceSearchKey) ToColumnName() string {
	switch SequenceSearchKey(key) {
	case SequenceSearchKeyViewName:
		return "view_name"
	case SequenceSearchKeyAggregateType:
		return "aggregate_type"
	case SequenceSearchKeyInstanceID:
		return "instance_id"
	default:
		return ""
	}
}

type sequenceSearchQuery struct {
	key   sequenceSearchKey
	value string
}

func (q *sequenceSearchQuery) GetKey() ColumnKey {
	return q.key
}

func (q *sequenceSearchQuery) GetMethod() domain.SearchMethod {
	return domain.SearchMethodEquals
}

func (q *sequenceSearchQuery) GetValue() interface{} {
	return q.value
}

type sequenceSearchRequest struct {
	queries []sequenceSearchQuery
}

func (s *sequenceSearchRequest) GetLimit() uint64 {
	return 0
}

func (s *sequenceSearchRequest) GetOffset() uint64 {
	return 0
}

func (s *sequenceSearchRequest) GetSortingColumn() ColumnKey {
	return nil
}

func (s *sequenceSearchRequest) GetAsc() bool {
	return false
}

func (s *sequenceSearchRequest) GetQueries() []SearchQuery {
	result := make([]SearchQuery, len(s.queries))
	for i, q := range s.queries {
		result[i] = &sequenceSearchQuery{key: q.key, value: q.value}
	}
	return result
}

func CurrentSequenceToModel(sequence *CurrentSequence) *model.View {
	dbView := strings.Split(sequence.ViewName, ".")
	return &model.View{
		Database:                 dbView[0],
		ViewName:                 dbView[1],
		CurrentSequence:          sequence.CurrentSequence,
		EventTimestamp:           sequence.EventTimestamp,
		LastSuccessfulSpoolerRun: sequence.LastSuccessfulSpoolerRun,
	}
}

func SaveCurrentSequence(db *gorm.DB, table, viewName, instanceID string, sequence uint64, eventTimestamp time.Time) error {
	return UpdateCurrentSequence(db, table, &CurrentSequence{viewName, sequence, eventTimestamp, time.Now(), instanceID})
}

func SaveCurrentSequences(db *gorm.DB, table, viewName string, sequence uint64, eventTimestamp time.Time) error {
	err := db.Table(table).Where("view_name = ?", viewName).
		Updates(map[string]interface{}{"current_sequence": sequence, "event_timestamp": eventTimestamp, "last_successful_spooler_run": time.Now()}).Error
	if err != nil {
		return caos_errs.ThrowInternal(err, "VIEW-Sfdqs", "unable to updated processed sequence")
	}
	return nil
}

func UpdateCurrentSequence(db *gorm.DB, table string, currentSequence *CurrentSequence) (err error) {
	save := PrepareSave(table)
	err = save(db, currentSequence)
	if err != nil {
		return caos_errs.ThrowInternal(err, "VIEW-5kOhP", "unable to updated processed sequence")
	}
	return nil
}

func UpdateCurrentSequences(db *gorm.DB, table string, currentSequences []*CurrentSequence) (err error) {
	save := PrepareBulkSave(table)
	s := make([]interface{}, len(currentSequences))
	for i, currentSequence := range currentSequences {
		s[i] = currentSequence
	}
	err = save(db, s...)
	if err != nil {
		return caos_errs.ThrowInternal(err, "VIEW-5kOhP", "unable to updated processed sequence")
	}
	return nil
}

func LatestSequence(db *gorm.DB, table, viewName, instanceID string) (*CurrentSequence, error) {
	searchQueries := []SearchQuery{
		&sequenceSearchQuery{key: sequenceSearchKey(SequenceSearchKeyViewName), value: viewName},
		&sequenceSearchQuery{key: sequenceSearchKey(SequenceSearchKeyInstanceID), value: instanceID},
	}

	// ensure highest sequence of view
	db = db.Order("current_sequence DESC")

	query := PrepareGetByQuery(table, searchQueries...)
	sequence := new(CurrentSequence)
	err := query(db, sequence)

	if err == nil {
		return sequence, nil
	}

	if caos_errs.IsNotFound(err) {
		return sequence, nil
	}
	return nil, caos_errs.ThrowInternalf(err, "VIEW-9LyCB", "unable to get latest sequence of %s", viewName)
}

func LatestSequences(db *gorm.DB, table, viewName string) ([]*CurrentSequence, error) {
	searchQueries := make([]SearchQuery, 0, 2)
	searchQueries = append(searchQueries)
	searchRequest := &sequenceSearchRequest{
		queries: []sequenceSearchQuery{
			{key: sequenceSearchKey(SequenceSearchKeyViewName), value: viewName},
		},
	}

	// ensure highest sequence of view
	db = db.Order("current_sequence DESC")

	sequences := make([]*CurrentSequence, 0)
	query := PrepareSearchQuery(table, searchRequest)
	_, err := query(db, &sequences)
	if err != nil {
		return nil, err
	}
	return sequences, nil
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
	return SaveCurrentSequences(db, sequenceTable, truncateView, 0, time.Now())
}
