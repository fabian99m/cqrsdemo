package config

import (
	"cqrsdemo/adapter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client() adapter.S3Operations {
	S3Client := s3.NewFromConfig(loadAWSConfig(), func(o *s3.Options) {
		o.BaseEndpoint = aws.String("https://s3.localhost.localstack.cloud:4566")
	})

	return adapter.NewS3Actions(S3Client)
}
