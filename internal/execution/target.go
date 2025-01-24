package execution

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"time"
)

type server struct {
	server *httptest.Server
	called bool
}

func (s *server) URL() string {
	return s.server.URL
}

func (s *server) Close() {
	s.server.Close()
}

func (s *server) Called() bool {
	return s.called
}

func TestServerCall(
	reqBody interface{},
	sleep time.Duration,
	statusCode int,
	respBody interface{},
) (string, func(), func() bool) {
	server := &server{
		called: false,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		server.called = true
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
			if _, err := io.WriteString(w, string(resp)); err != nil {
				http.Error(w, "error", http.StatusInternalServerError)
				return
			}
		} else {
			io.WriteString(w, "finished successfully")
		}
	}

	server.server = httptest.NewServer(http.HandlerFunc(handler))
	return server.URL(), server.Close, server.Called
}
