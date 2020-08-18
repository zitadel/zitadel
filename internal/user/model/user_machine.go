package model

type Machine struct {
	Name        string
	Description string
}

func (sa *Machine) IsValid() bool {
	return sa.Name != ""
}
