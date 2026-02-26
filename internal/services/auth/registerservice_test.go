package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Syrix42/link-shortener/internal/entities"
)

// ---- Fakes ----

type fakeUserRepo struct {
	GetByEmailFn func(ctx context.Context, email string) (*entities.User, error)
	saveFn       func(ctx context.Context, u entities.User) error

	// calls
	GetByEmailCalled int
	GetByEmailArgs   string

	saveCalled int
	savedUser  entities.User
}

func (f *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	f.GetByEmailCalled++
	f.GetByEmailArgs = email
	return f.GetByEmailFn(ctx, email)
}

func (f *fakeUserRepo) Save(ctx context.Context, u entities.User) error {
	f.saveCalled++
	f.savedUser = u
	return f.saveFn(ctx, u)

}

type fakeHahser struct {
	hashFn      func(ctx context.Context, password string) (string, error)
	hashCalled  int
	hashArgPass string
}

func (f *fakeHahser) Hash(ctx context.Context, password string) (string, error) {
	f.hashCalled++
	f.hashArgPass = password
	return f.hashFn(ctx, password)
}

// Tests

func TestRegister_InvalidEmail_ReturnsErrInvalidEmailFormat_AndDoesNotTouchRepo(t *testing.T) {
	repo := &fakeUserRepo{
		GetByEmailFn: func(ctx context.Context, email string) (*entities.User, error) {
			t.Fatalf("Getby Email should not be called for invalid email")
			return nil, nil
		},
		saveFn: func(ctx context.Context, u entities.User) error {
			t.Fatalf("Save Should not be called for invalid email ")
			return nil
		},
	}
	hasher := &fakeHahser{
		hashFn: func(ctx context.Context, password string) (string, error) {
			t.Fatalf("Hahser . Hash should not be called for invalid email")
			return "", nil
		},
	}
	svc := NewRegisterService(repo, hasher)

	err := svc.Register(context.Background(), "not-an-email", "pw123")

	if !errors.Is(err, ErrInvalidEmailFormat) {
		t.Fatalf("expected ErrInvalidEmailFormat , got: %v", err)
	}
	if repo.GetByEmailCalled != 0 || repo.saveCalled != 0 {
		t.Fatalf("expected repo not called, got GetByEmail=%d Save=%d", repo.GetByEmailCalled, repo.saveCalled)
	}
	if hasher.hashCalled != 0 {
		t.Fatalf("expected hasher not called, got %d", hasher.hashCalled)
	}

}

func TestRegister_UserAlreadyExists_ReturnsErrUserAlreadyExists_AndDoesNotSave(t *testing.T) {
	existing := entities.NewUser("id1", "ali@example.com", "hash", true, false, 0, time.Now().UTC(), time.Now().UTC())

	repo := &fakeUserRepo{
		GetByEmailFn: func(ctx context.Context, email string) (*entities.User, error) {
			return existing, nil
		},
		saveFn: func(ctx context.Context, u entities.User) error {
			t.Fatalf("Save should not be called when user exists")
			return nil
		},
	}
	hasher := &fakeHahser{
		hashFn: func(ctx context.Context, password string) (string, error) {
			t.Fatalf("Hasher.Hash should not be called when user exists")
			return "", nil
		},
	}

	svc := NewRegisterService(repo, hasher)

	err := svc.Register(context.Background(), "ali@example.com", "pw123")
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf("expected ErrUserAlreadyExists, got: %v", err)
	}
	if repo.saveCalled != 0 {
		t.Fatalf("expected Save not called, got %d", repo.saveCalled)
	}
}

func TestRegister_GetByEmailError_ReturnsThatError_AndDoesNotHashOrSave(t *testing.T) {
	repoErr := errors.New("db down")

	repo := &fakeUserRepo{
		GetByEmailFn: func(ctx context.Context, email string) (*entities.User, error) {
			return nil, repoErr
		},
		saveFn: func(ctx context.Context, u entities.User) error {
			t.Fatalf("Save should not be called when GetByEmail fails")
			return nil
		},
	}
	hasher := &fakeHahser{
		hashFn: func(ctx context.Context, password string) (string, error) {
			t.Fatalf("Hasher.Hash should not be called when GetByEmail fails")
			return "", nil
		},
	}

	svc := NewRegisterService(repo, hasher)

	err := svc.Register(context.Background(), "ali@example.com", "pw123")
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repoErr, got: %v", err)
	}
	if hasher.hashCalled != 0 {
		t.Fatalf("expected hasher not called, got %d", hasher.hashCalled)
	}
	if repo.saveCalled != 0 {
		t.Fatalf("expected Save not called, got %d", repo.saveCalled)
	}
}

func TestRegister_HashError_ReturnsThatError_AndDoesNotSave(t *testing.T) {
	hashErr := errors.New("hash failed")

	repo := &fakeUserRepo{
		GetByEmailFn: func(ctx context.Context, email string) (*entities.User, error) {
			return nil, nil
		},
		saveFn: func(ctx context.Context, u entities.User) error {
			t.Fatalf("Save should not be called when Hash fails")
			return nil
		},
	}
	hasher := &fakeHahser{
		hashFn: func(ctx context.Context, password string) (string, error) {
			return "", hashErr
		},
	}

	svc := NewRegisterService(repo, hasher)

	err := svc.Register(context.Background(), "ali@example.com", "pw123")
	if !errors.Is(err, hashErr) {
		t.Fatalf("expected hashErr, got: %v", err)
	}
	if repo.saveCalled != 0 {
		t.Fatalf("expected Save not called, got %d", repo.saveCalled)
	}
}

func TestRegister_Success_SetsInvariants_AndSavesUser(t *testing.T) {
	repo := &fakeUserRepo{
		GetByEmailFn: func(ctx context.Context, email string) (*entities.User, error) {
			return nil, nil
		},
		saveFn: func(ctx context.Context, u entities.User) error {
			return nil
		},
	}
	hasher := &fakeHahser{
		hashFn: func(ctx context.Context, password string) (string, error) {
			if password != "pw123" {
				t.Fatalf("expected password pw123, got %q", password)
			}
			return "hashed_pw123", nil
		},
	}

	svc := NewRegisterService(repo, hasher)

	err := svc.Register(context.Background(), "ali@example.com", "pw123")
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if repo.saveCalled != 1 {
		t.Fatalf("expected Save called once, got %d", repo.saveCalled)
	}

	u := repo.savedUser

	// Preconditions & invariants you listed
	if u.Email != "ali@example.com" {
		t.Fatalf("expected email ali@example.com, got %q", u.Email)
	}
	if u.IsActive != true {
		t.Fatalf("expected IsActive=true, got %v", u.IsActive)
	}
	if u.IsAdmin != false {
		t.Fatalf("expected IsAdmin=false, got %v", u.IsAdmin)
	}
	if u.ActiveSession != 0 {
		t.Fatalf("expected ActiveSessions=0, got %d", u.ActiveSession)
	}
	// Also verify password is the hashed value
	if u.HashedPassword != "hashed_pw123" {
		t.Fatalf("expected Password=hashed_pw123, got %q", u.HashedPassword)
	}
	// sanity checks for created fields
	if u.ID == "" {
		t.Fatalf("expected non-empty ID")
	}
	if u.CreatedAt.IsZero() || u.UpdatedAt.IsZero() {
		t.Fatalf("expected non-zero timestamps")
	}
	if u.UpdatedAt.Before(u.CreatedAt) {
		t.Fatalf("expected UpdatedAt >= CreatedAt")
	}
}
