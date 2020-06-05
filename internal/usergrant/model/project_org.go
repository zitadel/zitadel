package model

type ProjectOrgSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*Org
}

type Org struct {
	OrgID   string
	OrgName string
}
