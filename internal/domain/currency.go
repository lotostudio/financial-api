package domain

type Currency struct {
	ID   int64  `json:"id" db:"id"`
	Code string `json:"code" db:"code"`
}
