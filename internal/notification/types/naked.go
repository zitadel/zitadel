package types

func (notify Notify) Naked() error {
	return notify("", nil, "", false)
}
