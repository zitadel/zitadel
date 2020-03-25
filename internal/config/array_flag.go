package config

import "strings"

type ArrayFlags []string

func (i *ArrayFlags) String() string {
	return strings.Join(*i, ";")
}

func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
