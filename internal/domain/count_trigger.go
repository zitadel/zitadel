package domain

//go:generate enumer -type CountParentType -transform lower -trimprefix CountParentType -sql
type CountParentType int

const (
	CountParentTypeInstance CountParentType = iota
	CountParentTypeOrganization
)
