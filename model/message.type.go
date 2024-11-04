package model

type MessageType int

const (
	COMMAND MessageType = iota
	EVENT
)

func ToMessageType(typeMessage string) MessageType {
	if typeMessage == "COMMAND" {
		return COMMAND
	}

	if typeMessage == "EVENT" {
		return EVENT
	}

	return -1
}
