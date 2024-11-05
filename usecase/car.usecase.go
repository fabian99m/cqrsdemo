package usecase

import (
	"cqrsdemo/model"
	m "cqrsdemo/model"
	"cqrsdemo/util"
	"log/slog"
)

type CarUseCase struct {
	roleOperations model.RoleOperations
}

func NewCarUseCase(roleOperations model.RoleOperations) MessageHandler[m.Command] {
	return CarUseCase{
		roleOperations: roleOperations,
	}
}

func (r CarUseCase) Process(cmd m.Command) (evt m.EventResult, err error) {
	payload, err := util.UnmarshalTo[m.Dni](*cmd.Payload)
	if err != nil {
		return
	}

	slog.Info("carUseCase process", "id", cmd.IdTrazabilidad, "payload", *payload)

	id, err := r.roleOperations.SaveRole(m.Role{
		Service: "/upload",
		Role:    "testrole",
	})

	if err != nil {
		return
	}

	return event("car.out", map[string]any{
		"message": id,
	})
}

func event(name string, payload any) (m.EventResult, error) {
	return m.EventResult{
		Name:    name,
		Scope:   "EXT",
		Payload: payload,
	}, nil
}
