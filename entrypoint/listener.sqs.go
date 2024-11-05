package entrypoint

import (
	"github.com/fabian99m/cqrsdemo/adapter"
	handler "github.com/fabian99m/cqrsdemo/handler/base"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

var (
	maxMessages int32 = 10
	waitSeconds int32 = 2
)

type SqsListener struct {
	sqs adapter.SqsOperations
}

func NewSqsListener(sqs adapter.SqsOperations) *SqsListener {
	return &SqsListener{sqs: sqs}
}

func (r SqsListener) Listen(handler handler.SqsHandler, queueName string) {
	slog.Info("listening messages", "queue", queueName)

	queueUrl, err := r.sqs.GetQueueUrl(queueName)
	if err != nil {
		panic(err)
	}

	chnMessages := make(chan *types.Message, maxMessages)
	go func() {
		for {
			r.getMessages(chnMessages, queueUrl)
		}
	}()

	for message := range chnMessages {
		go r.processMessages(message, handler, queueUrl)
	}
}

func (r SqsListener) getMessages(chnMessages chan<- *types.Message, queueUrl string) {
	messages, err := r.sqs.GetMessages(queueUrl, maxMessages, waitSeconds)
	if err != nil {
		slog.Error("error getting messages", "error", err)
	} else if len(messages) > 0 {
		for _, message := range messages {
			chnMessages <- &message
		}
	}
}

func (r SqsListener) processMessages(message *types.Message, handler handler.SqsHandler, queueUrl string) {
	ok := handler.ReciveMessage(*message)
	if ok {
		err := r.sqs.DeleteMessage(queueUrl, *message.ReceiptHandle)
		if err != nil {
			slog.Error("error removing message", "error", err)
		}
	}
}
