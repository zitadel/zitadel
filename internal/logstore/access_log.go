package logstore

import (
	"net/http"
	"time"
)

type Protocol uint8

const (
	GRPC Protocol = iota
	HTTP
	// TODO: GRPC-Web?
	// TODO: HTTPS?
)

type AccessLogRecord struct {
	Timestamp       time.Time
	Protocol        Protocol
	RequestURL      string
	ResponseStatus  uint32
	RequestHeaders  http.Header
	ResponseHeaders http.Header
	InstanceID      string
	ProjectID       string
	RequestedDomain string
	RequestedHost   string
}
