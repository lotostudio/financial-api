package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lotostudio/financial-api/internal/domain"
	"time"
)

type BalancesRepo struct {
	db *sqlx.DB
}

func newBalancesRepo(db *sqlx.DB) *BalancesRepo {
	return &BalancesRepo{
		db: db,
	}
}

func (r *BalancesRepo) Get(ctx context.Context, accountID int64, date time.Time) (domain.Balance, error) {
	var b domain.Balance

	if err := r.db.GetContext(ctx, &b, `
	SELECT account_id, date, value
	FROM balances
	WHERE account_id = $1 AND date < $2
	ORDER BY date DESC LIMIT 1`, accountID, date); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Balance{}, ErrBalanceNotFound
		}

		return b, err
	}

	return b, nil
}

// updateBalance actualize account balance into balances table with recent value
func updateBalance(ctx context.Context, tx *sql.Tx, id int64, balance float64) error {
	// update new entry in balances table for today
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO balances(account_id, date, value) VALUES ($1, $2, $3) 
				ON CONFLICT (account_id, date) DO UPDATE SET value = excluded.value`,
		id, time.Now(), balance); err != nil {

		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	return nil
}
