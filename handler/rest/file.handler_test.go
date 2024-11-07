package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	adr "github.com/fabian99m/cqrsdemo/adapter"
	e "github.com/fabian99m/cqrsdemo/errors"
	"github.com/fabian99m/cqrsdemo/model"
	"github.com/fabian99m/cqrsdemo/util"

	queryBuilder "github.com/google/go-querystring/query"
	mock "github.com/ovechkin-dm/mockio/mock"
	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name       string
		serviceOut string
		errService error
		maxSize    int32
		upload     bool
		errorCode  int
		statusCode int
	}{
		{
			name:       "file not found",
			serviceOut: "filekey",
			upload:     false,
			maxSize:    10,
			errorCode:  e.FileNotFound.Code,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "maxSize",
			serviceOut: "filekey",
			upload:     true,
			maxSize:    0,
			errorCode:  e.FileSizeInvalid.Code,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "sucesss",
			serviceOut: "filekey",
			upload:     true,
			maxSize:    10,
			statusCode: http.StatusOK,
		},
		{
			name:       "service error",
			upload:     true,
			maxSize:    10,
			errService: fmt.Errorf("aws s3 error"),
			errorCode:  e.GenericError.Code,
			statusCode: http.StatusInternalServerError,
		},
	}

	mock.SetUp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			t.Log("running test", tt.name)

			mockS3 := mock.Mock[adr.S3Operations]()
			mock.When(mockS3.UploadFile(mock.AnyString(), mock.AnyString(), mock.Any[io.Reader]())).ThenReturn(tt.serviceOut, tt.errService)

			underTest := NewFileHandler(mockS3, &model.BucketProps{Name: "Test", MaxSize: tt.maxSize})

			recorder := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/file/upload", nil)
			if tt.upload {
				formFile(r, subtest)
			}

			err := underTest.ProcessUpload(recorder, r)

			if err != nil {
				var reqErr e.RequestError
				if errors.As(err, &reqErr) {
					assert.Equal(subtest, tt.statusCode, reqErr.StatusCode)
					assert.Equal(subtest, tt.errorCode, reqErr.Status.Code)
				}
			} else {
				assert.Equal(subtest, tt.statusCode, recorder.Code)
			}
		})
	}
}

func TestDownloadFile(t *testing.T) {
	t.Parallel()

	fileContent := "hello world"
	var tests = []struct {
		name       string
		key        string
		serviceOut *adr.FileContent
		errService error
		errorCode  int
		statusCode int
	}{
		{
			name:       "empty key",
			key:        "",
			errorCode:  e.InvalidParams.Code,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "service error",
			key:        "test.txt",
			serviceOut: nil,
			errService: fmt.Errorf("service error"),
			errorCode:  e.GenericError.Code,
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "ok",
			key:  "test.txt",
			serviceOut: &adr.FileContent{
				Name:          "asaaasas_text.key",
				ContentType:   "text/plain",
				ContentLength: 123,
				Body:          io.NopCloser(strings.NewReader(fileContent)),
			},
			statusCode: http.StatusOK,
		},
		{
			name:       "file not found",
			key:        "test.txt",
			serviceOut: nil,
			errService: nil,
			statusCode: http.StatusNotFound,
		},
	}

	mock.SetUp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			t.Log("running test: ", tt.name)

			mockS3 := mock.Mock[adr.S3Operations]()
			mock.When(mockS3.DownloadFile(mock.AnyString(), mock.AnyString())).ThenReturn(tt.serviceOut, tt.errService)

			underTest := NewFileHandler(mockS3, &model.BucketProps{Name: "Test"})

			recorder := httptest.NewRecorder()
			err := underTest.ProcessDownloadFile(recorder, httptest.NewRequest("GET", "/file/test?key="+tt.key, nil))

			if err != nil {
				var reqErr e.RequestError
				if errors.As(err, &reqErr) {
					assert.Equal(subtest, tt.statusCode, reqErr.StatusCode)
				}
			} else {
				assert.Equal(subtest, tt.statusCode, recorder.Code)
				assert.Equal(subtest, fileContent, recorder.Body.String())
			}
		})
	}
}

func TestProcessList(t *testing.T) {
	t.Parallel()

	type queryList struct {
		Next    string `url:"next"`
		MaxKeys string `url:"maxKeys"`
	}

	var tests = []struct {
		name       string
		query      queryList
		out        *adr.FileInfoResults
		errService error
		errorCode  int
		statusCode int
	}{
		{
			name: "list success",
			query: queryList{
				Next:    "1213",
				MaxKeys: "5",
			},
			out: &adr.FileInfoResults{
				FilesInfo: []adr.FileInfo{
					{
						Name: "test.txt",
						Size: 123,
					},
				},
			},
			statusCode: http.StatusOK,
		},
		{
			name: "list badrequest",
			query: queryList{
				Next:    "1213",
				MaxKeys: "abc",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "list service error",
			query: queryList{
				Next:    "1213",
				MaxKeys: "1",
			},
			errService: fmt.Errorf("service error"),
			statusCode: http.StatusInternalServerError,
		},
	}

	mock.SetUp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			t.Log("running test: ", tt.name)

			mockS3 := mock.Mock[adr.S3Operations]()
			mock.When(mockS3.ListFiles(mock.AnyString(), mock.NotNil[*adr.S3Pagination]())).ThenReturn(tt.out, tt.errService)

			underTest := NewFileHandler(mockS3, &model.BucketProps{Name: "Test"})

			v, err := queryBuilder.Values(tt.query)
			if err != nil {
				subtest.Error(err)
			}

			recorder := httptest.NewRecorder()
			err = underTest.ProcessList(recorder, httptest.NewRequest("GET", "/file/test?"+v.Encode(), nil))

			if err != nil {
				var reqErr e.RequestError
				if errors.As(err, &reqErr) {
					assert.Equal(subtest, tt.statusCode, reqErr.StatusCode)
				}
			} else {
				response, err := util.UnmarshalTo[adr.FileInfoResults](recorder.Body.Bytes())
				if err != nil {
					subtest.Error(err)
				}

				assert.Equal(subtest, tt.statusCode, recorder.Code)
				assert.Equal(subtest, len(tt.out.FilesInfo), len(response.FilesInfo))
			}
		})
	}
}

func formFile(r *http.Request, t *testing.T) {
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Error(err)
	}
	defer writer.Close()

	_, err = io.Copy(part, strings.NewReader("hello world"))
	if err != nil {
		t.Error(err)
	}

	r.Body = io.NopCloser(body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
}
