package adapter

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type sqsApi interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	GetQueueUrl(ctx context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

type sqsActions struct {
	sqsClient sqsApi
}

func NewSqsActions(sqsClient sqsApi) SqsOperations{
	return &sqsActions{sqsClient: sqsClient}
}

func (r *sqsActions) GetMessages(queueUrl string, maxMessages int32, waitTime int32) ([]types.Message, error) {
	
	result, err := r.sqsClient.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(queueUrl),
		MaxNumberOfMessages:   maxMessages,
		WaitTimeSeconds:       waitTime,
		MessageAttributeNames: []string{".*"},
	})

	if err != nil {
		slog.Error("couldn't get messages from queue", "error", err)
		return nil, err
	}

	return result.Messages, nil
}

func (r *sqsActions) GetQueueUrl(queueName string) (string, error) {
	resp, err := r.sqsClient.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		slog.Error("couldn't get queue URL", "queueName", queueName, "error", err)
		return "", err
	}

	return *resp.QueueUrl, err
}

func (r *sqsActions) DeleteMessage(queueUrl string, receiptHandle string) error {
	slog.Info("deleteMessages start...")

	_, err := r.sqsClient.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: aws.String(receiptHandle),
	})

	if err != nil {
		slog.Error("couldn't delete message from queue", "error", err)
		return err
	}

	return nil
}
