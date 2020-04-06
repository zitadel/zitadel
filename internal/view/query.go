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
			var ok bool
			query, ok = SetQuery(query, q.GetKey(), q.GetValue(), q.GetMethod())
			if !ok {
				return count, caos_errs.ThrowInvalidArgument(nil, "VIEW-KaGue", "query is invalid")
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

func SetQuery(query *gorm.DB, key ColumnKey, value interface{}, method model.SearchMethod) (*gorm.DB, bool) {
	column := key.ToColumnName()
	if column == "" {
		return nil, false
	}

	switch method {
	case model.Equals:
		query = query.Where("LOWER("+column+") = LOWER(?)", value)
	case model.EqualsCaseSensitive:
		query = query.Where(""+column+" = ?", value)
	case model.StartsWith:
		valueText, ok := value.(string)
		if !ok {
			return nil, false
		}
		query = query.Where("LOWER("+column+") LIKE LOWER(?)", valueText+"%")
	case model.StartsWithCaseSensitive:
		valueText, ok := value.(string)
		if !ok {
			return nil, false
		}
		query = query.Where(column+" LIKE ?", valueText+"%")
	case model.Contains:
		valueText, ok := value.(string)
		if !ok {
			return nil, false
		}
		query = query.Where("LOWER("+column+") LIKE LOWER(?)", "%"+valueText+"%")
	case model.ContainsCaseSensitive:
		valueText, ok := value.(string)
		if !ok {
			return nil, false
		}
		query = query.Where(column+" LIKE ?", "%"+valueText+"%")
	default:
		return nil, false
	}
	return query, true
}
