package adapter

import (
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	mock "github.com/ovechkin-dm/mockio/mock"
	"github.com/stretchr/testify/assert"
)

func TestDownloadFile(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name       string
		fileName   string
		out        *s3.GetObjectOutput
		errService error
		errRes     error
	}{
		{
			name:     "download success",
			fileName: "test.text",
			out: &s3.GetObjectOutput{
				Body:          io.NopCloser(strings.NewReader("hello world")),
				ContentType:   aws.String("text"),
				ContentLength: aws.Int64(12),
			},
			errService: nil,
			errRes:     nil,
		},
		{
			name:       "key not found",
			fileName:   "test.text",
			out:        nil,
			errService: &types.NoSuchKey{},
			errRes:     nil,
		},
		{
			name:       "key not found",
			fileName:   "test.text",
			out:        nil,
			errService: &types.NoSuchBucket{},
			errRes:     &types.NoSuchBucket{},
		},
	}

	mock.SetUp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			apiMock := mock.Mock[s3Api]()
			mock.When(apiMock.GetObject(mock.AnyContext(), mock.Any[*s3.GetObjectInput]())).ThenReturn(tt.out, tt.errService)

			underTest := NewS3Actions(apiMock)
			_, err := underTest.DownloadFile("testbucket", tt.fileName)
			assert.Equal(subtest, tt.errRes, err)
		})
	}
}

func TestUploadFile(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name     string
		fileName string
		file     io.Reader
		err      error
	}{
		{
			name:     "updload success",
			fileName: "test.text",
			file:     strings.NewReader("hola mundo"),
			err:      nil,
		},
		{
			name:     "updload error",
			fileName: "test.text",
			file:     strings.NewReader("hola mundo"),
			err:      &types.NoSuchBucket{},
		},
	}

	mock.SetUp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			apiMock := mock.Mock[s3Api]()
			mock.When(apiMock.PutObject(mock.AnyContext(), mock.Any[*s3.PutObjectInput]())).ThenReturn(&s3.PutObjectOutput{}, tt.err)
			underTest := NewS3Actions(apiMock)

			res, err := underTest.UploadFile("testbucket", tt.fileName, tt.file)
			assert.Equal(subtest, tt.err, err)
			if err == nil {
				assert.True(subtest, len(res) == 11+len(tt.fileName))
			}
		})
	}
}

func TestListFiles(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name     string
		input    *S3Pagination
		response *s3.ListObjectsV2Output
		err      error
	}{
		{
			name: "success 1",
			input: &S3Pagination{
				MaxKeys: 1,
				Next:    "",
			},
			response: &s3.ListObjectsV2Output{
				IsTruncated:           aws.Bool(false),
				NextContinuationToken: aws.String("avc"),
				Contents: []types.Object{
					{
						Key:  aws.String("test"),
						Size: aws.Int64(123),
					},
				},
			},
			err: nil,
		},
		{
			name: "error test",
			input: &S3Pagination{
				MaxKeys: 1,
				Next:    "test",
			},
			response: nil,
			err:      &types.NoSuchBucket{},
		},
	}

	mock.SetUp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			apiMock := mock.Mock[s3Api]()
			mock.When(apiMock.ListObjectsV2(mock.AnyContext(), mock.Any[*s3.ListObjectsV2Input]())).ThenReturn(tt.response, tt.err)
			underTest := NewS3Actions(apiMock)

			res, err := underTest.ListFiles("testbucket", tt.input)
			assert.Equal(subtest, tt.err, err)
			if err == nil {
				assert.Equal(subtest, len(tt.response.Contents), len(res.FilesInfo))
			}
		})
	}
}
