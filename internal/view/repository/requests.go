package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func PrepareGetByQuery(table string, queries ...SearchQuery) func(db *gorm.DB, res interface{}) error {
	return func(db *gorm.DB, res interface{}) error {
		query := db.Table(table)
		for _, q := range queries {
			var err error
			query, err = SetQuery(query, q.GetKey(), q.GetValue(), q.GetMethod())
			if err != nil {
				return zerrors.ThrowInvalidArgument(err, "VIEW-KaGue", "query is invalid")
			}
		}

		err := query.Take(res).Error
		if err == nil {
			return nil
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return zerrors.ThrowNotFound(err, "VIEW-hodc6", "object not found")
		}
		logging.LogWithFields("VIEW-Mg6la", "table ", table).WithError(err).Warn("get from cache error")
		return zerrors.ThrowInternal(err, "VIEW-qJBg9", "cache error")
	}
}

func PrepareBulkSave(table string) func(db *gorm.DB, objects ...interface{}) error {
	return func(db *gorm.DB, objects ...interface{}) error {
		db = db.Table(table)
		db = db.Begin()
		defer db.RollbackUnlessCommitted()
		if err := db.Error; err != nil {
			return zerrors.ThrowInternal(err, "REPOS-Fl0Is", "unable to begin")
		}
		for _, object := range objects {
			err := db.Save(object).Error
			if err != nil {
				return zerrors.ThrowInternal(err, "VIEW-oJJSm", "unable to put object to view")
			}
		}
		if err := db.Commit().Error; err != nil {
			return zerrors.ThrowInternal(err, "REPOS-IfhUE", "unable to commit")
		}
		return nil
	}
}

func PrepareSave(table string) func(db *gorm.DB, object interface{}) error {
	return func(db *gorm.DB, object interface{}) error {
		err := db.Table(table).Save(object).Error
		if err != nil {
			return zerrors.ThrowInternal(err, "VIEW-2m9fs", "unable to put object to view")
		}
		return nil
	}
}

func PrepareSaveOnConflict(table string, conflictColumns, updateColumns []string) func(db *gorm.DB, object interface{}) error {
	updates := make([]string, len(updateColumns))
	for i, column := range updateColumns {
		updates[i] = column + "=excluded." + column
	}
	onConflict := fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s", strings.Join(conflictColumns, ","), strings.Join(updates, ","))
	return func(db *gorm.DB, object interface{}) error {
		err := db.Table(table).Set("gorm:insert_option", onConflict).Save(object).Error
		if err != nil {
			return zerrors.ThrowInternal(err, "VIEW-AfC7G", "unable to put object to view")
		}
		return nil
	}
}

func PrepareDeleteByKey(table string, key ColumnKey, id interface{}) func(db *gorm.DB) error {
	return func(db *gorm.DB) error {
		err := db.Table(table).
			Where(fmt.Sprintf("%s = ?", key.ToColumnName()), id).
			Delete(nil).
			Error
		if err != nil {
			return zerrors.ThrowInternal(err, "VIEW-die73", "could not delete object")
		}
		return nil
	}
}

func PrepareUpdateByKeys(table string, column ColumnKey, value interface{}, keys ...Key) func(db *gorm.DB) error {
	return func(db *gorm.DB) error {
		for _, key := range keys {
			db = db.Table(table).
				Where(fmt.Sprintf("%s = ?", key.Key.ToColumnName()), key.Value)
		}
		err := db.
			Update(column.ToColumnName(), value).
			Error
		if err != nil {
			return zerrors.ThrowInternal(err, "VIEW-ps099xj", "could not update object")
		}
		return nil
	}
}

type Key struct {
	Key   ColumnKey
	Value interface{}
}

func PrepareDeleteByKeys(table string, keys ...Key) func(db *gorm.DB) error {
	return func(db *gorm.DB) error {
		for _, key := range keys {
			db = db.Table(table).
				Where(fmt.Sprintf("%s = ?", key.Key.ToColumnName()), key.Value)
		}
		err := db.
			Delete(nil).
			Error
		if err != nil {
			return zerrors.ThrowInternal(err, "VIEW-die73", "could not delete object")
		}
		return nil
	}
}
