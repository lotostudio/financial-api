package domain

import "time"

type Account struct {
	ID        int64     `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Balance   float64   `json:"balance" db:"balance"`
	Currency  string    `json:"currency" db:"currency"`
	Type      string    `json:"type" db:"type"`
	OwnerId   int64     `json:"-" db:"owner_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	Term      *uint8    `json:"term,omitempty" db:"term"`
	Rate      *float32  `json:"rate,omitempty" db:"rate"`
}

type AccountToCreate struct {
	Title   string  `json:"title"`
	Balance float64 `json:"balance"`
	Type    string  `json:"type"`
	Term    uint8   `json:"term"`
	Rate    float32 `json:"rate"`
}

type AccountToUpdate struct {
	Title   *string  `json:"title"`
	Balance *float64 `json:"balance"`
	Term    *uint8   `json:"term"`
	Rate    *float32 `json:"rate"`
}
