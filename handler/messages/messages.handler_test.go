package handler

import (
	"github.com/fabian99m/cqrsdemo/adapter"
	"github.com/fabian99m/cqrsdemo/model"
	"github.com/fabian99m/cqrsdemo/usecase"
	"fmt"

	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	mock "github.com/ovechkin-dm/mockio/mock"
	"github.com/stretchr/testify/assert"
)

type test struct {
	name         string
	message      *sqstypes.Message
	handlerError error
	publishError error
	exit         bool
}

func TestReciveMessageCommands(t *testing.T) {
	t.Parallel()
	
	var tests = []test{
		{
			name:    "noMessage atributes",
			message: &sqstypes.Message{},
			exit:    true,
		},
		{
			name: "messsage without typeMessage atribute",
			message: &sqstypes.Message{
				MessageAttributes: map[string]sqstypes.MessageAttributeValue{
					"xxxx": {
						StringValue: aws.String("text"),
					},
				},
			},
			exit: true,
		},
		{
			name: "bad message struct",
			message: &sqstypes.Message{
				Body: aws.String("{name : test}"),
				MessageAttributes: map[string]sqstypes.MessageAttributeValue{
					"typeMessage": {
						StringValue: aws.String("COMMAND"),
					},
				},
			},
			exit: false,
		},
		{
			name: "bad message struct2",
			message: &sqstypes.Message{
				Body: aws.String("{name : test}"),
				MessageAttributes: map[string]sqstypes.MessageAttributeValue{
					"typeMessage": {
						StringValue: aws.String("COMMANDx"),
					},
				},
			},
			exit: true,
		},
		{
			name: "command not found",
			message: &sqstypes.Message{
				Body: aws.String("{\"name\":\"holacommandx\",\"payload\":{\"id\":\"121213123\",\"type\":\"CE\"}}"),
				MessageAttributes: map[string]sqstypes.MessageAttributeValue{
					"typeMessage": {
						StringValue: aws.String("COMMAND"),
					},
				},
			},
			handlerError: nil,
			exit:         true,
		},
		{
			name: "command process error",
			message: &sqstypes.Message{
				Body: aws.String("{\"name\":\"holacommand\",\"payload\":{\"id\":\"121213123\",\"type\":\"CE\"}}"),
				MessageAttributes: map[string]sqstypes.MessageAttributeValue{
					"typeMessage": {
						StringValue: aws.String("COMMAND"),
					},
				},
			},
			handlerError: fmt.Errorf("error handling"),
			exit:         false,
		},
		{
			name: "command process success",
			message: &sqstypes.Message{
				Body: aws.String("{\"name\":\"holacommand\",\"payload\":{\"id\":\"121213123\",\"type\":\"CE\"}}"),
				MessageAttributes: map[string]sqstypes.MessageAttributeValue{
					"typeMessage": {
						StringValue: aws.String("COMMAND"),
					},
				},
			},
			handlerError: nil,
			exit:         true,
		},
		{
			name: "command publish event error",
			message: &sqstypes.Message{
				Body: aws.String("{\"name\":\"holacommand\",\"payload\":{\"id\":\"121213123\",\"type\":\"CE\"}}"),
				MessageAttributes: map[string]sqstypes.MessageAttributeValue{
					"typeMessage": {
						StringValue: aws.String("COMMAND"),
					},
				},
			},
			handlerError: nil,
			publishError: fmt.Errorf("error handling"),
			exit:         false,
		},
	}

	mock.SetUp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			mockSns := mock.Mock[adapter.SnsOperations]()
			mockHandler := mock.Mock[usecase.MessageHandler[model.Command]]()

			mock.When(mockSns.Publish(mock.AnyString(), mock.NotNil[*string](), mock.AnyString())).ThenReturn("", tt.publishError)
			mock.When(mockHandler.Process(mock.Any[model.Command]())).ThenReturn(model.EventResult{}, tt.handlerError)

			commands := CmdMapper{
				"holacommand": mockHandler,
			}

			underTest := NewMessageHandler(commands, EvtMapper{}, mockSns, &model.EventProps{QueueName: "test", TopicArn: "testarn"})

			exit := underTest.ReciveMessage(*tt.message)

			assert.Equal(subtest, tt.exit, exit)
		})
	}
}

func TestReciveMessageEvents(t *testing.T) {
	t.Parallel()

	var tests = []test{
		{
			name: "event not found",
			message: &sqstypes.Message{
				Body: aws.String("{\"name\":\"121313\",\"payload\":{\"id\":\"121213123\",\"type\":\"CE\"}}"),
				MessageAttributes: map[string]sqstypes.MessageAttributeValue{
					"typeMessage": {
						StringValue: aws.String("EVENT"),
					},
				},
			},
			handlerError: nil,
			exit:         true,
		},
		{
			name: "event process error",
			message: &sqstypes.Message{
				Body: aws.String("{\"name\":\"event1\",\"payload\":{\"id\":\"121213123\",\"type\":\"CE\"}}"),
				MessageAttributes: map[string]sqstypes.MessageAttributeValue{
					"typeMessage": {
						StringValue: aws.String("EVENT"),
					},
				},
			},
			handlerError: fmt.Errorf("error handling"),
			exit:         false,
		},
		{
			name: "event success",
			message: &sqstypes.Message{
				Body: aws.String("{\"name\":\"event1\",\"payload\":{\"id\":\"121213123\",\"type\":\"CE\"}}"),
				MessageAttributes: map[string]sqstypes.MessageAttributeValue{
					"typeMessage": {
						StringValue: aws.String("EVENT"),
					},
				},
			},
			handlerError: nil,
			exit:         true,
		},
		{
			name:  "event publish event error",
			message: &sqstypes.Message{
				Body: aws.String("{\"name\":\"event1\",\"payload\":{\"id\":\"121213123\",\"type\":\"CE\"}}"),
				MessageAttributes: map[string]sqstypes.MessageAttributeValue{
					"typeMessage": {
						StringValue: aws.String("EVENT"),
					},
				},
			},
			handlerError: nil,
			publishError: fmt.Errorf("error handling"),
			exit:         false,
		},
	}

	mock.SetUp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			mockSns := mock.Mock[adapter.SnsOperations]()
			mockHandler := mock.Mock[usecase.MessageHandler[model.Event]]()

			mock.When(mockSns.Publish(mock.AnyString(), mock.NotNil[*string](), mock.AnyString())).ThenReturn("", tt.publishError)
			mock.When(mockHandler.Process(mock.Any[model.Event]())).ThenReturn(model.EventResult{}, tt.handlerError)

			events := EvtMapper{
				"event1": mockHandler,
			}

			underTest := NewMessageHandler(CmdMapper{}, events, mockSns, &model.EventProps{QueueName: "test", TopicArn: "testarn"})

			exit := underTest.ReciveMessage(*tt.message)

			assert.Equal(subtest, tt.exit, exit)
		})
	}
}
