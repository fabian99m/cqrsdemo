package handler

import (
	handler "github.com/fabian99m/cqrsdemo/handler/base"
	"log/slog"

	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type carHandler struct{}

func NewCarHandler() handler.SqsHandler {
	return &carHandler{}
}

func (r carHandler) ReciveMessage(message sqstypes.Message) bool {
	bodyJson := *message.Body
	id := *message.MessageId

	slog.Info("handling car meesage ", "body", bodyJson, "id", id)

	return bodyJson == "test"
}
