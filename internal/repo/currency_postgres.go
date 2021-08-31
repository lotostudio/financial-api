package repo

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lotostudio/financial-api/internal/domain"
)

type CurrenciesRepo struct {
	db *sqlx.DB
}

func newCurrenciesRepo(db *sqlx.DB) *CurrenciesRepo {
	return &CurrenciesRepo{
		db: db,
	}
}

func (r *CurrenciesRepo) List(ctx context.Context) ([]domain.Currency, error) {
	currencies := make([]domain.Currency, 0)

	if err := r.db.SelectContext(ctx, &currencies, "SELECT c.id, c.code FROM currencies c"); err != nil {
		return nil, err
	}

	return currencies, nil
}

func (r *CurrenciesRepo) Get(ctx context.Context, id int) (domain.Currency, error) {
	var currency domain.Currency

	if err := r.db.GetContext(ctx, &currency, `SELECT c.id, c.code FROM currencies c WHERE c.id = $1`, id); err != nil {
		if err == sql.ErrNoRows {
			return domain.Currency{}, ErrCurrencyNotFound
		}

		return currency, err
	}

	return currency, nil
}
