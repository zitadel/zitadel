package sink

//go:generate enumer -type Channel -trimprefix Channel -transform snake
type Channel int

const (
	ChannelMilestone Channel = iota
	ChannelQuota
)
