package repositories

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type UserDb struct {
	ID             string    `db:"id"`
	Email          string    `db:"user_id"`
	HashedPassword string    `db:"hashed_password"`
	ActiveSession  int       `db:"active_session"`
	IsActive       bool      `db:"is_active"`
	IsAdmin        bool      `db:"is_admin"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type UserRepository struct {
	db *sqlx.DB
}

// func (d *UserRepository) Save(ctx context.Context , u entities.User) error{

// }
