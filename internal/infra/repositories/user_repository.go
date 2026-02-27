package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Syrix42/link-shortener/internal/entities"
	"github.com/jmoiron/sqlx"
)

type UserDb struct {
	ID             string    `db:"id"`
	Email          string    `db:"email"`
	HashedPassword string    `db:"hashed_password"`
	ActiveSession  int       `db:"active_sessions"`
	IsActive       bool      `db:"is_active"`
	IsAdmin        bool      `db:"is_admin"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(database *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: database,
	}

}

func (d *UserRepository) Save(ctx context.Context, u entities.User) error {

	user := UserDb{
		ID:             u.ID,
		Email:          u.Email,
		HashedPassword: u.HashedPassword,
		IsActive:       u.IsActive,
		IsAdmin:        u.IsAdmin,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
	query := "Insert Into users (id , email , hashed_password , active_sessions , is_active , is_admin , created_at , updated_at) VALUES(:id , :email ,:hashed_password ,:active_sessions , :is_active , :is_admin , :created_at , :updated_at) "

	_, err := d.db.NamedExec(query, user)

	if err != nil {
		return err
	}
	return nil

}

func (d *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {

	user := UserDb{}
	query := "SELECT * FROM users WHERE email = $1"

	err := d.db.Get(&user, query, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err // real DB error
		}
	}
	return &entities.User{
		ID:             user.ID,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		IsActive:       user.IsActive,
		IsAdmin:        user.IsAdmin,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}, nil
}
