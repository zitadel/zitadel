package model

import (
	"time"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
)

type CustomTextView struct {
	AggregateID string
	Template    string
	Language    language.Tag
	Key         string
	Text        string

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type CustomTextSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn CustomTextSearchKey
	Asc           bool
	Queries       []*CustomTextSearchQuery
}

type CustomTextSearchKey int32

const (
	CustomTextSearchKeyUnspecified CustomTextSearchKey = iota
	CustomTextSearchKeyAggregateID
	CustomTextSearchKeyTemplate
	CustomTextSearchKeyLanguage
	CustomTextSearchKeyKey
)

type CustomTextSearchQuery struct {
	Key    CustomTextSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type CustomTextSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*CustomTextView
	Sequence    uint64
	Timestamp   time.Time
}
