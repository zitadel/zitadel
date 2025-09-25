package sink

//go:generate go tool enumer -type Channel -trimprefix Channel -transform snake
type Channel int

const (
	ChannelMilestone Channel = iota
	ChannelQuota
)
