package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lotostudio/financial-api/internal/domain"
	log "github.com/sirupsen/logrus"
	"strings"
)

type UsersRepo struct {
	db *sqlx.DB
}

func newUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) List(ctx context.Context) ([]domain.User, error) {
	users := make([]domain.User, 0)

	if err := r.db.Select(&users, "SELECT * FROM users"); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UsersRepo) Create(ctx context.Context, user domain.User) (int64, error) {
	tx, err := r.db.Begin()

	if err != nil {
		return 0, err
	}

	var userId int64

	row := tx.QueryRow("INSERT INTO users (email, first_name, last_name, password) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Email, user.FirstName, user.LastName, user.Password)

	if err = row.Scan(&userId); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error(err)
		}

		// If user already exists
		if err, ok := err.(*pq.Error); ok || err.Code.Name() == "unique_violation" {
			return 0, ErrUserAlreadyExists
		}

		return 0, err
	}

	return userId, tx.Commit()
}

func (r *UsersRepo) Get(ctx context.Context, id int64) (domain.User, error) {
	var item domain.User

	if err := r.db.Get(&item, `SELECT * FROM users WHERE users.id = $1`, id); err != nil {
		return item, err
	}

	return item, nil
}

func (r *UsersRepo) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var item domain.User

	if err := r.db.Get(&item, `SELECT * FROM users WHERE users.email = $1 AND users.password = $2`,
		email, password); err != nil {

		if err == sql.ErrNoRows {
			return item, ErrUserNotFound
		}

		return item, err
	}

	return item, nil
}

func (r *UsersRepo) UpdatePassword(ctx context.Context, userID int64, toUpdate domain.UserToUpdate) (domain.User, error) {
	var user domain.User

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if toUpdate.FirstName != nil {
		setValues = append(setValues, fmt.Sprintf("first_name=$%d", argId))
		args = append(args, *toUpdate.FirstName)
		argId++
	}

	if toUpdate.LastName != nil {
		setValues = append(setValues, fmt.Sprintf("last_name=$%d", argId))
		args = append(args, *toUpdate.LastName)
		argId++
	}

	if toUpdate.Password != nil {
		setValues = append(setValues, fmt.Sprintf("password=$%d", argId))
		args = append(args, *toUpdate.Password)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE users u SET %s WHERE u.id = $%d RETURNING u.*`, setQuery, argId)
	args = append(args, userID)

	err := r.db.GetContext(ctx, &user, query, args...)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

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

func (r *SessionsRepo) Update(ctx context.Context, toUpdate domain.SessionToUpdate, userID int64) (domain.Session, error) {
	var session domain.Session

	err := r.db.GetContext(ctx, &session,
		"UPDATE sessions s SET refresh_token = $1, expires_at = $2 WHERE s.user_id =$3 RETURNING *",
		toUpdate.RefreshToken, toUpdate.ExpiresAt, userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return session, ErrSessionNotFound
		}

		return session, err
	}

	return session, nil
}
