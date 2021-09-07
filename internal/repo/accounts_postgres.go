package repo

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lotostudio/financial-api/internal/domain"
)

type AccountsRepo struct {
	db *sqlx.DB
}

func newAccountsRepo(db *sqlx.DB) *AccountsRepo {
	return &AccountsRepo{
		db: db,
	}
}

func (r *AccountsRepo) List(ctx context.Context, userID int64) ([]domain.Account, error) {
	accounts := make([]domain.Account, 0)

	if err := r.db.SelectContext(ctx, &accounts, `
	SELECT a.id, a.title, a.balance, cur.code currency, a.type, a.created_at, 
	       coalesce(l.term, d.term) AS term, coalesce(l.rate, d.rate) AS rate, c.number
	FROM accounts a 
    LEFT JOIN loans l ON a.id = l.account_id 
    LEFT JOIN deposits d ON a.id = d.account_id
	LEFT JOIN cards c ON a.id = c.account_id
	JOIN currencies cur ON a.currency_id = cur.id
	WHERE a.owner_id = $1`, userID); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *AccountsRepo) Create(ctx context.Context, toCreate domain.AccountToCreate, userID int64, currencyID int) (domain.Account, error) {
	tx, err := r.db.Begin()

	if err != nil {
		return domain.Account{}, err
	}

	row := tx.QueryRowContext(ctx,
		`INSERT INTO accounts(title, balance, type, currency_id, owner_id) 
				VALUES ($1, $2, $3, $4, $5) 
				RETURNING id, title, balance, type, created_at`,
		toCreate.Title, toCreate.Balance, toCreate.Type, currencyID, userID)

	var account domain.Account

	if err = row.Scan(&account.ID, &account.Title, &account.Balance, &account.Type, &account.CreatedAt); err != nil {
		if err := tx.Rollback(); err != nil {
			return domain.Account{}, err
		}

		return domain.Account{}, err
	}

	if account.Type == "loan" {
		row = tx.QueryRowContext(ctx,
			`INSERT INTO loans(term, rate, account_id) VALUES ($1, $2, $3) RETURNING term, rate`,
			toCreate.Term, toCreate.Rate, account.ID)

		if err = row.Scan(&account.Term, &account.Rate); err != nil {
			if err := tx.Rollback(); err != nil {
				return domain.Account{}, err
			}

			return domain.Account{}, err
		}
	}

	if account.Type == "deposit" {
		row = tx.QueryRowContext(ctx,
			`INSERT INTO deposits(term, rate, account_id) VALUES ($1, $2, $3) RETURNING term, rate`,
			toCreate.Term, toCreate.Rate, account.ID)

		if err = row.Scan(&account.Term, &account.Rate); err != nil {
			if err := tx.Rollback(); err != nil {
				return domain.Account{}, err
			}

			return domain.Account{}, err
		}
	}

	if account.Type == "card" {
		row = tx.QueryRowContext(ctx,
			`INSERT INTO cards(number, account_id) VALUES ($1, $2) RETURNING number`,
			toCreate.Number, account.ID)

		if err = row.Scan(&account.Number); err != nil {
			if err := tx.Rollback(); err != nil {
				return domain.Account{}, err
			}

			return domain.Account{}, err
		}
	}

	return account, tx.Commit()
}

func (r *AccountsRepo) Get(ctx context.Context, id int64) (domain.Account, error) {
	var accounts domain.Account

	if err := r.db.GetContext(ctx, &accounts, `
	SELECT a.id, a.title, a.balance, cur.code currency, a.type, a.owner_id, a.created_at, 
	       coalesce(l.term, d.term) AS term, coalesce(l.rate, d.rate) AS rate, c.number
	FROM accounts a 
    LEFT JOIN loans l ON a.id = l.account_id 
    LEFT JOIN deposits d ON a.id = d.account_id 
	LEFT JOIN cards c ON a.id = c.account_id 
	JOIN currencies cur ON a.currency_id = cur.id 
	WHERE a.id = $1`, id); err != nil {
		if err == sql.ErrNoRows {
			return domain.Account{}, ErrAccountNotFound
		}

		return accounts, err
	}

	return accounts, nil
}

func (r *AccountsRepo) Update(ctx context.Context, toUpdate domain.AccountToUpdate, id int64, _type string) (domain.Account, error) {
	var account domain.Account

	tx, err := r.db.Begin()

	if err != nil {
		return account, err
	}

	row := tx.QueryRowContext(ctx, `UPDATE accounts a SET title = $1, balance = $2 WHERE a.id = $3 
	RETURNING a.id, a.title, a.balance, a.type, a.created_at`, toUpdate.Title, toUpdate.Balance, id)

	if err = row.Scan(&account.ID, &account.Title, &account.Balance, &account.Type, &account.CreatedAt); err != nil {
		if err := tx.Rollback(); err != nil {
			return account, err
		}

		return account, err
	}

	if _type == "loan" {
		row := tx.QueryRowContext(ctx, `UPDATE loans l SET term = $1, rate = $2 WHERE l.account_id = $3 
		RETURNING l.term, l.rate`, toUpdate.Term, toUpdate.Rate, id)

		if err = row.Scan(&account.Term, &account.Rate); err != nil {
			if err := tx.Rollback(); err != nil {
				return account, err
			}

			return account, err
		}
	}

	if _type == "deposit" {
		row := tx.QueryRowContext(ctx, `UPDATE deposits d SET term = $1, rate = $2 WHERE d.account_id = $3 
		RETURNING d.term, d.rate`, toUpdate.Term, toUpdate.Rate, id)

		if err = row.Scan(&account.Term, &account.Rate); err != nil {
			if err := tx.Rollback(); err != nil {
				return account, err
			}

			return account, err
		}
	}

	if _type == "card" {
		row := tx.QueryRowContext(ctx, `UPDATE cards c SET number = $1 WHERE c.account_id = $2 
		RETURNING c.number`, toUpdate.Number, id)

		if err = row.Scan(&account.Number); err != nil {
			if err := tx.Rollback(); err != nil {
				return account, err
			}

			return account, err
		}
	}

	return account, tx.Commit()
}

func (r *AccountsRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM accounts WHERE id = $1", id)

	return err
}
