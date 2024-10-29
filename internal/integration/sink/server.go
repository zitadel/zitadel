//go:build integration

package sink

import (
	"errors"
	"io"
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

// CallURL returns the full URL to the handler of a [Channel].
func CallURL(ch Channel) string {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   rootPath(ch),
	}
	return u.String()
}

// StartServer starts a simple HTTP server on localhost:8081
// ZITADEL can use the server to send HTTP requests which can be
// used to validate tests through [Subscribe]rs.
// For each [Channel] a route is registered on http://localhost:8081/<channel_name>.
// The route must be used to send the HTTP request to be validated.
// [CallURL] can be used to obtain the full URL for a given Channel.
//
// This function is only active when the `integration` build tag is enabled
func StartServer() (close func()) {
	router := chi.NewRouter()
	for _, ch := range ChannelValues() {
		fwd := &forwarder{
			channelID:   ch,
			subscribers: make(map[int64]chan<- *Request),
		}
		router.HandleFunc(rootPath(ch), fwd.receiveHandler)
		router.HandleFunc(subscribePath(ch), fwd.subscriptionHandler)
	}
	s := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}

	logging.WithFields("listen_addr", listenAddr).Warn("!!!! A sink server is started which may expose sensitive data on a public endpoint. Make sure the `integration` build tag is disabled for production builds. !!!!")
	go func() {
		err := s.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			logging.WithError(err).Fatal("sink server")
		}
	}()
	return func() {
		logging.OnError(s.Close()).Error("sink server")
	}
}

func rootPath(c Channel) string {
	return path.Join("/", c.String())
}

func subscribePath(c Channel) string {
	return path.Join("/", c.String(), "subscribe")
}

// forwarder handles incoming HTTP requests from ZITADEL and
// forwards them to all subscribed web sockets.
type forwarder struct {
	channelID   Channel
	id          atomic.Int64
	mtx         sync.RWMutex
	subscribers map[int64]chan<- *Request
	upgrader    websocket.Upgrader
}

// receiveHandler receives a simple HTTP for a single [Channel]
// and forwards them on all active subscribers of that Channel.
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
	for _, reqChan := range c.subscribers {
		reqChan <- req
	}
	c.mtx.RUnlock()
	w.WriteHeader(http.StatusOK)
}

// subscriptionHandler upgrades HTTP request to a websocket connection for subscribers.
// All received HTTP requests on a subscriber's channel are send on the websocket to the client.
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
	c.subscribers[id] = reqChannel
	c.mtx.Unlock()

	logging.WithFields("id", id, "channel", c.channelID).Info("websocket opened")

	defer func() {
		c.mtx.Lock()
		delete(c.subscribers, id)
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
