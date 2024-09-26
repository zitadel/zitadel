package socket

import (
	"fmt"
	"net"
	"os"

	"github.com/zitadel/logging"
)

func Listen[T any](handleFunc HandleFunc[T]) (chan<- T, func() error, error) {
	listener, err := listen()
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("cannot start socket listener: %w", err)
	}
	serverStartedUp := make(chan T)
	go acceptSocketConnections[T](listener, serverStartedUp, handleFunc)
	return serverStartedUp, listener.Close, nil
}

func ListenAndIgnore() (func() error, error) {
	listener, err := listen()
	if err != nil {
		return func() error { return nil }, fmt.Errorf("cannot start socket listener: %w", err)
	}
	go acceptSocketConnections[any](listener, make(chan any), nil)
	return listener.Close, nil
}

func listen() (net.Listener, error) {
	if err := os.Remove(Path); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("cannot remove socket file: %w", err)
	}
	return net.Listen("unix", Path)
}

type Handler func(request SocketRequest) (SocketResponse, error)

func acceptSocketConnections[T any](listener net.Listener, startupDone <-chan T, handle HandleFunc[T]) {
	var server T
	for {
		conn, err := listener.Accept()
		if err != nil {
			logging.Errorf("accept socket error: %v", err)
			continue
		}
		if server == nil {
			server = <-startupDone
		}
		if handle == nil {
			conn.Close()
			return
		}
		go func() {
			if handleErr := respond(conn, handle, server); handleErr != nil {
				logging.Errorf("socket handle error: %v", handleErr)
			}
		}()
	}
}
