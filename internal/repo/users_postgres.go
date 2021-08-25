package repo

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lotostudio/financial-api/internal/domain"
	log "github.com/sirupsen/logrus"
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
