package record

func cutString(str string, pos int) string {
	if len(str) <= pos {
		return str
	}
	return str[:pos-1]
}
