package handler

import (
	"github.com/fabian99m/cqrsdemo/adapter"
	e "github.com/fabian99m/cqrsdemo/errors"
	"github.com/fabian99m/cqrsdemo/model"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
)

type FileHandler struct {
	s3Operations adapter.S3Operations
	props        *model.BucketProps
}

func NewFileHandler(s3Action adapter.S3Operations, props *model.BucketProps) *FileHandler {
	return &FileHandler{
		s3Operations: s3Action, props: props,
	}
}

func (fh *FileHandler) ProcessList(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()

	next := query.Get("next")
	maxkeys, err := strconv.Atoi(query.Get("maxKeys"))

	if err != nil {
		return e.RequestError{
			StatusCode: http.StatusBadRequest, Status: e.GenericError.Fmt(err),
		}
	}

	ls, err := fh.s3Operations.ListFiles(fh.props.Name, &adapter.S3Pagination{
		MaxKeys: int32(maxkeys),
		Next:    next,
	})

	if err != nil {
		return e.RequestError{
			StatusCode: http.StatusInternalServerError, Status: e.GenericError.Fmt(err),
		}
	}

	render.JSON(w, r, ls)

	return nil
}

func (fh *FileHandler) ProcessDownloadFile(w http.ResponseWriter, r *http.Request) error {
	key := r.URL.Query().Get("key")

	if strings.TrimSpace(key) == "" {
		return e.RequestError{
			StatusCode: http.StatusBadRequest, Status: e.ParamsNotFound.Fmt("key"),
		}
	}

	file, err := fh.s3Operations.DownloadFile(fh.props.Name, key)
	if err != nil {
		return e.RequestError{
			StatusCode: http.StatusInternalServerError, Status: e.GenericError.Fmt(err),
		}
	}

	if file == nil {
		return e.RequestError{
			StatusCode: http.StatusNotFound, Status: e.FileNotFound,
		}
	}
	defer file.Body.Close()

	w.Header().Set("Content-Type", file.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", key))

	_, err = io.Copy(w, file.Body)
	if err != nil {
		return e.RequestError{
			StatusCode: http.StatusInternalServerError, Status: e.GenericError.Fmt(err),
		}
	}

	return nil
}

func (fh *FileHandler) ProcessUpload(w http.ResponseWriter, r *http.Request) error {
	file, header, err := r.FormFile("file")
	if err != nil {
		return e.RequestError{
			StatusCode: http.StatusBadRequest, Status: e.FileNotFound,
		}
	}

	sizeMb := float64(header.Size) / (1024 * 1024)
	slog.Info("file info", "sizeMb in mb", sizeMb, "name", header.Filename)

	if sizeMb > float64(fh.props.MaxSize) {
		return e.RequestError{
			StatusCode: http.StatusBadRequest, Status: e.FileSizeInvalid.Fmt(sizeMb),
		}
	}

	id, err := fh.s3Operations.UploadFile(fh.props.Name, header.Filename, file)
	if err != nil {
		return e.RequestError{
			StatusCode: http.StatusInternalServerError, Status: e.GenericError.Fmt(err),
		}
	}

	render.JSON(w, r, map[string]string{
		"fileId": id,
		"name":   header.Filename,
		"sizeMb": fmt.Sprintf("%.2f mb", sizeMb),
	})

	return nil
}
