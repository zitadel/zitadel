package domain

type SearchMethod int32

const (
	SearchMethodEquals SearchMethod = iota
	SearchMethodStartsWith
	SearchMethodContains
	SearchMethodEqualsIgnoreCase
	SearchMethodStartsWithIgnoreCase
	SearchMethodContainsIgnoreCase
	SearchMethodNotEquals
	SearchMethodGreaterThan
	SearchMethodLessThan
	SearchMethodIsOneOf
	SearchMethodListContains
	SearchMethodEndsWith
	SearchMethodEndsWithIgnoreCase
)
