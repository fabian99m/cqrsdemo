package usecase

import (
	m "github.com/fabian99m/cqrsdemo/model"
	"github.com/fabian99m/cqrsdemo/util"
	"log/slog"
)

type UserUseCase struct {
}

func NewUserUseCase() MessageHandler[m.Event] {
	return &UserUseCase{}
}

func (r UserUseCase) Process(evt m.Event) (result m.EventResult, err error) {
	payload, err := util.UnmarshalTo[map[string]string](*evt.Payload)
	if err != nil {
		return
	}

	slog.Info("userUseCase event process", "payload", *payload)

	return m.EventResult{
		Name:  "car.test",
		Scope: "UI",
		Payload: map[string]string{
			"message": "Test",
		},
	}, nil
}
