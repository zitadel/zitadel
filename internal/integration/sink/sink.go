// Package sink provides a simple HTTP server where Zitadel can send HTTP based messages,
// which are then possible to be observed using observers on websockets.
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
	"path"
	"sync"
	"sync/atomic"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/zitadel/logging"
)

const (
	port       = "8081"
	listenAddr = "127.0.0.1:" + port
	host       = "localhost:" + port
)

//go:generate enumer -type Channel -trimprefix Channel -transform snake
type Channel int

const (
	ChannelMilestone Channel = iota
	ChannelQuota
)

// CallURL returns the full URL to the handler of a Channel.
func CallURL(ch Channel) string {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   rootPath(ch),
	}
	return u.String()
}

func rootPath(c Channel) string {
	return path.Join("/", c.String())
}

func subscribePath(c Channel) string {
	return path.Join("/", c.String(), "subscribe")
}

// Request is a message forwarded from the handler to Subscriptions.
type Request struct {
	Header http.Header
	Body   json.RawMessage
}

func ListenAndServe() error {
	router := chi.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(next http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(next, r)
		})
	})
	for _, ch := range ChannelValues() {
		fwd := &forwarder{
			channelID: ch,
			channels:  make(map[int64]chan<- *Request),
		}
		router.HandleFunc(rootPath(ch), fwd.receiveHandler)
		router.HandleFunc(subscribePath(ch), fwd.subscriptionHandler)
	}
	return http.ListenAndServe(listenAddr, router)
}

type forwarder struct {
	channelID Channel
	id        atomic.Int64
	mtx       sync.RWMutex
	channels  map[int64]chan<- *Request
	upgrader  websocket.Upgrader
}

func (c *forwarder) receiveHandler(w http.ResponseWriter, r *http.Request) {
	req := &Request{
		Header: r.Header.Clone(),
	}
	var err error
	req.Body, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	c.mtx.RLock()
	for _, reqChan := range c.channels {
		reqChan <- req
	}
	c.mtx.RUnlock()
	w.WriteHeader(http.StatusOK)
}

func (c *forwarder) subscriptionHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := c.upgrader.Upgrade(w, r, nil)
	logging.OnError(err).Error("websocket upgrade")
	if err != nil {
		return
	}
	done := readLoop(ws)

	id := c.id.Add(1)
	reqChannel := make(chan *Request, 100)

	c.mtx.Lock()
	c.channels[id] = reqChannel
	c.mtx.Unlock()

	logging.WithFields("id", id, "channel", c.channelID).Info("websocket opened")

	defer func() {
		c.mtx.Lock()
		delete(c.channels, id)
		c.mtx.Unlock()

		ws.Close()
		close(reqChannel)
	}()

	for {
		select {
		case err := <-done:
			logging.WithError(err).WithFields(logrus.Fields{"id": id, "channel": c.channelID}).Info("websocket closed")
			return
		case req := <-reqChannel:
			if err := ws.WriteJSON(req); err != nil {
				logging.WithError(err).WithFields(logrus.Fields{"id": id, "channel": c.channelID}).Error("websocket write json")
				return
			}
		}
	}
}

// readLoop makes sure we can receive close messages
func readLoop(ws *websocket.Conn) (done chan error) {
	done = make(chan error, 1)

	go func(done chan<- error) {
		for {
			_, _, err := ws.NextReader()
			if err != nil {
				done <- err
				break
			}
		}
		close(done)
	}(done)

	return done
}

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

// Recv returns the channel over which Requests are send.
func (s *Subscription) Recv() <-chan *Request {
	return s.reqChannel
}

func (s *Subscription) Close() error {
	s.closed.Store(true)
	return s.conn.Close()
}
