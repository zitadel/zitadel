//go:build integration

package net

import (
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

func ListenConfig() *net.ListenConfig {
	return &net.ListenConfig{
		Control: reusePort,
	}
}

func reusePort(network, address string, conn syscall.RawConn) error {
	return conn.Control(func(descriptor uintptr) {
		err := syscall.SetsockoptInt(int(descriptor), syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1)
		if err != nil {
			panic(err)
		}
	})
}
