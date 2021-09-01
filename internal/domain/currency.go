package domain

type Currency struct {
	ID   int    `json:"id" binding:"required" db:"id" example:"1"`
	Code string `json:"code" binding:"required,max=10" maxLength:"10" db:"code" example:"KZT"`
} // @name Currency
