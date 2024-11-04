package adapter

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/smithy-go/middleware"
	mock "github.com/ovechkin-dm/mockio/mock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteMessage(t *testing.T) {
	var tests = []struct {
		name   string
		out    *sqs.DeleteMessageOutput
		errRes error
	}{
		{
			name: "DeleteMessage suceess",
			out: &sqs.DeleteMessageOutput{
				ResultMetadata: middleware.Metadata{},
			},
			errRes: nil,
		},
		{
			name:   "DeleteMessage suceess",
			out:    nil,
			errRes: fmt.Errorf("aws error"),
		},
	}

	mock.SetUp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			apiMock := mock.Mock[sqsApi]()
			mock.When(apiMock.DeleteMessage(mock.AnyContext(), mock.Any[*sqs.DeleteMessageInput]())).ThenReturn(tt.out, tt.errRes)

			underTest := NewSqsActions(apiMock)

			err := underTest.DeleteMessage("queueName", "recipt")
			assert.Equal(subtest, tt.errRes, err)
		})
	}
}

func TestGetQueueUrl(t *testing.T) {
	var tests = []struct {
		name   string
		out    *sqs.GetQueueUrlOutput
		errRes error
	}{
		{
			name: "GetQueueUrl succes",
			out: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("url"),
			},
			errRes: nil,
		},
		{
			name:   "GetQueueUrl error",
			out:    nil,
			errRes: fmt.Errorf("aws error"),
		},
	}

	mock.SetUp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			apiMock := mock.Mock[sqsApi]()
			mock.When(apiMock.GetQueueUrl(mock.AnyContext(), mock.Any[*sqs.GetQueueUrlInput]())).ThenReturn(tt.out, tt.errRes)

			underTest := NewSqsActions(apiMock)

			url, err := underTest.GetQueueUrl("queueName")
			assert.Equal(subtest, tt.errRes, err)
			if err == nil {
				assert.Equal(subtest, *tt.out.QueueUrl, url)
			}
		})
	}
}

func TestGetMessages(t *testing.T) {
	var tests = []struct {
		name   string
		out    *sqs.ReceiveMessageOutput
		errRes error
	}{
		{
			name: "getMessages succes",
			out: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						Body:      aws.String("123"),
						MessageId: aws.String("abc123"),
					},
					{
						Body:      aws.String("asaas"),
						MessageId: aws.String("abca1fsfqs23"),
					},
				},
			},
			errRes: nil,
		},
		{
			name: "getMessages succes empty",
			out: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{},
			},
			errRes: nil,
		},
		{
			name:   "getMessages error",
			out:    nil,
			errRes: fmt.Errorf("aws error"),
		},
	}

	mock.SetUp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			apiMock := mock.Mock[sqsApi]()
			mock.When(apiMock.ReceiveMessage(mock.AnyContext(), mock.Any[*sqs.ReceiveMessageInput]())).ThenReturn(tt.out, tt.errRes)

			underTest := NewSqsActions(apiMock)

			res, err := underTest.GetMessages("url", 5, 1)
			assert.Equal(subtest, tt.errRes, err)
			if err == nil {
				assert.Equal(subtest, len(res), len(tt.out.Messages))
			}
		})
	}
}
