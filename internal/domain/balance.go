package domain

import "time"

type Balance struct {
	AccountID int64     `json:"-" db:"account_id"`
	Date      time.Time `json:"date" db:"date"`
	Value     int64     `json:"value" db:"value"`
}
