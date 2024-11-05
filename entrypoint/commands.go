package entrypoint

import (
	"github.com/fabian99m/cqrsdemo/adapter"
	bdCfg "github.com/fabian99m/cqrsdemo/config/db"
	config "github.com/fabian99m/cqrsdemo/config/props"
	"github.com/fabian99m/cqrsdemo/handler/messages"
	"github.com/fabian99m/cqrsdemo/repository"
	"github.com/fabian99m/cqrsdemo/usecase"
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
