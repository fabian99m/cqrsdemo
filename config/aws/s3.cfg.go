package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client() *s3.Client {
	s3Client := s3.NewFromConfig(loadAWSConfig(), func(o *s3.Options) {
		o.BaseEndpoint = aws.String("https://s3.localhost.localstack.cloud:4566")
	})

	return s3Client
}
