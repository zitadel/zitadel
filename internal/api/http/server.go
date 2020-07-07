package http

import (
	"context"
	"net/http"

	"github.com/caos/logging"
)

func Serve(ctx context.Context, handler http.Handler, port, servername string) {
	server := &http.Server{
		Handler: handler,
	}

	listener := CreateListener(port)

	go func() {
		<-ctx.Done()
		err := server.Shutdown(ctx)
		logging.LogWithFields("HTTP-m7kBlq", "name", servername).OnError(err).Warnf("error during graceful shutdown of http server (%s)", servername)
	}()

	go func() {
		err := server.Serve(listener)
		logging.LogWithFields("HTTP-tBHR60", "name", servername).OnError(err).Panicf("http serve (%s) failed", servername)
	}()
	logging.LogWithFields("HTTP-KHh0Cb", "name", servername, "port", port).Infof("http server (%s) is listening", servername)
}

func RegisterHandler(mux *http.ServeMux, prefix string, handler http.Handler) {
	mux.Handle(prefix+"/", http.StripPrefix(prefix, handler))
}
