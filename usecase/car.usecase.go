package usecase

import (
	m "cqrsdemo/model"
	"cqrsdemo/util"
	"log/slog"
)

type CarUseCase struct {
}

func NewCarUseCase() MessageHandler[m.Command] {
	return CarUseCase{}
}

func (r CarUseCase) Process(cmd m.Command) (evt m.EventResult, err error) {
	payload, err := util.UnmarshalTo[m.Dni](*cmd.Payload)
	if err != nil {
		return
	}

	slog.Info("carUseCase process", "payload", *payload)


	return m.EventResult{
		Name:  "car.out",
		Scope: "EXT",
		Payload: map[string]string{
			"message": "Test",
		},
	}, nil
}
