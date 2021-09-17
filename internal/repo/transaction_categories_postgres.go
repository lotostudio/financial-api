package repo

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lotostudio/financial-api/internal/domain"
)

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
