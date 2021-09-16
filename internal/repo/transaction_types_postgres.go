package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lotostudio/financial-api/internal/domain"
)

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
