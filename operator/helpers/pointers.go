package helpers

func PointerInt32(value int32) *int32 {
	pointer := value
	return &pointer
}

func PointerInt64(value int64) *int64 {
	pointer := value
	return &pointer
}

func PointerBool(value bool) *bool {
	pointer := value
	return &pointer
}
