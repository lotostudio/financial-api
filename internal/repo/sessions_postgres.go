package repo

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lotostudio/financial-api/internal/domain"
)

type SessionsRepo struct {
	db *sqlx.DB
}

func newSessionsRepo(db *sqlx.DB) *SessionsRepo {
	return &SessionsRepo{db: db}
}

func (r *SessionsRepo) Create(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO sessions (user_id) VALUES ($1)", userID)

	return err
}

func (r *SessionsRepo) GetByToken(ctx context.Context, token string) (domain.Session, error) {
	var item domain.Session

	if err := r.db.GetContext(ctx, &item, `SELECT s.* FROM sessions s WHERE s.refresh_token = $1`, token); err != nil {

		if err == sql.ErrNoRows {
			return item, ErrSessionNotFound
		}

		return item, err
	}

	return item, nil
}

func (r *SessionsRepo) Update(ctx context.Context, toUpdate domain.SessionToUpdate, id int64) (domain.Session, error) {
	var session domain.Session

	err := r.db.GetContext(ctx, &session,
		"UPDATE sessions s SET refresh_token = $1, expires_at = $2 WHERE s.id =$3 RETURNING *",
		toUpdate.RefreshToken, toUpdate.ExpiresAt, id)

	if err != nil {
		return session, err
	}

	return session, nil
}
