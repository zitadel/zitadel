package model

//Deprecated: Enum is useless, better use normal enums, because we rarely need string value
type Enum interface {
	String() string
}
