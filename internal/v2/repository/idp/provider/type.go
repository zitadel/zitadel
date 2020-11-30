package provider

type Type int8

const (
	TypeSystem Type = iota
	TypeOrg

	typeCount
)

func (f Type) Valid() bool {
	return f >= 0 && f < typeCount
}
