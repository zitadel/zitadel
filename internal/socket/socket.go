package socket

import (
	"fmt"
	"io"
	"net"
)

const Path = "/tmp/zitadel.sock"

type SocketRequest byte

const (
	unknown SocketRequest = iota
	ReadinessQuery
)

type SocketResponse byte

const (
	unknownResponse SocketResponse = iota
	UnknownRequest
	True
	False
)

func (s SocketRequest) Request() (resp SocketResponse, err error) {
	conn, err := net.Dial("unix", Path)
	if err != nil {
		return resp, fmt.Errorf("dial error: %w", err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte{byte(s)})
	if err != nil {
		return resp, fmt.Errorf("write error: %w", err)
	}
	response, err := io.ReadAll(conn)
	if err != nil {
		return resp, fmt.Errorf("read error: %w", err)
	}
	if len(response) != 1 {
		return resp, fmt.Errorf("invalid response length")
	}
	return SocketResponse(response[0]), nil
}

type HandleFunc[T any] func(T, SocketRequest) (SocketResponse, error)

func respond[T any](conn net.Conn, handler HandleFunc[T], server T) error {
	defer conn.Close()
	buf := make([]byte, 1)
	_, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("could not read from socket: %v", err)
	}
	if len(buf) != 1 {
		return fmt.Errorf("invalid request length: %d", len(buf))
	}
	req := SocketRequest(buf[0])
	resp, err := handler(server, req)
	if err != nil {
		return fmt.Errorf("handler error: %w", err)
	}
	_, err = conn.Write([]byte{byte(resp)})
	if err != nil {
		return fmt.Errorf("could not write response: %v", err)
	}
	return nil
}
