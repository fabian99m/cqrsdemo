package usecase

import (
    m "cqrsdemo/model"
)

type MessageHandler[T Message] interface {
    Process(T) (m.EventResult, error)
}

type Message interface {
    m.Command | m.Event
}
