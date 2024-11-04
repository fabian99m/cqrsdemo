package handler

import (
	"cqrsdemo/adapter"
	baseSqs "cqrsdemo/handler/base"
	m "cqrsdemo/model"
	"cqrsdemo/util"
	"fmt"
	"strings"

	uc "cqrsdemo/usecase"
	"encoding/json"

	"log/slog"

	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type CmdMapper = map[string]uc.MessageHandler[m.Command]
type EvtMapper = map[string]uc.MessageHandler[m.Event]

type messageHandler struct {
	commands      CmdMapper
	events        EvtMapper
	snsOperations adapter.SnsOperations
	props         *m.EventProps
}

func NewMessageHandler(commands CmdMapper, events EvtMapper, snsActions adapter.SnsOperations, props *m.EventProps) baseSqs.SqsHandler {
	return &messageHandler{
		commands:      commands,
		events:        events,
		snsOperations: snsActions,
		props:         props,
	}
}

func (r messageHandler) ReciveMessage(message sqstypes.Message) bool {
	slog.Debug("Message received...", "atrr", &message)

	typeMessage, exit := messageValidation(message)
	if exit {
		return true
	}

	var result m.EventResult
	exit, err := r.proccesMessage(typeMessage, message, &result)
	if err != nil {
		slog.Error("error procesing command", "error", err)
		return exit
	}

	return r.publishEvent(&result) == nil
}

func (r messageHandler) proccesMessage(typeMessage m.MessageType, message sqstypes.Message, result *m.EventResult) (bool, error) {
	switch {
	case typeMessage == m.COMMAND:
		return r.processCommand(message, result)
	default:
		return r.processEvent(message, result)
	}
}

func (r messageHandler) processCommand(message sqstypes.Message, result *m.EventResult) (exit bool, _ error) {
	jsonBody := message.Body
	slog.Debug("comando recibido", "body", *jsonBody)

	command, err := util.UnmarshalTo[m.Command]([]byte(*jsonBody))
	if err != nil {
		return false, err
	}

	commandHandler, found := r.commands[command.Name]
	if !found {
		return true, fmt.Errorf("command {%s} no defined", command.Name)
	}

	*result, err = commandHandler.Process(*command)
	if err != nil {
		return false, err
	}

	return false, nil
}

func (r messageHandler) processEvent(message sqstypes.Message, result *m.EventResult) (exit bool, _ error) {
	jsonBody := message.Body
	slog.Debug("evento recibido", "body", *jsonBody)

	event, err := util.UnmarshalTo[m.Event]([]byte(*jsonBody))
	if err != nil {
		return false, err
	}

	eventHandler, found := r.events[event.Name]
	if !found {
		return true, fmt.Errorf("event %s no defined", event.Name)
	}

	*result, err = eventHandler.Process(*event)
	if err != nil {
		return false, err
	}

	return false, nil
}

func (r messageHandler) publishEvent(payload *m.EventResult) error {
	slog.Info("publishing event ", "name", payload.Name)

	jsonString, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	jsonMessage := string(jsonString)
	id, err := r.snsOperations.Publish(r.props.TopicArn, &jsonMessage, "Event")
	if err != nil {
		return err
	}

	slog.Info("event published successfully", "id", id)

	return nil
}

func messageValidation(message sqstypes.Message) (m.MessageType, bool) {
	if len(message.MessageAttributes) == 0 {
		slog.Error("error messsage without messageAttributes")
		return -1, true
	}

	typeMessageAtribute, found := message.MessageAttributes["typeMessage"]
	if !found {
		slog.Error("error messsage without typeMessage atribute")
		return -1, true
	}

	typeMessage := strings.ToUpper(*typeMessageAtribute.StringValue)
	if typeMessage != "COMMAND" && typeMessage != "EVENT" {
		slog.Error("error invalid message type", "type", typeMessage)
		return -1, true
	}

	return m.ToMessageType(typeMessage), false
}