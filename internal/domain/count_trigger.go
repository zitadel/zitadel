package domain

//go:generate go tool github.com/dmarkham/enumer -type CountParentType -transform lower -trimprefix CountParentType -sql
type CountParentType int

const (
	CountParentTypeInstance CountParentType = iota
	CountParentTypeOrganization
)
