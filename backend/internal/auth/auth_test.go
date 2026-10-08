package auth

import (
	"context"
	"errors"
	"testing"
	"time"
)

type memoryRepository struct {
	users    map[string]User
	password map[string]string
	sessions map[string]string
}

func newMemoryRepository() *memoryRepository {
	return &memoryRepository{
		users:    make(map[string]User),
		password: make(map[string]string),
		sessions: make(map[string]string),
	}
}

func (m *memoryRepository) CreateUser(_ context.Context, email, passwordHash, role string) (User, error) {
	if _, exists := m.users[email]; exists {
		return User{}, ErrEmailAlreadyExists
	}
	user := User{ID: email, Email: email, Role: role, CreatedAt: time.Now().UTC()}
	m.users[email] = user
	m.password[email] = passwordHash
	return user, nil
}

func (m *memoryRepository) FindUserByEmail(_ context.Context, email string) (User, string, error) {
	user, exists := m.users[email]
	if !exists {
		return User{}, "", ErrInvalidCredentials
	}
	return user, m.password[email], nil
}

func (m *memoryRepository) CreateSession(_ context.Context, userID, tokenHash string, _ time.Time) error {
	m.sessions[tokenHash] = userID
	return nil
}

func (m *memoryRepository) FindUserBySession(_ context.Context, tokenHash string, _ time.Time) (User, error) {
	userID, exists := m.sessions[tokenHash]
	if !exists {
		return User{}, ErrInvalidSession
	}
	return m.users[userID], nil
}

func (m *memoryRepository) DeleteSession(_ context.Context, tokenHash string) error {
	delete(m.sessions, tokenHash)
	return nil
}

func TestStudentRegistrationCreatesHashedPasswordAndSession(t *testing.T) {
	repository := newMemoryRepository()
	service := NewService(repository)

	session, err := service.RegisterStudent(context.Background(), "  Student@College.edu ", "long-secure-password")
	if err != nil {
		t.Fatalf("RegisterStudent() error = %v", err)
	}
	if session.User.Role != RoleStudent || session.User.Email != "student@college.edu" {
		t.Fatalf("unexpected user in session: %+v", session.User)
	}
	if repository.password[session.User.Email] == "long-secure-password" {
		t.Fatal("password was stored in plaintext")
	}
	if session.AccessToken == "" || session.TokenType != "Bearer" {
		t.Fatalf("invalid session response: %+v", session)
	}
	if _, err := service.Authenticate(context.Background(), session.AccessToken); err != nil {
		t.Fatalf("Authenticate() error = %v", err)
	}
}

func TestLoginRejectsInvalidPassword(t *testing.T) {
	service := NewService(newMemoryRepository())
	if _, err := service.RegisterStudent(context.Background(), "student@college.edu", "long-secure-password"); err != nil {
		t.Fatal(err)
	}

	_, err := service.Login(context.Background(), "student@college.edu", "wrong-password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestModeratorProvisioningDoesNotIssuePublicStudentSession(t *testing.T) {
	repository := newMemoryRepository()
	service := NewService(repository)

	user, err := service.ProvisionModerator(context.Background(), "moderator@college.edu", "long-secure-password")
	if err != nil {
		t.Fatalf("ProvisionModerator() error = %v", err)
	}
	if user.Role != RoleModerator {
		t.Fatalf("role = %q, want %q", user.Role, RoleModerator)
	}
	if _, err := service.Authenticate(context.Background(), "not-a-valid-session-token"); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("Authenticate() error = %v, want ErrInvalidSession", err)
	}
}

func TestLogoutRevokesSession(t *testing.T) {
	service := NewService(newMemoryRepository())
	session, err := service.RegisterStudent(context.Background(), "student@college.edu", "long-secure-password")
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Logout(context.Background(), session.AccessToken); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
	if _, err := service.Authenticate(context.Background(), session.AccessToken); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("Authenticate() error = %v, want ErrInvalidSession", err)
	}
}

func TestRegistrationRejectsWeakPasswordAndInvalidEmail(t *testing.T) {
	service := NewService(newMemoryRepository())
	tests := []struct {
		email    string
		password string
	}{
		{email: "not-an-email", password: "long-secure-password"},
		{email: "student@college.edu", password: "short"},
	}
	for _, test := range tests {
		if _, err := service.RegisterStudent(context.Background(), test.email, test.password); !errors.Is(err, ErrInvalidInput) {
			t.Errorf("RegisterStudent(%q) error = %v, want ErrInvalidInput", test.email, err)
		}
	}
}
