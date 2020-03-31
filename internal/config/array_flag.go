package config

import (
	"flag"
	"strings"
)

var _ flag.Value = (*ArrayFlags)(nil)

//ArrayFlags implements the flag/Value interface
//allowing to set multiple string flags with the same name
type ArrayFlags []string

func (i *ArrayFlags) String() string {
	return strings.Join(*i, ";")
}

func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
