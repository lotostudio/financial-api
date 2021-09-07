package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lotostudio/financial-api/internal/domain"
)

type AccountTypesRepo struct {
	db *sqlx.DB
}

func newAccountTypesRepo(db *sqlx.DB) *AccountTypesRepo {
	return &AccountTypesRepo{
		db: db,
	}
}

func (r *AccountTypesRepo) List(ctx context.Context) ([]domain.AccountType, error) {
	types := make([]domain.AccountType, 0)

	if err := r.db.SelectContext(ctx, &types, "SELECT unnest(enum_range(NULL::account_type))"); err != nil {
		return nil, err
	}

	return types, nil
}
