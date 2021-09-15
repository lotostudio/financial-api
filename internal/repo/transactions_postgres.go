package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lotostudio/financial-api/internal/domain"
	"time"
)

type TransactionsRepo struct {
	db *sqlx.DB
}

func newTransactionsRepo(db *sqlx.DB) *TransactionsRepo {
	return &TransactionsRepo{
		db: db,
	}
}

func (r *TransactionsRepo) List(ctx context.Context, userID int64) ([]domain.Transaction, error) {
	transactions := make([]domain.Transaction, 0)

	rows, err := r.db.QueryContext(ctx, `
	SELECT t.id, t.amount, t.type, tc.title AS category, t.created_at, 
	       cr.id, cr.title, cr.balance, cr_c.code, cr.type, cr.created_at, 
	       db.id, db.title, db.balance, db_c.code, db.type, db.created_at
	FROM transactions t
	LEFT JOIN transaction_categories tc ON t.category_id = tc.id
	LEFT JOIN accounts cr ON t.credit_id = cr.id
	LEFT JOIN currencies cr_c ON cr.currency_id = cr_c.id
	LEFT JOIN accounts db ON t.debit_id = db.id
	LEFT JOIN currencies db_c ON db.currency_id = db_c.id
	WHERE cr.owner_id = $1 OR db.owner_id = $1`, userID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		tr := domain.Transaction{}

		var creditId, debitId *int64
		var creditTitle, debitTitle, creditCurr, debitCurr *string
		var creditBalance, debitBalance *float64
		var creditType, debitType *domain.AccountType
		var creditCreatedAt, debitCreatedAt *time.Time

		if err = rows.Scan(&tr.ID, &tr.Amount, &tr.Type, &tr.Category, &tr.CreatedAt,
			&creditId, &creditTitle, &creditBalance, &creditCurr, &creditType, &creditCreatedAt,
			&debitId, &debitTitle, &debitBalance, &debitCurr, &debitType, &debitCreatedAt); err != nil {
			return nil, err
		}

		if creditId != nil {
			tr.Credit = &domain.Account{
				ID:        *creditId,
				Title:     *creditTitle,
				Balance:   *creditBalance,
				Currency:  *creditCurr,
				Type:      *creditType,
				CreatedAt: *creditCreatedAt,
			}
		}

		if debitId != nil {
			tr.Debit = &domain.Account{
				ID:        *debitId,
				Title:     *debitTitle,
				Balance:   *debitBalance,
				Currency:  *debitCurr,
				Type:      *debitType,
				CreatedAt: *debitCreatedAt,
			}
		}

		transactions = append(transactions, tr)
	}

	return transactions, nil
}

func (r *TransactionsRepo) Create(ctx context.Context, toCreate domain.TransactionToCreate, categoryId *int64,
	creditId *int64, debitId *int64) (domain.Transaction, error) {
	tx, err := r.db.Begin()

	if err != nil {
		return domain.Transaction{}, err
	}

	row := tx.QueryRowContext(ctx,
		`INSERT INTO transactions(amount, type, created_at, category_id, credit_id, debit_id) 
				VALUES ($1, $2, $3, $4, $5, $6) 
				RETURNING id, amount, type, created_at`,
		toCreate.Amount, toCreate.Type, toCreate.CreatedAt, categoryId, creditId, debitId)

	var transaction domain.Transaction

	if err = row.Scan(&transaction.ID, &transaction.Amount, &transaction.Type, &transaction.CreatedAt); err != nil {
		if err := tx.Rollback(); err != nil {
			return transaction, err
		}

		return transaction, err
	}

	if toCreate.Type == domain.Expense || toCreate.Type == domain.Transfer {
		row = tx.QueryRowContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2 RETURNING balance",
			transaction.Amount, creditId)
		var balance float64

		if err = row.Scan(&balance); err != nil {
			if err := tx.Rollback(); err != nil {
				return transaction, err
			}

			return transaction, err
		}

		if balance < 0 {
			if err = tx.Rollback(); err != nil {
				return transaction, err
			}

			return domain.Transaction{}, ErrAccountNotEnoughBalance
		}
	}

	if toCreate.Type == domain.Income || toCreate.Type == domain.Transfer {
		if _, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2",
			transaction.Amount, debitId); err != nil {
			if err := tx.Rollback(); err != nil {
				return transaction, err
			}

			return transaction, err
		}
	}

	return transaction, tx.Commit()
}
