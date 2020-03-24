package repository

import "github.com/jinzhu/gorm"

type aggregateRoot interface {
	ID() string
	Type() string
}

func whereAggregateRoot(db *gorm.DB, root aggregateRoot) *gorm.DB {
	return db.Where("aggregate_id = ? and aggregate_type = ?", root.ID(), root.Type())
}
