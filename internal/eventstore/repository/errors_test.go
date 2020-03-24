package repository

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type isExpectedError func(err error) bool

func isUnaddressable(err error) bool {
	return errors.Is(err, gorm.ErrUnaddressable)
}

func isRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
