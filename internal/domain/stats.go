package domain

type Statement struct {
	Account      Account       `json:"account"`
	BalanceIn    Balance       `json:"balanceIn"`
	BalanceOut   Balance       `json:"balanceOut"`
	Transactions []Transaction `json:"transactions"`
}
