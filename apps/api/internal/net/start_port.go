//go:build !integration

package net

import (
	"net"
)

func ListenConfig() *net.ListenConfig {
	return &net.ListenConfig{}
}
