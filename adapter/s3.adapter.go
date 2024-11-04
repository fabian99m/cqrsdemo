package adapter

import (
	"bytes"
	"context"
	"errors"
	"io"
	"math/rand/v2"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gabriel-vasile/mimetype"
)

const objectKeyCharset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"

type s3Api interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type s3Actions struct {
	s3Client s3Api
}

func NewS3Actions(s3Client s3Api) *s3Actions {
	return &s3Actions{s3Client: s3Client}
}

func (r s3Actions) UploadFile(bucketName string, fileName string, file io.Reader) (string, error) {
	key := ramdomString(10) + "_" + fileName

	var buf bytes.Buffer
	mime, err := mimetype.DetectReader(io.TeeReader(file, &buf))
	if err != nil {
		return "", err
	}

	_, err = r.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      &bucketName,
		Key:         &key,
		Body:        io.MultiReader(&buf, file),
		ContentType: aws.String(mime.String()),
	})

	if err != nil {
		return "", err
	}

	return key, nil
}

type S3Pagination struct {
	MaxKeys int32
	Next    string
}

func (r s3Actions) ListFiles(bucketName string, pagination *S3Pagination) (*FileInfoResults, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:  &bucketName,
		MaxKeys: &pagination.MaxKeys,
	}

	if pagination.Next != "" {
		input.ContinuationToken = &pagination.Next
	}

	listObjectsOutput, err := r.s3Client.ListObjectsV2(context.Background(), input)
	if err != nil {
		return nil, err
	}

	fileInfoResults := FileInfoResults{
		All: !*listObjectsOutput.IsTruncated,
	}

	if listObjectsOutput.NextContinuationToken != nil {
		fileInfoResults.Next = *listObjectsOutput.NextContinuationToken
	}

	listFiles := make([]FileInfo, len(listObjectsOutput.Contents))
	for i, file := range listObjectsOutput.Contents {
		listFiles[i] = FileInfo{Name: *file.Key, Size: *file.Size}
	}

	fileInfoResults.FilesInfo = listFiles

	return &fileInfoResults, nil
}

func (r s3Actions) DownloadFile(bucketName string, key string) (*FileContent, error) {
	file, err := r.s3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})

	if err != nil {
		var NoSuchKey *types.NoSuchKey
		if errors.As(err, &NoSuchKey) {
			return nil, nil
		}

		return nil, err
	}

	return &FileContent{
		Name:          key,
		Body:          file.Body,
		ContentType:   *file.ContentType,
		ContentLength: *file.ContentLength,
	}, nil
}

func ramdomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = objectKeyCharset[rand.N(len(objectKeyCharset))]
	}
	return string(b)
}

type FileInfo struct {
	Name string
	Size int64
}

type FileInfoResults struct {
	FilesInfo []FileInfo
	Next      string
	All       bool
}

type FileContent struct {
	Name          string
	ContentType   string
	ContentLength int64
	Body          io.ReadCloser
}
