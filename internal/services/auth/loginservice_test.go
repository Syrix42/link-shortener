package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"
	"time"

	"github.com/Syrix42/link-shortener/internal/domain"
)

// ---- Login fakes (prefixed to avoid clashing with register fakes) ----

type loginFakeUserRepo struct {
	getByEmailFn func(ctx context.Context, email string) (*domain.User, error)

	getByEmailCalled int
	getByEmailArg    string
}

func (f *loginFakeUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	f.getByEmailCalled++
	f.getByEmailArg = email
	return f.getByEmailFn(ctx, email)
}

func (f *loginFakeUserRepo) Save(ctx context.Context, u domain.User) error {
	// not used by Login
	return nil
}

type loginFakeComparer struct {
	compareFn     func(ctx context.Context, hashed, plaintext string) error
	compareCalled int
	hashedArg     string
	plainArg      string
}

func (f *loginFakeComparer) Compare(ctx context.Context, hashed string, plaintext string) error {
	f.compareCalled++
	f.hashedArg = hashed
	f.plainArg = plaintext
	return f.compareFn(ctx, hashed, plaintext)
}

type loginFakeSessionQueryRepo struct {
	countFn     func(ctx context.Context, userID string) (int, error)
	countCalled int
	userIDArg   string
}

func (f *loginFakeSessionQueryRepo) CountSessionsByUserID(ctx context.Context, userID string) (int, error) {
	f.countCalled++
	f.userIDArg = userID
	return f.countFn(ctx, userID)
}

type loginFakeSessionRepo struct {
	createFn     func(ctx context.Context, s *domain.Session) error
	createCalled int
	created      *domain.Session
}

func (f *loginFakeSessionRepo) CreateSession(ctx context.Context, s *domain.Session) error {
	f.createCalled++
	f.created = s
	return f.createFn(ctx, s)
}

func (f *loginFakeSessionRepo) RotateRefreshToken(ctx context.Context, s *domain.Session) error {
	return nil
}

func (f *loginFakeSessionRepo) DeleteSession(ctx context.Context, sessionID string) error {
	return nil
}

// ---- Helpers ----

func mustRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey: %v", err)
	}
	return k
}

// ---- Tests ----

func TestLogin_Success_IssuesTokens_AndCreatesSession(t *testing.T) {
	refreshKey := mustRSAKey(t)
	accessKey := mustRSAKey(t)

	u := &domain.User{
		ID:             "user-1",
		Email:          "ali@example.com",
		HashedPassword: "hashed_pw123",
		IsAdmin:        false,
	}

	userRepo := &loginFakeUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			if email != "ali@example.com" {
				t.Fatalf("expected email ali@example.com, got %q", email)
			}
			return u, nil
		},
	}

	comparer := &loginFakeComparer{
		compareFn: func(ctx context.Context, hashed, plaintext string) error {
			if hashed != "hashed_pw123" {
				t.Fatalf("expected hashed_pw123, got %q", hashed)
			}
			if plaintext != "pw123" {
				t.Fatalf("expected pw123, got %q", plaintext)
			}
			return nil
		},
	}

	query := &loginFakeSessionQueryRepo{
		countFn: func(ctx context.Context, userID string) (int, error) {
			if userID != "user-1" {
				t.Fatalf("expected userID user-1, got %q", userID)
			}
			return 0, nil
		},
	}

	sessionRepo := &loginFakeSessionRepo{
		createFn: func(ctx context.Context, s *domain.Session) error {
			return nil
		},
	}

	svc := NewLoginService(userRepo, comparer, query, sessionRepo, refreshKey, accessKey)

	accessTok, refreshTok, err := svc.Login(context.Background(), "ali@example.com", "pw123")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if accessTok == "" || refreshTok == "" {
		t.Fatalf("expected non-empty tokens, got access=%q refresh=%q", accessTok, refreshTok)
	}

	if userRepo.getByEmailCalled != 1 {
		t.Fatalf("expected GetByEmail called once, got %d", userRepo.getByEmailCalled)
	}
	if comparer.compareCalled != 1 {
		t.Fatalf("expected Compare called once, got %d", comparer.compareCalled)
	}
	if query.countCalled != 1 {
		t.Fatalf("expected CountSessionsByUserID called once, got %d", query.countCalled)
	}
	if sessionRepo.createCalled != 1 {
		t.Fatalf("expected CreateSession called once, got %d", sessionRepo.createCalled)
	}

	// Session invariants we can check without knowing internal token format
	s := sessionRepo.created
	if s == nil {
		t.Fatalf("expected session to be created")
	}
	if s.UserID != "user-1" {
		t.Fatalf("expected session.UserID=user-1, got %q", s.UserID)
	}
	if s.RefreshHash != refreshTok {
		t.Fatalf("expected session.RefreshToken to equal returned refresh token")
	}
	if s.ID == "" {
		t.Fatalf("expected non-empty session ID")
	}
	now := time.Now()
	min := now.Add(7*24*time.Hour - 2*time.Second) // 2s tolerance for test runtime
	if s.ExpiresAt.Before(min) {
		t.Fatalf("expected expires in ~>=7 days, got %v (min acceptable %v)", s.ExpiresAt, min)
	}
}

func TestLogin_InvalidEmail_ReturnsErrInvalidEmailFormat_AndDoesNotTouchDeps(t *testing.T) {
	refreshKey := mustRSAKey(t)
	accessKey := mustRSAKey(t)

	userRepo := &loginFakeUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			t.Fatalf("GetByEmail should not be called on invalid email")
			return nil, nil
		},
	}
	comparer := &loginFakeComparer{
		compareFn: func(ctx context.Context, hashed, plaintext string) error {
			t.Fatalf("Compare should not be called on invalid email")
			return nil
		},
	}
	query := &loginFakeSessionQueryRepo{
		countFn: func(ctx context.Context, userID string) (int, error) {
			t.Fatalf("CountSessionsByUserID should not be called on invalid email")
			return 0, nil
		},
	}
	sessionRepo := &loginFakeSessionRepo{
		createFn: func(ctx context.Context, s *domain.Session) error {
			t.Fatalf("CreateSession should not be called on invalid email")
			return nil
		},
	}

	svc := NewLoginService(userRepo, comparer, query, sessionRepo, refreshKey, accessKey)

	_, _, err := svc.Login(context.Background(), "not-an-email", "pw123")
	if !errors.Is(err, ErrInvalidEmailFormat) {
		t.Fatalf("expected ErrInvalidEmailFormat, got %v", err)
	}
}

func TestLogin_UserNotFound_ReturnsErrUserNotFound(t *testing.T) {
	refreshKey := mustRSAKey(t)
	accessKey := mustRSAKey(t)

	userRepo := &loginFakeUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, nil // not found
		},
	}
	comparer := &loginFakeComparer{
		compareFn: func(ctx context.Context, hashed, plaintext string) error {
			t.Fatalf("Compare should not be called when user not found")
			return nil
		},
	}
	query := &loginFakeSessionQueryRepo{
		countFn: func(ctx context.Context, userID string) (int, error) {
			t.Fatalf("CountSessionsByUserID should not be called when user not found")
			return 0, nil
		},
	}
	sessionRepo := &loginFakeSessionRepo{
		createFn: func(ctx context.Context, s *domain.Session) error {
			t.Fatalf("CreateSession should not be called when user not found")
			return nil
		},
	}

	svc := NewLoginService(userRepo, comparer, query, sessionRepo, refreshKey, accessKey)

	_, _, err := svc.Login(context.Background(), "ali@example.com", "pw123")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestLogin_InvalidPassword_ReturnsErrInvalidPassword(t *testing.T) {
	refreshKey := mustRSAKey(t)
	accessKey := mustRSAKey(t)

	u := &domain.User{
		ID:             "user-1",
		Email:          "ali@example.com",
		HashedPassword: "hashed_pw123",
		IsAdmin:        false,
	}

	userRepo := &loginFakeUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return u, nil
		},
	}

	comparer := &loginFakeComparer{
		compareFn: func(ctx context.Context, hashed, plaintext string) error {
			return errors.New("mismatch")
		},
	}

	query := &loginFakeSessionQueryRepo{
		countFn: func(ctx context.Context, userID string) (int, error) {
			t.Fatalf("CountSessionsByUserID should not be called when password invalid")
			return 0, nil
		},
	}
	sessionRepo := &loginFakeSessionRepo{
		createFn: func(ctx context.Context, s *domain.Session) error {
			t.Fatalf("CreateSession should not be called when password invalid")
			return nil
		},
	}

	svc := NewLoginService(userRepo, comparer, query, sessionRepo, refreshKey, accessKey)

	_, _, err := svc.Login(context.Background(), "ali@example.com", "wrongpw")
	if !errors.Is(err, ErrInvalidPassword) {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
}

func TestLogin_TooManyActiveSessions_ReturnsErrTooManyActiveSessions(t *testing.T) {
	refreshKey := mustRSAKey(t)
	accessKey := mustRSAKey(t)

	u := &domain.User{
		ID:             "user-1",
		Email:          "ali@example.com",
		HashedPassword: "hashed_pw123",
		IsAdmin:        false,
	}

	userRepo := &loginFakeUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return u, nil
		},
	}

	comparer := &loginFakeComparer{
		compareFn: func(ctx context.Context, hashed, plaintext string) error {
			return nil
		},
	}

	query := &loginFakeSessionQueryRepo{
		countFn: func(ctx context.Context, userID string) (int, error) {
			return 5, nil // maxed out
		},
	}

	sessionRepo := &loginFakeSessionRepo{
		createFn: func(ctx context.Context, s *domain.Session) error {
			t.Fatalf("CreateSession should not be called when too many sessions")
			return nil
		},
	}

	svc := NewLoginService(userRepo, comparer, query, sessionRepo, refreshKey, accessKey)

	_, _, err := svc.Login(context.Background(), "ali@example.com", "pw123")
	if !errors.Is(err, ErrTooManyActiveSessions) {
		t.Fatalf("expected ErrTooManyActiveSessions, got %v", err)
	}
}
