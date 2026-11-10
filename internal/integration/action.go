package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type server struct {
	server *httptest.Server
	mu     sync.Mutex
	called int
}

func (s *server) URL() string {
	return s.server.URL
}

func (s *server) Close() {
	s.server.Close()
}

func (s *server) Called() int {
	s.mu.Lock()
	called := s.called
	s.mu.Unlock()
	return called
}

func (s *server) Increase() {
	s.mu.Lock()
	s.called++
	s.mu.Unlock()
}

func (s *server) ResetCalled() {
	s.mu.Lock()
	s.called = 0
	s.mu.Unlock()
}

func TestServerCall(
	reqBody interface{},
	sleep time.Duration,
	statusCode int,
	respBody interface{},
) (url string, closeF func(), calledF func() int, resetCalledF func()) {
	server := &server{
		called: 0,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		server.Increase()
		if reqBody != nil {
			data, err := json.Marshal(reqBody)
			if err != nil {
				http.Error(w, "error, marshall: "+err.Error(), http.StatusInternalServerError)
				return
			}
			sentBody, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "error, read body: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if !reflect.DeepEqual(data, sentBody) {
				http.Error(w, "error, equal:\n"+string(data)+"\nsent:\n"+string(sentBody), http.StatusInternalServerError)
				return
			}
		}
		if statusCode != http.StatusOK {
			http.Error(w, "error, statusCode", statusCode)
			return
		}

		time.Sleep(sleep)

		if respBody != nil {
			w.Header().Set("Content-Type", "application/json")
			resp, err := json.Marshal(respBody)
			if err != nil {
				http.Error(w, "error", http.StatusInternalServerError)
				return
			}
			if _, err := io.Writer.Write(w, resp); err != nil {
				http.Error(w, "error", http.StatusInternalServerError)
				return
			}
		} else {
			if _, err := io.WriteString(w, "finished successfully"); err != nil {
				http.Error(w, "error", http.StatusInternalServerError)
				return
			}
		}
	}

	server.server = httptest.NewServer(http.HandlerFunc(handler))
	return server.URL(), server.Close, server.Called, server.ResetCalled
}

func TestServerCallProto(
	reqBody interface{},
	sleep time.Duration,
	statusCode int,
	respBody proto.Message,
) (url string, closeF func(), calledF func() int, resetCalledF func()) {
	server := &server{
		called: 0,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		server.Increase()
		if reqBody != nil {
			data, err := json.Marshal(reqBody)
			if err != nil {
				http.Error(w, "error, marshall: "+err.Error(), http.StatusInternalServerError)
				return
			}
			sentBody, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "error, read body: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if !reflect.DeepEqual(data, sentBody) {
				http.Error(w, "error, equal:\n"+string(data)+"\nsent:\n"+string(sentBody), http.StatusInternalServerError)
				return
			}
		}
		if statusCode != http.StatusOK {
			http.Error(w, "error, statusCode", statusCode)
			return
		}

		time.Sleep(sleep)

		if respBody != nil {
			w.Header().Set("Content-Type", "application/json")
			resp, err := protojson.Marshal(respBody)
			if err != nil {
				http.Error(w, "error", http.StatusInternalServerError)
				return
			}
			if _, err := io.Writer.Write(w, resp); err != nil {
				http.Error(w, "error", http.StatusInternalServerError)
				return
			}
		} else {
			if _, err := io.WriteString(w, "finished successfully"); err != nil {
				http.Error(w, "error", http.StatusInternalServerError)
				return
			}
		}
	}

	server.server = httptest.NewServer(http.HandlerFunc(handler))
	return server.URL(), server.Close, server.Called, server.ResetCalled
}
