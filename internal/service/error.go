package service

import "errors"

var (
	ErrInvalidLoanData    = errors.New("account with type 'loan' must have valid term and rate")
	ErrInvalidDepositData = errors.New("account with type 'deposit' must have valid term and rate")
	ErrInvalidCardData    = errors.New("account with type 'card' must have valid number")
	ErrAccountForbidden   = errors.New("account forbidden to access")

	ErrRefreshTokenExpired = errors.New("refresh token expired")

	ErrAccountsHaveDifferenceCurrencies = errors.New("accounts have different currencies")
	ErrCreditAccountForbidden           = errors.New("sender account forbidden to access")
	ErrDebitAccountForbidden            = errors.New("receiver account forbidden to access")
	ErrNoAccountSelected                = errors.New("no account selected")

	ErrTransactionForbidden                = errors.New("transaction forbidden to access")
	ErrTransactionAndCategoryTypesMismatch = errors.New("type of transaction and category does not match")
)
