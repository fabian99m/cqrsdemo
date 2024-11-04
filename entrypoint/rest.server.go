package entrypoint

import (
	"log/slog"
	"net/http"
)

func RestServer(httpHandler http.Handler) {
	slog.Info("restServer start", "port", 9090)
	err := http.ListenAndServe("localhost:9090", httpHandler)
	if err != nil {
		panic(err)
	}
}
