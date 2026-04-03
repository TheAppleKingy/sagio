package dtos

type InitTransactionDTO struct {
	Name string `json:"name" binding:"required"`
}
