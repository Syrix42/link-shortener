package repositories

import (
	"context"
	"time"

	"github.com/Syrix42/link-shortener/internal/domain"
	"github.com/jmoiron/sqlx"
)

type SessionDB struct {
	ID          string     `db:"id"`
	UserID      string     `db:"user_id"`
	RefreshHash string     `db:"refresh_token_hashed"`
	ExpiresAt   time.Time  `db:"expires_at"`
	RevokedAt   *time.Time `db:"revoked_at"`
	CreatedAt   time.Time  `db:"created_at"`
	LastUsedAt  *time.Time `db:"last_used_at"`
}

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(database *sqlx.DB) *SessionRepository {
	return &SessionRepository{
		db: database,
	}

}

func (d *SessionRepository) CreateSession(ctx context.Context, s *domain.Session) error {
	session := SessionDB{
		ID:          s.ID,
		UserID:      s.UserID,
		RefreshHash: s.RefreshHash,
		ExpiresAt:   s.ExpiresAt,
		RevokedAt:   s.RevokedAt,
		CreatedAt:   s.CreatedAt,
		LastUsedAt:  s.LastUsedAt,
	}

	query := "INSERT INTO sessions (id, user_id, refresh_token_hash, created_at, last_used_at, revoked_at, expires_at) VALUES (:id, :user_id, :refresh_token_hash, :created_at, :last_used_at, :revoked_at, :expires_at)"
	_, err := d.db.NamedExecContext(ctx, query, session)
	return err
}

func (d *SessionRepository) RotateRefreshToken(ctx context.Context, s *domain.Session) error {
	session := SessionDB{
		ID:          s.ID,
		RefreshHash: s.RefreshHash,
		ExpiresAt:   s.ExpiresAt,
		RevokedAt:   s.RevokedAt,
		LastUsedAt:  s.LastUsedAt,
	}

	query := "UPDATE sessions SET refresh_token_hash = :refresh_token_hash, revoked_at = :revoked_at, expires_at = :expires_at, last_used_at = :last_used_at WHERE id = :id"
	_, err := d.db.NamedExecContext(ctx, query, session)
	return err
}

func (d *SessionRepository) DeleteSession(ctx context.Context, SessionId string) error {
	query := "DELETE FROM sessions WHERE id = :id"
	_, err := d.db.NamedExecContext(ctx, query, SessionId)
	return err
}

func (d *SessionRepository) CountSessionsByUserID(ctx context.Context, userID string) (int, error) {
	query := "SELECT COUNT(*) FROM sessions WHERE user_id = $1"
	var count int
	if err := d.db.GetContext(ctx, &count, query, userID); err != nil {
		return 0, err
	}
	return count, nil
}
