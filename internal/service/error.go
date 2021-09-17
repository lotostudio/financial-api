package service

import "errors"

var (
	ErrInvalidLoanData    = errors.New("account with type 'loan' must have valid term and rate")
	ErrInvalidDepositData = errors.New("account with type 'deposit' must have valid term and rate")
	ErrInvalidCardData    = errors.New("account with type 'card' must have valid number")
	ErrAccountForbidden   = errors.New("account forbidden to access")

	ErrRefreshTokenExpired = errors.New("refresh token expired")
)
