package view

import (
	"errors"
	"fmt"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/jinzhu/gorm"
)

func PrepareGetByID(table string, key ColumnKey, id string) func(db *gorm.DB, res interface{}) error {
	return func(db *gorm.DB, res interface{}) error {
		err := db.Table(table).
			Where(fmt.Sprintf("%s = ?", key.ToColumnName()), id).
			Take(res).
			Error
		if err == nil {
			return nil
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return caos_errs.ThrowNotFound(err, "VIEW-XRI9c", "object not found")
		}
		logging.LogWithFields("VIEW-xVShS", "ID", id).WithError(err).Warn("get from view error")
		return caos_errs.ThrowInternal(err, "VIEW-J92Td", "view error")
	}
}

func PrepareGetByQuery(table string, queries ...SearchQuery) func(db *gorm.DB, res interface{}) error {
	return func(db *gorm.DB, res interface{}) error {
		query := db.Table(table)
		for _, q := range queries {
			var ok bool
			query, ok = SetQuery(query, q.GetKey(), q.GetValue(), q.GetMethod())
			if !ok {
				return caos_errs.ThrowInvalidArgument(nil, "VIEW-KaGue", "query is invalid")
			}
		}

		err := query.Take(res).Error
		if err == nil {
			return nil
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return caos_errs.ThrowNotFound(err, "VIEW-hodc6", "object not found")
		}
		logging.LogWithFields("VIEW-Mg6la", "table ", table).WithError(err).Warn("get from cache error")
		return caos_errs.ThrowInternal(err, "VIEW-qJBg9", "cache error")
	}
}

func PreparePut(table string) func(db *gorm.DB, res interface{}) error {
	return func(db *gorm.DB, object interface{}) error {
		err := db.Table(table).Save(object).Error
		if err != nil {
			return caos_errs.ThrowInternal(err, "VIEW-AfC7G", "unable to put object to view")
		}
		return nil
	}
}

func PrepareDelete(table string, key ColumnKey, id string) func(db *gorm.DB) error {
	return func(db *gorm.DB) error {
		err := db.Table(table).
			Where(fmt.Sprintf("%s = ?", key), id).
			Delete(nil).
			Error
		if err == nil {
			return caos_errs.ThrowInternal(err, "VIEW-die73", "could not delete object")
		}
		return nil
	}
}
