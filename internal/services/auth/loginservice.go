package auth

import (
	"context"
	"crypto/rsa"
	"net/mail"
	"time"

	"github.com/Syrix42/link-shortener/internal/domain"
	crypto "github.com/Syrix42/link-shortener/internal/infra/crypto/tokenfactory"
	"github.com/Syrix42/link-shortener/internal/services/repositories"
	"github.com/google/uuid"
)

type LoginService struct {
	UserRepo             repositories.UserRepository
	Comparer             Comparer
	QuerySession         SessionQueryRepository
	SessionRepo          repositories.SessionRepository
	RefreshPrivateSecret *rsa.PrivateKey
	AcesssPrivateSecret  *rsa.PrivateKey
}

func NewLoginService(Userrepository repositories.UserRepository,
	Comparer Comparer,
	QuerySession SessionQueryRepository,
	SessionRepostiry repositories.SessionRepository,
	RefreshKey *rsa.PrivateKey, AccessKey *rsa.PrivateKey) *LoginService {
	return &LoginService{
		UserRepo:             Userrepository,
		Comparer:             Comparer,
		QuerySession:         QuerySession,
		SessionRepo:          SessionRepostiry,
		RefreshPrivateSecret: RefreshKey,
		AcesssPrivateSecret:  AccessKey,
	}
}

func (l *LoginService) Login(ctx context.Context, Email, Password string) (string, string, error) {
	ValidEmail, err := mail.ParseAddress(Email)

	if err != nil {
		return "", "", ErrInvalidEmailFormat
	}

	existance, err := l.UserRepo.GetByEmail(ctx, ValidEmail.Address)
	if existance == nil {

		return "", "", ErrUserNotFound

	}
	if err != nil {

		return "", "", err
	}
	err = l.Comparer.Compare(ctx, existance.HashedPassword, Password)

	if err != nil {

		return "", "", ErrInvalidPassword
	}
	ActiveSessions, err := l.QuerySession.CountSessionsByUserID(ctx, existance.ID)

	if err != nil {

		return "", "", err
	}
	err = domain.EnsureMaxActiveSessions(ActiveSessions, 5)
	if err != nil {
		return "", "", ErrTooManyActiveSessions
	}
	SessionId := uuid.NewString()
	AccessToken, err := crypto.IssueAccessToken(ctx, existance.IsAdmin, SessionId, existance.ID, l.AcesssPrivateSecret)
	if err != nil {
		return "", "", err
	}
	RefreshToken, err := crypto.IssueRefreshToken(ctx, SessionId, existance.ID, l.RefreshPrivateSecret)
	if err != nil {

		return "", "", err
	}
	Session := domain.NewSession(SessionId, existance.ID, RefreshToken, time.Now().Add(time.Minute*10080), time.Now())
	err = l.SessionRepo.CreateSession(ctx, Session)
	if err != nil {

		return "", "", err
	}
	return AccessToken, RefreshToken, nil

}
