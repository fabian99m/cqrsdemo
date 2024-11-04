package adapter

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SnsOperations interface {
	Publish(topicArn string, message *string, typeMessage string) (string, error)
}

type SqsOperations interface {
	GetMessages(queueUrl string, maxMessages int32, waitTime int32) ([]types.Message, error)
	GetQueueUrl(queueName string) (string, error)
	DeleteMessage(queueUrl string, rangeeceiptHandle string) error
}

type S3Operations interface {
	UploadFile(bucketName string, fileName string, file io.Reader) (string, error)
	ListFiles(bucketName string, pagination *S3Pagination) (*FileInfoResults, error)
	DownloadFile(bucketName string, key string) (*FileContent, error)
}