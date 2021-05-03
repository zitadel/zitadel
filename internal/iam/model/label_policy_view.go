package model

import (
	"time"

	"github.com/caos/zitadel/internal/domain"
)

type LabelPolicyView struct {
	AggregateID    string
	PrimaryColor   string
	SecondaryColor string
	WarnColor      string
	LogoURL        string
	IconURL        string

	PrimaryColorDark   string
	SecondaryColorDark string
	WarnColorDark      string
	LogoURLDark        string
	IconURLDark        string

	HideLoginNameSuffix bool
	ErrorMsgPopup       bool
	DisableWatermark    bool

	Default bool

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type LabelPolicySearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn LabelPolicySearchKey
	Asc           bool
	Queries       []*LabelPolicySearchQuery
}

type LabelPolicySearchKey int32

const (
	LabelPolicySearchKeyUnspecified LabelPolicySearchKey = iota
	LabelPolicySearchKeyAggregateID
)

type LabelPolicySearchQuery struct {
	Key    LabelPolicySearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type LabelPolicySearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*LabelPolicyView
	Sequence    uint64
	Timestamp   time.Time
}
