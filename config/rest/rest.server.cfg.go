package config

import (
	errs "cqrsdemo/errors"
	restHandler "cqrsdemo/handler/rest"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

func NewRestServer(group *restHandler.GruopHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/hi", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello World!"))
		if err != nil {
			return
		}
	})

	r.Route("/file", func(fileRouter chi.Router) {
		fileRouter.Method("POST", "/upload", HandleMethod(group.FileRestHandler.ProcessUpload))
		fileRouter.Method("GET", "/list", HandleMethod(group.FileRestHandler.ProcessList))
		fileRouter.Method("GET", "/download", HandleMethod(group.FileRestHandler.ProcessDownloadFile))
	})

	r.Method("POST", "/commands", HandleMethod(group.CommandRestHandler.Process))

	return r
}

type HandleMethod func(w http.ResponseWriter, r *http.Request) error

func (h HandleMethod) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		slog.Error("error processing", "path", r.URL, "error", err.Error())
		var queryError errs.RequestError
		switch {
		case errors.As(err, &queryError):
			render.Status(r, queryError.StatusCode)
			render.JSON(w, r, queryError)
		default:
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{
				"error": err.Error(),
			})
		}
	}
}
