package handler

import (
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SqsHandler interface {
	ReciveMessage(message sqsTypes.Message) bool
}