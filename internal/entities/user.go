package entities

import "time"

type User struct {
	ID             string
	Email          string
	HashedPassword string
	ActiveSession  int
	IsActive       bool
	IsAdmin        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewUser(Id, Email, HashedPassword string,
	IsActive, IsAdmin bool, ActiveSession int, CreatedAt, UpdatedAt time.Time) *User {

	return &User{
		ID:             Id,
		Email:          Email,
		HashedPassword: HashedPassword,
		ActiveSession:  ActiveSession,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
}
