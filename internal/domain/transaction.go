package domain

import "time"

// Transaction types
const (
	Income   = TransactionType("income")
	Expense  = TransactionType("expense")
	Transfer = TransactionType("transfer")
)

type TransactionType string // @name TransactionType

type TransactionCategory struct {
	ID    int64           `json:"id"  binding:"required" db:"id" example:"1"`
	Title string          `json:"title" binding:"required" example:"Food"`
	Type  TransactionType `json:"type" binding:"required,oneof=income expense transfer" enums:"income,expense,transfer" example:"income"`
}

type Transaction struct {
	ID        int64           `json:"id" binding:"required" db:"id" example:"1"`
	Amount    float64         `json:"amount" binding:"required,gte=0" db:"amount" example:"1230.23"`
	Type      TransactionType `json:"type" binding:"required,oneof=income expense transfer" enums:"income,expense,transfer" example:"income"`
	Category  *string         `json:"category,omitempty"`
	CreatedAt time.Time       `json:"createdAt" binding:"required,date" db:"created_at" format:"yyyy-MM-dd" example:"2021-09-01"`
	Credit    *Account        `json:"credit,omitempty" db:"credit"`
	Debit     *Account        `json:"debit,omitempty" db:"debit"`
}

type TransactionsFilter struct {
	AccountId   *int64
	OwnerId     *int64
	Category    *string
	CreatedFrom *time.Time
	CreatedTo   *time.Time
}

type TransactionToCreate struct {
	Amount    float64         `json:"amount" binding:"required,gte=0" db:"amount" example:"1230.23"`
	Type      TransactionType `json:"type" binding:"required,oneof=income expense transfer" enums:"income,expense,transfer" example:"income"`
	CreatedAt time.Time       `json:"createdAt" binding:"required" db:"created_at" format:"yyyy-MM-dd" example:"2021-09-01"`
}

func (t TransactionType) Validate() error {
	if t != Income && t != Expense && t != Transfer {
		return ErrInvalidTransactionType
	}

	return nil
}
