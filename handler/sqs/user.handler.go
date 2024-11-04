package handler

import (
	base "cqrsdemo/handler/base"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"log/slog"
)

type userHandler struct{}

func NewUserHandler() base.SqsHandler {
	return &userHandler{}
}

func (r userHandler) ReciveMessage(message sqstypes.Message) bool {
	bodyJson := *message.Body
	id := *message.MessageId

	slog.Info("handling user meesage ", "body", bodyJson, "id", id)

	return bodyJson == "test"
}
