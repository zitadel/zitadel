package cache

/* TBD, where would we put these definitions? */

type InstanceIndex int16

//go:generate enumer -type InstanceIndex -trimprefix InstanceIndex
const (
	InstanceIndexByID InstanceIndex = iota
	InstanceIndexByHost
)
