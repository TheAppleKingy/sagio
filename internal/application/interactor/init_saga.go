package interactor

import "context"

type InitSagaInteractor struct{}

func (i InitSagaInteractor) Execute(ctx context.Context) (any, error) {
	return nil, nil
}
