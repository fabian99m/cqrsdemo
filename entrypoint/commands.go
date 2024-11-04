package entrypoint

import (
	"cqrsdemo/adapter"
	base "cqrsdemo/handler/base"
	"cqrsdemo/handler/messages"
	"cqrsdemo/model"
	"cqrsdemo/usecase"
)

func NewMessagesHandler(snsActions adapter.SnsOperations, props *model.EventProps) base.SqsHandler {
	return handler.NewMessageHandler(CommandsHandler(), EventsHandler(), snsActions, props)
}

func CommandsHandler() handler.CmdMapper {
	return handler.CmdMapper{
		"holacommand": usecase.NewCarUseCase(),
	}
}

func EventsHandler() handler.EvtMapper {
	return handler.EvtMapper{
		"car.out": usecase.NewUserUseCase(),
	}
}
