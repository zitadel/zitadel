package model

import (
	"time"

	"github.com/caos/zitadel/internal/domain"
)

type MailTextsView struct {
	Texts   []*MailTextView
	Default bool
}
type MailTextView struct {
	AggregateID  string
	MailTextType string
	Language     string
	Title        string
	PreHeader    string
	Subject      string
	Greeting     string
	Text         string
	ButtonText   string
	FooterText   string
	Default      bool

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type MailTextSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn MailTextSearchKey
	Asc           bool
	Queries       []*MailTextSearchQuery
}

type MailTextSearchKey int32

const (
	MailTextSearchKeyUnspecified MailTextSearchKey = iota
	MailTextSearchKeyAggregateID
	MailTextSearchKeyMailTextType
	MailTextSearchKeyLanguage
)

type MailTextSearchQuery struct {
	Key    MailTextSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type MailTextSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*MailTextView
	Sequence    uint64
	Timestamp   time.Time
}
