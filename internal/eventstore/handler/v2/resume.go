package handler

import "github.com/shopspring/decimal"

type resume struct {
	afterSortKey bool
	position     decimal.Decimal
	offset       uint32
}

func (h *Handler) resumeFrom(s *state, minPosition decimal.Decimal) resume {
	if minPosition.GreaterThan(decimal.NewFromInt(0)) {
		return resume{position: minPosition}
	}
	if s.cursor.IsZero() {
		return resume{}
	}
	if h.filterOffsetIsCursor || s.inTxOrderSet {
		return resume{afterSortKey: true}
	}
	r := resume{position: s.cursor.Position}
	if s.offset > 0 {
		r.offset = s.offset
	}
	return r
}
