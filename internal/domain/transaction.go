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
	// Unique ID
	ID int64 `json:"id"  binding:"required" db:"id" example:"1"`
	// Name of category
	Title string `json:"title" binding:"required" example:"Food"`
	// Applicable type
	Type TransactionType `json:"type" binding:"required,oneof=income expense transfer" enums:"income,expense,transfer" example:"income"`
} // @name TransactionCategory

type Transaction struct {
	// Unique ID
	ID int64 `json:"id" binding:"required" db:"id" example:"1"`
	// Amount of transaction in currency of linked account
	Amount float64 `json:"amount" binding:"required,gte=0" db:"amount" example:"1230.23"`
	// Type
	Type TransactionType `json:"type" binding:"required,oneof=income expense transfer" enums:"income,expense,transfer" example:"income"`
	// Category
	Category *string `json:"category,omitempty"`
	// Date of creation
	CreatedAt time.Time `json:"createdAt" binding:"required,date" db:"created_at" format:"yyyy-MM-dd" example:"2021-09-01"`
	// Account transfer from
	Credit *Account `json:"credit,omitempty" db:"credit"`
	// Account transfer to
	Debit *Account `json:"debit,omitempty" db:"debit"`
} // @name Transaction

type TransactionsFilter struct {
	AccountId   *int64
	OwnerId     *int64
	Category    *string
	Type        *TransactionType
	CreatedFrom *time.Time
	CreatedTo   *time.Time
}

type TransactionToCreate struct {
	// Amount (in currency of accounts)
	Amount float64 `json:"amount" binding:"required,gte=0" db:"amount" example:"1230.23"`
	// Type
	Type TransactionType `json:"type" binding:"required,oneof=income expense transfer" enums:"income,expense,transfer" example:"income"`
	// Date of creation
	CreatedAt time.Time `json:"createdAt" binding:"required" db:"created_at" format:"yyyy-MM-dd" example:"2021-09-01"`
} // @name TransactionToCreate

func (t TransactionType) Validate() error {
	if t != Income && t != Expense && t != Transfer {
		return ErrInvalidTransactionType
	}

	return nil
}

type TransactionStat struct {
	// Category
	Category string `json:"category" binding:"required" db:"category" example:"food"`
	// Count
	Value int64 `json:"value" binding:"required" db:"value" example:"12"`
} // @name TransactionStat
