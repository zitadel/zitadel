package keys

type MachineKeyType int32

const (
	MachineKeyTypeNONE = iota
	MachineKeyTypeJSON

	keyCount
)

func (f MachineKeyType) Valid() bool {
	return f >= 0 && f < keyCount
}
