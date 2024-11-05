package handler

import (
	e "github.com/fabian99m/cqrsdemo/errors"
	h "github.com/fabian99m/cqrsdemo/handler/messages"

	"github.com/go-chi/render"

	"github.com/fabian99m/cqrsdemo/adapter"
	"github.com/fabian99m/cqrsdemo/model"
	"encoding/json"
	"io"
	"net/http"
)

type CommandRestHandler struct {
	snsOperations adapter.SnsOperations
	commands      h.CmdMapper
	props         *model.EventProps
}

func NewCommandRestHandler(snsActions adapter.SnsOperations, commands h.CmdMapper, props *model.EventProps) *CommandRestHandler {
	return &CommandRestHandler{
		snsOperations: snsActions,
		commands:      commands,
		props:         props,
	}
}

func (rh CommandRestHandler) Process(w http.ResponseWriter, r *http.Request) error {
	byteData, err := io.ReadAll(r.Body)
	if err != nil {
		return e.RequestError{
			StatusCode: http.StatusInternalServerError, Status: e.GenericError.Fmt(err),
		}
	}

	var cmd model.Command
	err = json.Unmarshal(byteData, &cmd)
	if err != nil {
		return e.RequestError{
			StatusCode: http.StatusInternalServerError, Status: e.GenericError.Fmt(err),
		}
	}

	if cmd.Name == "" {
		return e.RequestError{
			StatusCode: http.StatusBadRequest, Status: e.MissingCommandName,
		}
	}

	_, found := rh.commands[cmd.Name]
	if !found {
		return e.RequestError{
			StatusCode: http.StatusBadRequest, Status: e.CommandNotRegistered.Fmt(cmd.Name),
		}
	}

	jsonMessage := string(byteData)
	idPublish, err := rh.snsOperations.Publish(rh.props.TopicArn, &jsonMessage, "COMMAND")
	if err != nil {
		return e.RequestError{
			StatusCode: http.StatusInternalServerError, Status: e.GenericError.Fmt(err),
		}
	}

	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, map[string]string{
		"id":      idPublish,
		"command": cmd.Name,
	})

	return nil
}
