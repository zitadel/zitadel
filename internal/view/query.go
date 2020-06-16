package view

import (
	"fmt"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/model"
	"github.com/jinzhu/gorm"
)

type SearchRequest interface {
	GetLimit() uint64
	GetOffset() uint64
	GetSortingColumn() ColumnKey
	GetAsc() bool
	GetQueries() []SearchQuery
}

type SearchQuery interface {
	GetKey() ColumnKey
	GetMethod() model.SearchMethod
	GetValue() interface{}
}

type ColumnKey interface {
	ToColumnName() string
}

func PrepareSearchQuery(table string, request SearchRequest) func(db *gorm.DB, res interface{}) (int, error) {
	return func(db *gorm.DB, res interface{}) (int, error) {
		count := 0
		query := db.Table(table)
		if column := request.GetSortingColumn(); column != nil {
			order := "DESC"
			if request.GetAsc() {
				order = "ASC"
			}
			query = query.Order(fmt.Sprintf("%s %s", column.ToColumnName(), order))
		}
		for _, q := range request.GetQueries() {
			var err error
			query, err = SetQuery(query, q.GetKey(), q.GetValue(), q.GetMethod())
			if err != nil {
				return count, caos_errs.ThrowInvalidArgument(err, "VIEW-KaGue", "query is invalid")
			}
		}

		query = query.Count(&count)
		if request.GetLimit() != 0 {
			query = query.Limit(request.GetLimit())
		}
		query = query.Offset(request.GetOffset())
		err := query.Find(res).Error
		if err != nil {
			return count, caos_errs.ThrowInternal(err, "VIEW-muSDK", "unable to find result")
		}
		return count, nil
	}
}

func SetQuery(query *gorm.DB, key ColumnKey, value interface{}, method model.SearchMethod) (*gorm.DB, error) {
	column := key.ToColumnName()
	if column == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "VIEW-7dz3w", "Column name missing")
	}

	switch method {
	case model.SEARCHMETHOD_EQUALS:
		query = query.Where(""+column+" = ?", value)
	case model.SEARCHMETHOD_EQUALS_IGNORE_CASE:
		valueText, ok := value.(string)
		if !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "VIEW-idu8e", "Equal ignore case only possible for strings")
		}
		query = query.Where("LOWER("+column+") = LOWER(?)", valueText)
	case model.SEARCHMETHOD_STARTS_WITH:
		valueText, ok := value.(string)
		if !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "VIEW-idu8e", "Starts with only possible for strings")
		}
		query = query.Where(column+" LIKE ?", valueText+"%")
	case model.SEARCHMETHOD_STARTS_WITH_IGNORE_CASE:
		valueText, ok := value.(string)
		if !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "VIEW-eidus", "Starts with ignore case only possible for strings")
		}
		query = query.Where("LOWER("+column+") LIKE LOWER(?)", valueText+"%")
	case model.SEARCHMETHOD_CONTAINS:
		valueText, ok := value.(string)
		if !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "VIEW-3ids", "Contains with only possible for strings")
		}
		query = query.Where(column+" LIKE ?", "%"+valueText+"%")
	case model.SEARCHMETHOD_CONTAINS_IGNORE_CASE:
		valueText, ok := value.(string)
		if !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "VIEW-eid73", "Contains with ignore case only possible for strings")
		}
		query = query.Where("LOWER("+column+") LIKE LOWER(?)", "%"+valueText+"%")
	case model.SEARCHMETHOD_NOT_EQUALS:
		query = query.Where(""+column+" <> ?", value)
	case model.SEARCHMETHOD_GREATER_THAN:
		query = query.Where(column+" > ?", value)
	case model.SEARCHMETHOD_LESS_THAN:
		query = query.Where(column+" < ?", value)
	case model.SEARCHMETHOD_IN:
		query = query.Where(column+" IN (?)", value)
	case model.SEARCHMETHOD_EQUALS_IN_ARRAY:
		query = query.Where("? <@ "+column, value)
	default:
		return nil, nil
	}
	return query, nil
}
