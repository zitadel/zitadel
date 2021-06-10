package model

import (
	"time"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
)

type MessageTextsView struct {
	Texts   []*MessageTextView
	Default bool
}
type MessageTextView struct {
	AggregateID     string
	MessageTextType string
	Language        language.Tag
	Title           string
	PreHeader       string
	Subject         string
	Greeting        string
	Text            string
	ButtonText      string
	FooterText      string
	Default         bool

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type MessageTextSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn MessageTextSearchKey
	Asc           bool
	Queries       []*MessageTextSearchQuery
}

type MessageTextSearchKey int32

const (
	MessageTextSearchKeyUnspecified MessageTextSearchKey = iota
	MessageTextSearchKeyAggregateID
	MessageTextSearchKeyMessageTextType
	MessageTextSearchKeyLanguage
)

type MessageTextSearchQuery struct {
	Key    MessageTextSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type MessageTextSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*MessageTextView
	Sequence    uint64
	Timestamp   time.Time
}
