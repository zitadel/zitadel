package types

func (notify Notify) WithoutTemplate() error {
	return notify("", nil, "", false)
}
