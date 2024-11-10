package entrypoint

import (
	"fmt"
	"testing"

	adr "github.com/fabian99m/cqrsdemo/adapter"
	handler "github.com/fabian99m/cqrsdemo/handler/base"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
	mock "github.com/ovechkin-dm/mockio/mock"
	"github.com/stretchr/testify/assert"
)

func TestListenQueueNotFound(t *testing.T) {
	sqsMock := mock.Mock[adr.SqsOperations]()
	mock.When(sqsMock.GetQueueUrl(mock.AnyString())).ThenReturn("", fmt.Errorf("aws error"))

	assert.Panicsf(t, func() { NewSqsListener(sqsMock, nil, "que") }, "must panic if not queue url")
}

func TestListenProcessMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		handlerSuccess bool
		deleteError    error
	}{
		{
			name:           "handler succes",
			handlerSuccess: true,
			deleteError:    nil,
		},
		{
			name:           "handler fail",
			handlerSuccess: false,
			deleteError:    nil,
		},
		{
			name:           "delete fail",
			handlerSuccess: true,
			deleteError:    fmt.Errorf("aws error"),
		},
	}

	mock.SetUp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			sqsMock := mock.Mock[adr.SqsOperations]()
			mock.When(sqsMock.DeleteMessage(mock.AnyString(), mock.AnyString())).ThenReturn(tt.deleteError)
			mock.When(sqsMock.GetQueueUrl(mock.AnyString())).ThenReturn("url", nil)

			handlerMock := mock.Mock[handler.SqsHandler]()
			mock.WhenSingle(handlerMock.ReciveMessage(mock.Any[types.Message]())).ThenReturn(tt.handlerSuccess)

			underTest := NewSqsListener(sqsMock, handlerMock, "queue")

			assert.NotPanics(subtest, func() {
				underTest.processMessage(&types.Message{
					ReceiptHandle: aws.String("abc"),
				})
			})
		})
	}
}

func TestListenReceiveMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		serviceOut   []types.Message
		serviceError error
	}{
		{
			name: "success",
			serviceOut: []types.Message{
				{
					Body:      aws.String("test"),
					MessageId: aws.String(uuid.New().String()),
				},
				{
					Body:      aws.String("test2"),
					MessageId: aws.String(uuid.New().String()),
				},
			},
			serviceError: nil,
		},
		{
			name:         "error getting messages",
			serviceOut:   nil,
			serviceError: fmt.Errorf("aws error"),
		},
	}

	mock.SetUp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			sqsMock := mock.Mock[adr.SqsOperations]()
			mock.WhenDouble(sqsMock.GetMessages(mock.AnyString(), mock.Any[int32](), mock.Any[int32]())).ThenReturn(tt.serviceOut, tt.serviceError)
			mock.When(sqsMock.GetQueueUrl(mock.AnyString())).ThenReturn("url", nil)

			underTest := NewSqsListener(sqsMock, nil, "quewue")

			chnMessages := make(chan<- *types.Message, 10)
			underTest.getMessages(chnMessages)

			assert.Equal(subtest, len(tt.serviceOut), len(chnMessages))
		})
	}
}
