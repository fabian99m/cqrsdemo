package config

import (
	"github.com/fabian99m/cqrsdemo/adapter"
	
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NewSqsClient() adapter.SqsOperations {
	sqsClient := sqs.NewFromConfig(loadAWSConfig(), func(o *sqs.Options) {
		o.BaseEndpoint = aws.String("http://localhost:4566")
	})

	return adapter.NewSqsActions(sqsClient)
}

func NewSnsClient() adapter.SnsOperations {
	snsClient := sns.NewFromConfig(loadAWSConfig(), func(o *sns.Options) {
		o.BaseEndpoint = aws.String("http://localhost:4566")
	})

	return adapter.NewSnsActions(snsClient)
}
