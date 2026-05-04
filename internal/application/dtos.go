package application

type InitTransactionDTO struct {
	Name  string        `json:"name" binding:"required"`
	Steps []SagaStepDTO `json:"steps" binding:"required"`
	ID    string        `json:"id"`
}

type ProcessTransactionDTO struct {
	Name string `json:"name" binding:"required"`
}

type SagaStepDTO struct {
	Name    string `json:"name" binding:"required"`
	StepNum int    `json:"step" binding:"required"`
}
