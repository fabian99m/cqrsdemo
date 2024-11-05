package config

import (
	"errors"
	errs "github.com/fabian99m/cqrsdemo/errors"
	restHandler "github.com/fabian99m/cqrsdemo/handler/rest"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func NewFileRouter(r *chi.Mux, file, handler *restHandler.FileHandler) {
	r.Route("/file", func(fileRouter chi.Router) {
		fileRouter.Post("/upload", HandleMethod(handler.ProcessUpload).ServeHTTP)
		fileRouter.Get("/list", HandleMethod(handler.ProcessList).ServeHTTP)
		fileRouter.Get("/download", HandleMethod(handler.ProcessDownloadFile).ServeHTTP)
	})
}

func NewCommandRouter(r *chi.Mux, file, handler *restHandler.CommandRestHandler) {
	r.Route("/commands", func(fileRouter chi.Router) {
		fileRouter.Post("/", HandleMethod(handler.Process).ServeHTTP)
	})
}

func NewBaseRestServer(group *restHandler.GroupHandler) http.Handler {
	r := chi.NewRouter()

	r.Route("/file", func(fileRouter chi.Router) {
		handler := group.FileRestHandler

		fileRouter.Post("/upload", HandleMethod(handler.ProcessUpload).ServeHTTP)
		fileRouter.Get("/list", HandleMethod(handler.ProcessList).ServeHTTP)
		fileRouter.Get("/download", HandleMethod(handler.ProcessDownloadFile).ServeHTTP)
	})

	r.Route("/commands", func(fileRouter chi.Router) {
		fileRouter.Post("/", HandleMethod(group.CommandRestHandler.Process).ServeHTTP)
	})

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
