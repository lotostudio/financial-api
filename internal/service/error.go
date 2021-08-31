package service

import "errors"

var (
	ErrInvalidLoanData    = errors.New("account with type 'loan' must have valid term and rate")
	ErrInvalidDepositData = errors.New("account with type 'deposit' must have valid term and rate")
	ErrAccountForbidden   = errors.New("account forbidden to access")
)
