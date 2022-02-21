package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lotostudio/financial-api/internal/domain"
	"strings"
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

func (r *TransactionsRepo) List(ctx context.Context, filter domain.TransactionsFilter) ([]domain.Transaction, error) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	var argId = 1

	setValues = append(setValues, "1 = 1")

	if filter.OwnerId != nil {
		setValues = append(setValues, fmt.Sprintf("(cr.owner_id = $%d OR db.owner_id = $%d)", argId, argId))
		args = append(args, *filter.OwnerId)
		argId++
	}

	if filter.AccountId != nil {
		setValues = append(setValues, fmt.Sprintf("(cr.id=$%d OR db.id=$%d)", argId, argId))
		args = append(args, *filter.AccountId)
		argId++
	}

	if filter.Category != nil {
		setValues = append(setValues, fmt.Sprintf("tc.title=$%d", argId))
		args = append(args, *filter.Category)
		argId++
	}

	if filter.Type != nil {
		setValues = append(setValues, fmt.Sprintf("t.type=$%d", argId))
		args = append(args, *filter.Type)
		argId++
	}

	if filter.CreatedFrom != nil && filter.CreatedTo != nil {
		setValues = append(setValues, fmt.Sprintf("t.created_at BETWEEN $%d AND $%d", argId, argId+1))
		args = append(args, *filter.CreatedFrom, *filter.CreatedTo)

		// last argId increasing must be marked as nolint
		argId = argId + 2 //nolint
	}

	// Create WHERE statement variables with separated by ANDs
	setQuery := strings.Join(setValues, " AND ")
	query := fmt.Sprintf(`
	SELECT t.id, t.amount, t.type, tc.title AS category, t.created_at, 
	       cr.id, cr.title, cr.balance, cr_c.code, cr.type, cr.created_at, 
	       db.id, db.title, db.balance, db_c.code, db.type, db.created_at
	FROM transactions t
	LEFT JOIN transaction_categories tc ON t.category_id = tc.id
	LEFT JOIN accounts cr ON t.credit_id = cr.id
	LEFT JOIN currencies cr_c ON cr.currency_id = cr_c.id
	LEFT JOIN accounts db ON t.debit_id = db.id
	LEFT JOIN currencies db_c ON db.currency_id = db_c.id
	WHERE %s`, setQuery)

	rows, err := r.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	transactions := make([]domain.Transaction, 0)

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

func (r *TransactionsRepo) Stats(ctx context.Context, filter domain.TransactionsFilter) ([]domain.TransactionStat, error) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	var argId = 1

	setValues = append(setValues, "1 = 1")

	if filter.OwnerId != nil {
		setValues = append(setValues, fmt.Sprintf("(cr.owner_id = $%d OR db.owner_id = $%d)", argId, argId))
		args = append(args, *filter.OwnerId)
		argId++
	}

	if filter.AccountId != nil {
		setValues = append(setValues, fmt.Sprintf("(cr.id=$%d OR db.id=$%d)", argId, argId))
		args = append(args, *filter.AccountId)
		argId++
	}

	if filter.Category != nil {
		setValues = append(setValues, fmt.Sprintf("tc.title=$%d", argId))
		args = append(args, *filter.Category)
		argId++
	}

	if filter.Type != nil {
		setValues = append(setValues, fmt.Sprintf("t.type=$%d", argId))
		args = append(args, *filter.Type)
		argId++
	}

	if filter.CreatedFrom != nil && filter.CreatedTo != nil {
		setValues = append(setValues, fmt.Sprintf("t.created_at BETWEEN $%d AND $%d", argId, argId+1))
		args = append(args, *filter.CreatedFrom, *filter.CreatedTo)

		// last argId increasing must be marked as nolint
		argId = argId + 2 //nolint
	}

	// Create WHERE statement variables with separated by ANDs
	setQuery := strings.Join(setValues, " AND ")
	query := fmt.Sprintf(`
	SELECT coalesce(tc.title, 'transfer') AS category, sum(t.amount) AS value
	FROM transactions t
	LEFT JOIN transaction_categories tc ON t.category_id = tc.id
	LEFT JOIN accounts cr ON t.credit_id = cr.id
	LEFT JOIN currencies cr_c ON cr.currency_id = cr_c.id
	LEFT JOIN accounts db ON t.debit_id = db.id
	LEFT JOIN currencies db_c ON db.currency_id = db_c.id
	WHERE %s
	GROUP BY tc.title`, setQuery)

	fmt.Println(query)

	rows, err := r.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	stats := make([]domain.TransactionStat, 0)

	for rows.Next() {
		st := domain.TransactionStat{}

		if err = rows.Scan(&st.Category, &st.Value); err != nil {
			return nil, err
		}

		stats = append(stats, st)
	}

	return stats, nil
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

			return transaction, ErrAccountNotEnoughBalance
		}

		if err := updateBalance(ctx, tx, *creditId, balance); err != nil {
			return transaction, err
		}
	}

	if toCreate.Type == domain.Income || toCreate.Type == domain.Transfer {
		row = tx.QueryRowContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2 RETURNING balance",
			transaction.Amount, debitId)
		var balance float64

		if err = row.Scan(&balance); err != nil {
			if err := tx.Rollback(); err != nil {
				return transaction, err
			}

			return transaction, err
		}

		if err := updateBalance(ctx, tx, *debitId, balance); err != nil {
			return transaction, err
		}
	}

	return transaction, tx.Commit()
}

func (r *TransactionsRepo) GetOwner(ctx context.Context, id int64) (int64, error) {
	rows, err := r.db.QueryContext(ctx, `
	SELECT cr.owner_id, db.owner_id 
	FROM transactions t
	LEFT JOIN accounts cr ON t.credit_id = cr.id
	LEFT JOIN accounts db ON t.debit_id = db.id
	WHERE t.id = $1`, id)

	if err != nil {
		return 0, err
	}

	if !rows.Next() {
		return 0, ErrTransactionNotFound
	}

	var creditOwner, debitOwner *int64

	if err = rows.Scan(&creditOwner, &debitOwner); err != nil {
		return 0, err
	}

	if creditOwner != nil {
		return *creditOwner, nil
	}

	if debitOwner != nil {
		return *debitOwner, nil
	}

	return 0, ErrTransactionOwnerNotFound
}

func (r *TransactionsRepo) Delete(ctx context.Context, id int64) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	var creditId, debitId *int64
	var amount float64

	row := tx.QueryRowContext(ctx,
		"DELETE FROM transactions t WHERE t.id = $1 RETURNING t.credit_id, t.debit_id, t.amount", id)

	if err = row.Scan(&creditId, &debitId, &amount); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	if creditId != nil {
		row = tx.QueryRowContext(ctx,
			"UPDATE accounts SET balance = balance + $1 WHERE id = $2 RETURNING balance", amount, creditId)
		var balance float64

		if err = row.Scan(&balance); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}

			return err
		}

		if err := updateBalance(ctx, tx, *creditId, balance); err != nil {
			return err
		}
	}

	if debitId != nil {
		row = tx.QueryRowContext(ctx,
			"UPDATE accounts SET balance = balance - $1 WHERE id = $2 RETURNING balance", amount, debitId)
		var balance float64

		if err = row.Scan(&balance); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}

			return err
		}

		if balance < 0 {
			if err = tx.Rollback(); err != nil {
				return err
			}

			return ErrAccountNotEnoughBalance
		}

		if err := updateBalance(ctx, tx, *debitId, balance); err != nil {
			return err
		}
	}

	return tx.Commit()
}

type TransactionCategoryRepo struct {
	db *sqlx.DB
}

func newTransactionCategoriesRepo(db *sqlx.DB) *TransactionCategoryRepo {
	return &TransactionCategoryRepo{
		db: db,
	}
}

func (r *TransactionCategoryRepo) List(ctx context.Context) ([]domain.TransactionCategory, error) {
	categories := make([]domain.TransactionCategory, 0)

	if err := r.db.SelectContext(ctx, &categories, "SELECT c.id, c.title, c.type FROM transaction_categories c"); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *TransactionCategoryRepo) ListByType(ctx context.Context, _type domain.TransactionType) ([]domain.TransactionCategory, error) {
	categories := make([]domain.TransactionCategory, 0)

	if err := r.db.SelectContext(ctx, &categories,
		"SELECT c.id, c.title, c.type FROM transaction_categories c WHERE c.type = $1", _type); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *TransactionCategoryRepo) Get(ctx context.Context, id int64) (domain.TransactionCategory, error) {
	var category domain.TransactionCategory

	if err := r.db.GetContext(ctx, &category,
		"SELECT c.id, c.title, c.type FROM transaction_categories c WHERE c.id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return category, ErrTransactionCategoryNotFound
		}

		return category, err
	}

	return category, nil
}

type TransactionTypesRepo struct {
	db *sqlx.DB
}

func newTransactionTypesRepo(db *sqlx.DB) *TransactionTypesRepo {
	return &TransactionTypesRepo{
		db: db,
	}
}

func (r *TransactionTypesRepo) List(ctx context.Context) ([]domain.TransactionType, error) {
	types := make([]domain.TransactionType, 0)

	if err := r.db.SelectContext(ctx, &types, "SELECT unnest(enum_range(NULL::transaction_type))"); err != nil {
		return nil, err
	}

	return types, nil
}
