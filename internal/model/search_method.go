package model

// code below could be generated
type SearchMethod Enum

var methods = []string{"Equals", "StartsWith", "Contains", "EqualsCaseSensitive", "StartsWithCaseSensitive", "ContainsCaseSensitive"}

type method int32

func (s method) String() string {
	return methods[s]
}

const (
	Equals method = iota
	StartsWith
	Contains
	EqualsCaseSensitive
	StartsWithCaseSensitive
	ContainsCaseSensitive
)

func SearchMethodToInt(s SearchMethod) int32 {
	return int32(s.(method))
}

func SearchMethodFromInt(index int32) SearchMethod {
	return method(index)
}
