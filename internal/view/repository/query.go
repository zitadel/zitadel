package repository

import (
	"fmt"
	"github.com/caos/zitadel/internal/domain"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
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
	GetMethod() domain.SearchMethod
	GetValue() interface{}
}

type ColumnKey interface {
	ToColumnName() string
}

func PrepareSearchQuery(table string, request SearchRequest) func(db *gorm.DB, res interface{}) (uint64, error) {
	return func(db *gorm.DB, res interface{}) (uint64, error) {
		var count uint64 = 0
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

func SetQuery(query *gorm.DB, key ColumnKey, value interface{}, method domain.SearchMethod) (*gorm.DB, error) {
	column := key.ToColumnName()
	if column == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "VIEW-7dz3w", "Column name missing")
	}

	switch method {
	case domain.SearchMethodEquals:
		query = query.Where(""+column+" = ?", value)
	case domain.SearchMethodEqualsIgnoreCase:
		valueText, ok := value.(string)
		if !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "VIEW-idu8e", "Equal ignore case only possible for strings")
		}
		query = query.Where("LOWER("+column+") = LOWER(?)", valueText)
	case domain.SearchMethodStartsWith:
		valueText, ok := value.(string)
		if !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "VIEW-SLj7s", "Starts with only possible for strings")
		}
		query = query.Where(column+" LIKE ?", valueText+"%")
	case domain.SearchMethodStartsWithIgnoreCase:
		valueText, ok := value.(string)
		if !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "VIEW-eidus", "Starts with ignore case only possible for strings")
		}
		query = query.Where("LOWER("+column+") LIKE LOWER(?)", valueText+"%")
	case domain.SearchMethodEndsWith:
		valueText, ok := value.(string)
		if !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "VIEW-Hswd3", "Ends with only possible for strings")
		}
		query = query.Where(column+" LIKE ?", "%"+valueText)
	case domain.SearchMethodEndsWithIgnoreCase:
		valueText, ok := value.(string)
		if !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "VIEW-dAG31", "Ends with ignore case only possible for strings")
		}
		query = query.Where("LOWER("+column+") LIKE LOWER(?)", "%"+valueText)
	case domain.SearchMethodContains:
		valueText, ok := value.(string)
		if !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "VIEW-3ids", "Contains with only possible for strings")
		}
		query = query.Where(column+" LIKE ?", "%"+valueText+"%")
	case domain.SearchMethodContainsIgnoreCase:
		valueText, ok := value.(string)
		if !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "VIEW-eid73", "Contains with ignore case only possible for strings")
		}
		query = query.Where("LOWER("+column+") LIKE LOWER(?)", "%"+valueText+"%")
	case domain.SearchMethodNotEquals:
		query = query.Where(""+column+" <> ?", value)
	case domain.SearchMethodGreaterThan:
		query = query.Where(column+" > ?", value)
	case domain.SearchMethodLessThan:
		query = query.Where(column+" < ?", value)
	case domain.SearchMethodIsOneOf:
		query = query.Where(column+" IN (?)", value)
	case domain.SearchMethodListContains:
		valueText, ok := value.(string)
		if !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "VIEW-Psois", "list contains only possible for strings")
		}
		query = query.Where("? <@ "+column, pq.Array([]string{valueText}))
	default:
		return nil, nil
	}
	return query, nil
}
