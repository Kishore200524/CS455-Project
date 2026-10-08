package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	RoleStudent       = "student"
	RoleModerator     = "moderator"
	RoleAdministrator = "administrator"
	SessionDuration   = 12 * time.Hour
	minPasswordLength = 12
	maxPasswordBytes  = 72
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email is already registered")
	ErrInvalidSession     = errors.New("invalid or expired session")
	ErrInvalidInput       = errors.New("invalid authentication input")
)

type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Email     string    `json:"email" bson:"email"`
	Role      string    `json:"role" bson:"role"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type Session struct {
	AccessToken string    `json:"accessToken"`
	TokenType   string    `json:"tokenType"`
	ExpiresAt   time.Time `json:"expiresAt"`
	User        User      `json:"user"`
}

type Repository interface {
	CreateUser(context.Context, string, string, string) (User, error)
	FindUserByEmail(context.Context, string) (User, string, error)
	CreateSession(context.Context, string, string, time.Time) error
	FindUserBySession(context.Context, string, time.Time) (User, error)
	DeleteSession(context.Context, string) error
}

type Service struct {
	repository Repository
	now        func() time.Time
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository, now: time.Now}
}

func (s *Service) RegisterStudent(ctx context.Context, email, password string) (Session, error) {
	email, err := normalizeEmail(email)
	if err != nil || validatePassword(password) != nil {
		return Session{}, ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Session{}, fmt.Errorf("hash student password: %w", err)
	}
	user, err := s.repository.CreateUser(ctx, email, string(hash), RoleStudent)
	if err != nil {
		return Session{}, err
	}

	return s.createSession(ctx, user)
}

func (s *Service) ProvisionModerator(ctx context.Context, email, password string) (User, error) {
	email, err := normalizeEmail(email)
	if err != nil || validatePassword(password) != nil {
		return User{}, ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("hash moderator password: %w", err)
	}
	user, err := s.repository.CreateUser(ctx, email, string(hash), RoleModerator)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Service) CreateAdministrator(ctx context.Context, email, password string) (User, error) {
	email, err := normalizeEmail(email)
	if err != nil || validatePassword(password) != nil {
		return User{}, ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("hash administrator password: %w", err)
	}
	user, err := s.repository.CreateUser(ctx, email, string(hash), RoleAdministrator)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (Session, error) {
	email, err := normalizeEmail(email)
	if err != nil || password == "" || len(password) > maxPasswordBytes {
		return Session{}, ErrInvalidCredentials
	}

	user, hash, err := s.repository.FindUserByEmail(ctx, email)
	if err != nil {
		return Session{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return Session{}, ErrInvalidCredentials
	}
	return s.createSession(ctx, user)
}

func (s *Service) Authenticate(ctx context.Context, token string) (User, error) {
	if len(token) < 32 || len(token) > 128 {
		return User{}, ErrInvalidSession
	}
	user, err := s.repository.FindUserBySession(ctx, hashToken(token), s.now().UTC())
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	if len(token) < 32 || len(token) > 128 {
		return ErrInvalidSession
	}
	if err := s.repository.DeleteSession(ctx, hashToken(token)); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (s *Service) createSession(ctx context.Context, user User) (Session, error) {
	rawToken := make([]byte, 32)
	if _, err := rand.Read(rawToken); err != nil {
		return Session{}, fmt.Errorf("generate session token: %w", err)
	}

	token := base64.RawURLEncoding.EncodeToString(rawToken)
	expiresAt := s.now().UTC().Add(SessionDuration)
	if err := s.repository.CreateSession(ctx, user.ID, hashToken(token), expiresAt); err != nil {
		return Session{}, fmt.Errorf("save session: %w", err)
	}
	return Session{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt,
		User:        user,
	}, nil
}

func normalizeEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if len(email) > 254 || strings.ContainsAny(email, "\r\n") {
		return "", ErrInvalidInput
	}
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email || !strings.Contains(email, "@") {
		return "", ErrInvalidInput
	}
	return email, nil
}

func validatePassword(password string) error {
	length := utf8.RuneCountInString(password)
	if length < minPasswordLength || len(password) > maxPasswordBytes {
		return ErrInvalidInput
	}
	return nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
