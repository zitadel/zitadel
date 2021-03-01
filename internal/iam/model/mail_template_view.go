package model

import (
	"github.com/caos/zitadel/internal/domain"
	"time"
)

type MailTemplateView struct {
	AggregateID string
	Template    []byte
	Default     bool

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type MailTemplateSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn MailTemplateSearchKey
	Asc           bool
	Queries       []*MailTemplateSearchQuery
}

type MailTemplateSearchKey int32

const (
	MailTemplateSearchKeyUnspecified MailTemplateSearchKey = iota
	MailTemplateSearchKeyAggregateID
)

type MailTemplateSearchQuery struct {
	Key    MailTemplateSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type MailTemplateSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*MailTemplateView
	Sequence    uint64
	Timestamp   time.Time
}
