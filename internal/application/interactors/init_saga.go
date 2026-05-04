package interactors

import "sagio/internal/application"

// InitSagaInteractor is an interactor that initializes new saga transaction
type InitSagaInteractor struct{}

func NewInitSagaInteractor() *InitSagaInteractor {
	return &InitSagaInteractor{}
}

func (i InitSagaInteractor) Execute(_ application.InitTransactionDTO) error {
	return nil
}
