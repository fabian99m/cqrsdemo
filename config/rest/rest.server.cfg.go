package config

import (
	"errors"
	"log/slog"
	"net/http"

	errs "github.com/fabian99m/cqrsdemo/errors"
	restHandler "github.com/fabian99m/cqrsdemo/handler/rest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func NewFileRouter(r *chi.Mux, handler *restHandler.FileHandler) {
	r.Route("/file", func(fileRouter chi.Router) {
		fileRouter.Post("/upload", HandleMethod(handler.ProcessUpload).ServeHTTP)
		fileRouter.Get("/list", HandleMethod(handler.ProcessList).ServeHTTP)
		fileRouter.Get("/download", HandleMethod(handler.ProcessDownloadFile).ServeHTTP)
	})
}

func NewCommandRouter(r *chi.Mux, handler *restHandler.CommandRestHandler) {
	r.Route("/commands", func(commandRouter chi.Router) {
		commandRouter.Post("/", HandleMethod(handler.Process).ServeHTTP)
	})
}

func NewBaseRestServer(group *restHandler.GroupHandler) http.Handler {
	r := chi.NewRouter()
	NewFileRouter(r, &group.FileRestHandler)
	NewCommandRouter(r, &group.CommandRestHandler)
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
