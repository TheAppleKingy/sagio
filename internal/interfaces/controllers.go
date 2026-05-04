package interfaces

import (
	"context"
	"encoding/json"

	"sagio/internal/application"
	"sagio/internal/application/interactors"
)

type Controller func(context.Context, []byte) error

func InitializeSaga(data []byte) error {
	var dto application.InitTransactionDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return err
	}
	return interactors.NewInitSagaInteractor().Execute(dto)
}
