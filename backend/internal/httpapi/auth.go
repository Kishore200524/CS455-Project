package httpapi

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Kishore200524/CS455-Project/backend/internal/auth"
)

type contextKey string

const userContextKey contextKey = "authenticated-user"

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) registerStudent(w http.ResponseWriter, r *http.Request) {
	var input credentials
	if err := decodeRequestJSON(w, r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "request must contain one valid JSON object"})
		return
	}

	session, err := s.authService.RegisterStudent(r.Context(), input.Email, input.Password)
	if err != nil {
		s.writeAuthError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, session)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var input credentials
	if err := decodeRequestJSON(w, r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "request must contain one valid JSON object"})
		return
	}

	session, err := s.authService.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		s.writeAuthError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	token, ok := bearerToken(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "bearer token required"})
		return
	}
	if err := s.authService.Logout(r.Context(), token); err != nil {
		s.writeAuthError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) currentUser(w http.ResponseWriter, r *http.Request) {
	user, ok := userFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "authenticated user unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (s *Server) provisionModerator(w http.ResponseWriter, r *http.Request) {
	var input credentials
	if err := decodeRequestJSON(w, r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "request must contain one valid JSON object"})
		return
	}

	user, err := s.authService.ProvisionModerator(r.Context(), input.Email, input.Password)
	if err != nil {
		s.writeAuthError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func (s *Server) requireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := s.authenticateRequest(w, r)
		if !ok {
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, user)))
	})
}

func (s *Server) requireRole(role string, next http.Handler) http.Handler {
	return s.requireUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := userFromContext(r.Context())
		if user.Role != role {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "insufficient permissions"})
			return
		}
		next.ServeHTTP(w, r)
	}))
}

func (s *Server) authenticateRequest(w http.ResponseWriter, r *http.Request) (auth.User, bool) {
	token, ok := bearerToken(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "bearer token required"})
		return auth.User{}, false
	}
	user, err := s.authService.Authenticate(r.Context(), token)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidSession) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired session"})
			return auth.User{}, false
		}
		log.Printf("authenticate request: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not authenticate request"})
		return auth.User{}, false
	}
	return user, true
}

func (s *Server) writeAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, auth.ErrInvalidInput):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email or password does not meet requirements"})
	case errors.Is(err, auth.ErrInvalidCredentials):
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
	case errors.Is(err, auth.ErrInvalidSession):
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired session"})
	case errors.Is(err, auth.ErrEmailAlreadyExists):
		writeJSON(w, http.StatusConflict, map[string]string{"error": "email is already registered"})
	default:
		log.Printf("authentication operation: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "authentication operation failed"})
	}
}

func bearerToken(r *http.Request) (string, bool) {
	parts := strings.Fields(r.Header.Get("Authorization"))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}

func userFromContext(ctx context.Context) (auth.User, bool) {
	user, ok := ctx.Value(userContextKey).(auth.User)
	return user, ok
}
