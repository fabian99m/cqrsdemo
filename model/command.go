package model

import "encoding/json"

type Command struct {
	BaseHandler
}

type Event struct {
	BaseHandler
}

type BaseHandler struct {
	Name    string           `json:"name"`
	Payload *json.RawMessage `json:"payload"`
}

type EventResult struct {
	Name    string `json:"name"`
	Payload any    `json:"payload"`
	Scope   string `json:"scope"`
}

type Dni struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}
