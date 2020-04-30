package model

import "time"

type GrantedProject struct {
	ProjectID     string
	Name          string
	CreationDate  time.Time
	ChangeDate    time.Time
	State         ProjectState
	Type          ProjectType
	ResourceOwner string
	OrgID         string
	OrgName       string
	OrgDomain     string
	Sequence      uint64
	GrantID       string
}

type ProjectType int32

const (
	PROJECTTYPE_OWNED ProjectType = iota
	PROJECTTYPE_GRANTED
)
