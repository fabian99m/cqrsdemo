package entrypoint

import (
	"log/slog"

	adp "github.com/fabian99m/cqrsdemo/adapter"
	handler "github.com/fabian99m/cqrsdemo/handler/base"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const (
	maxMessages int32 = 10
	waitSeconds int32 = 2
)

type SqsListener struct {
	sqs adp.SqsOperations
	handler  handler.SqsHandler
	queueUrl string
}

func NewSqsListener(sqs adp.SqsOperations, handler handler.SqsHandler, queueName string) *SqsListener {
	queueUrl, err := sqs.GetQueueUrl(queueName)
	if err != nil {
		panic(err)
	}

	return &SqsListener{
		sqs:      sqs,
		handler:  handler,
		queueUrl: queueUrl,
	}
}

func (l SqsListener) Listen() {
	slog.Info("listening messages", "queueUrl", l.queueUrl)

	chnMessages := make(chan *types.Message, maxMessages)
	go func() {
		for {
			l.getMessages(chnMessages)
		}
	}()

	for message := range chnMessages {
		go l.processMessage(message)
	}
}

func (l SqsListener) getMessages(chnMessages chan<- *types.Message) {
	messages, err := l.sqs.GetMessages(l.queueUrl, maxMessages, waitSeconds)
	if err != nil {
		slog.Error("error getting messages", "error", err)
	} else if len(messages) > 0 {
		for _, message := range messages {
			chnMessages <- &message
		}
	}
}

func (l SqsListener) processMessage(message *types.Message) {
	ok := l.handler.ReciveMessage(*message)
	if ok {
		err := l.sqs.DeleteMessage(l.queueUrl, *message.ReceiptHandle)
		if err != nil {
			slog.Error("error removing message", "error", err)
		}
	}
}
