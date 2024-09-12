//go:build integration

package sink

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/zitadel/logging"
)

// Request is a message forwarded from the handler to [Subscription]s.
type Request struct {
	Header http.Header
	Body   json.RawMessage
}

// Subscription is a websocket client to which [Request]s are forwarded by the server.
type Subscription struct {
	conn       *websocket.Conn
	closed     atomic.Bool
	reqChannel chan *Request
}

// Subscribe to a channel.
// The subscription forwards all requests it received on the channel's
// handler, after Subscribe has returned.
// Multiple subscription may be active on a single channel.
// Each request is always forwarded to each Subscription.
// Close must be called to cleanup up the Subscription's channel and go routine.
func Subscribe(ctx context.Context, ch Channel) *Subscription {
	u := url.URL{
		Scheme: "ws",
		Host:   listenAddr,
		Path:   subscribePath(ch),
	}
	conn, resp, err := websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		if resp != nil {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			err = fmt.Errorf("subscribe: %w, status: %s, body: %s", err, resp.Status, body)
		}
		panic(err)
	}

	sub := &Subscription{
		conn:       conn,
		reqChannel: make(chan *Request, 10),
	}
	go sub.readToChan()
	return sub
}

func (s *Subscription) readToChan() {
	for {
		if s.closed.Load() {
			break
		}
		req := new(Request)
		if err := s.conn.ReadJSON(req); err != nil {
			opErr := new(net.OpError)
			if errors.As(err, &opErr) {
				break
			}
			logging.WithError(err).Error("subscription read")
			break
		}
		s.reqChannel <- req
	}
	close(s.reqChannel)
}

// Recv returns the channel over which [Request]s are send.
func (s *Subscription) Recv() <-chan *Request {
	return s.reqChannel
}

func (s *Subscription) Close() error {
	s.closed.Store(true)
	return s.conn.Close()
}
