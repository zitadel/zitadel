package config

import (
	"flag"
	"strings"
)

var _ flag.Value = (*ArrayFlags)(nil)

//ArrayFlags implements the flag/Value interface
//allowing to set multiple string flags with the same name
type ArrayFlags struct {
	defaultValues []string
	values        []string
}

func NewArrayFlags(defaults ...string) *ArrayFlags {
	return &ArrayFlags{
		defaultValues: defaults,
	}
}

func (i *ArrayFlags) Values() []string {
	if len(i.values) == 0 {
		return i.defaultValues
	}
	return i.values
}

func (i *ArrayFlags) String() string {
	return strings.Join(i.Values(), ";")
}

func (i *ArrayFlags) Set(value string) error {
	i.values = append(i.values, value)
	return nil
}
