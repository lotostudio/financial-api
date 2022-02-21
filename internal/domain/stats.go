package domain

type Statement struct {
	// Account information
	Account Account `json:"account" binding:"required"`
	// Balance for start of period
	BalanceIn Balance `json:"balanceIn" binding:"required"`
	// Balance for end of period
	BalanceOut Balance `json:"balanceOut" binding:"required"`
	// Transactions for given period
	Transactions []Transaction `json:"transactions" binding:"required"`
} // @name Statement
