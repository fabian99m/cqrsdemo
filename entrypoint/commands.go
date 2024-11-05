package entrypoint

import (
	"cqrsdemo/adapter"
	bdCfg "cqrsdemo/config/db"
	config "cqrsdemo/config/props"
	"cqrsdemo/handler/messages"
	"cqrsdemo/repository"
	"cqrsdemo/usecase"
)

func NewMessagesHandler(snsActions adapter.SnsOperations, props *config.AppConfig) *handler.MessageHandler {
	bdCon := bdCfg.NewDbConnection(props)
	roleRepostory := repository.NewRoleRepository(bdCon)

	cmds := handler.CmdMapper{
		"holacommand": usecase.NewCarUseCase(roleRepostory),
	}

	evts := handler.EvtMapper{
		"car.out": usecase.NewUserUseCase(),
	}
	
	return handler.NewMessageHandler(cmds, evts, snsActions, props.EventsProps.Get())
}
