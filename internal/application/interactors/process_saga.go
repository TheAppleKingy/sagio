package interactors

import "context"

// ProcessSagaInteractor is an interactor that execute another one part of entire saga transaction
type ProcessSagaInteractor struct{}

func (i ProcessSagaInteractor) Execute(_ context.Context) (any, error) {
	return nil, nil
}
