package domain

import "time"

type AccountType string // @name AccountType

// Account types
const (
	Cash    = AccountType("cash")
	Card    = AccountType("card")
	Loan    = AccountType("loan")
	Deposit = AccountType("deposit")
)

type Account struct {
	// Unique id
	ID int64 `json:"id" binding:"required" db:"id" example:"1"`
	// Main purpose
	Title string `json:"title" binding:"required" db:"title" example:"Main savings"`
	// Current amount of money
	Balance float64 `json:"balance" binding:"required,gte=0" db:"balance" example:"123002.12"`
	// Currency
	Currency string `json:"currency" binding:"required" db:"currency" example:"KZT"`
	// Type (different types have distinct data)
	Type    AccountType `json:"type" binding:"required,oneof=card cash loan deposit" db:"type" enums:"card,cash,loan,deposit" example:"deposit"`
	OwnerId int64       `json:"-" db:"owner_id" swaggerignore:"true"`
	// Time of creation
	CreatedAt time.Time `json:"createdAt" binding:"required,datetime" db:"created_at" format:"yyyy-MM-ddThh:mm:ss.ZZZ" example:"2021-09-01T18:03:24.499198Z"`
	// Applicable for cards
	// * For cards - last 4 digit of card number
	Number *string `json:"number,omitempty" binding:"omitempty,numeric,len=4" db:"number" example:"0327"`
	// Applicable for loans and deposits
	// * For loans - loan term
	// * For deposits - term of deposit
	Term *uint8 `json:"term,omitempty" binding:"omitempty,gt=0" db:"term" example:"12"`
	// Applicable for loans and deposits
	// * For loans - loan interest
	// * For deposits - deposit percentage
	Rate *float32 `json:"rate,omitempty" binding:"omitempty,gt=0" db:"rate" example:"10.8"`
} // @name Account

type AccountToCreate struct {
	// Main purpose
	Title string `json:"title" binding:"required" example:"Main savings"`
	// Current amount of money
	Balance float64 `json:"balance" binding:"required,gte=0" example:"123002.12"`
	// Type (different types have distinct data)
	Type AccountType `json:"type" binding:"required,oneof=card cash loan deposit" enums:"card,cash,loan,deposit" example:"deposit"`
	// Applicable for cards
	// * For cards - last 4 digit of card number
	Number *string `json:"number" binding:"omitempty,numeric,len=4" example:"0327"`
	// In months. Applicable for loans and deposits
	// * For loans - loan term
	// * For deposits - term of deposit
	Term *uint8 `json:"term" binding:"omitempty,gt=0" example:"12"`
	// Applicable for loans and deposits
	// * For loans - loan interest
	// * For deposits - deposit percentage
	Rate *float32 `json:"rate" binding:"omitempty,gt=0" example:"10.8"`
} // @name AccountToCreate

type AccountToUpdate struct {
	// Main purpose
	Title *string `json:"title" example:"Secondary savings"`
	// Current amount of money
	Balance *float64 `json:"balance" binding:"omitempty,gte=0" example:"123002.12"`
	// Applicable for cards
	// * For cards - last 4 digit of card number
	Number *string `json:"number" binding:"omitempty,numeric,len=4" example:"0327"`
	// Applicable for loans and deposits
	// * For loans - loan term
	// * For deposits - term of deposit
	Term *uint8 `json:"term" binding:"omitempty,gt=0" example:"12"`
	// Applicable for loans and deposits
	// * For loans - loan interest
	// * For deposits - deposit percentage
	Rate *float32 `json:"rate" binding:"omitempty,gt=0" example:"10.8"`
} // @name AccountToUpdate

type GroupedAccounts map[AccountType][]Account // @name GroupedAccounts

type Balance struct {
	AccountID int64 `json:"-" db:"account_id" swaggerignore:"true"`
	// Date of balance
	Date time.Time `json:"date" binding:"required" db:"date" format:"yyyy-MM-ddThh:mm:ss.ZZZ" example:"2022-01-15T00:00:00Z"`
	// Amount of balance
	Value int64 `json:"value" binding:"required" db:"value" example:"123002.12"`
} // @name Balance

type Currency struct {
	ID   int    `json:"id" binding:"required" db:"id" example:"1"`
	Code string `json:"code" binding:"required,max=10" maxLength:"10" db:"code" example:"KZT"`
} // @name Currency
