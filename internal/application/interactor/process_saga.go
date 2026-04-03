package interactor

import "context"

type ProcessSagaInteractor struct{}

func (i ProcessSagaInteractor) Execute(ctx context.Context) (any, error) {
	return nil, nil
}
