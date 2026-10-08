package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Kishore200524/CS455-Project/backend/internal/auth"
)

func TestRegistrationLoginAndAuthenticatedIdentity(t *testing.T) {
	service := auth.NewService(newMemoryAuthRepository())
	handler := NewServer(nil, service)

	register := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(
		`{"email":"Student@College.edu","password":"long-secure-password"}`,
	))
	registerResponse := httptest.NewRecorder()
	handler.ServeHTTP(registerResponse, register)
	if registerResponse.Code != http.StatusCreated {
		t.Fatalf("registration status = %d, want %d: %s", registerResponse.Code, http.StatusCreated, registerResponse.Body.String())
	}

	var session auth.Session
	if err := json.NewDecoder(registerResponse.Body).Decode(&session); err != nil {
		t.Fatalf("decode registration response: %v", err)
	}
	if session.User.Role != auth.RoleStudent || session.User.Email != "student@college.edu" {
		t.Fatalf("unexpected registered user: %+v", session.User)
	}
	if strings.Contains(registerResponse.Body.String(), "password") {
		t.Fatal("authentication response must not include password data")
	}

	me := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	me.Header.Set("Authorization", "Bearer "+session.AccessToken)
	meResponse := httptest.NewRecorder()
	handler.ServeHTTP(meResponse, me)
	if meResponse.Code != http.StatusOK {
		t.Fatalf("identity status = %d, want %d", meResponse.Code, http.StatusOK)
	}

	login := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(
		`{"email":"student@college.edu","password":"long-secure-password"}`,
	))
	loginResponse := httptest.NewRecorder()
	handler.ServeHTTP(loginResponse, login)
	if loginResponse.Code != http.StatusOK {
		t.Fatalf("login status = %d, want %d: %s", loginResponse.Code, http.StatusOK, loginResponse.Body.String())
	}
}

func TestProtectedRoutesRequireRole(t *testing.T) {
	service := auth.NewService(newMemoryAuthRepository())
	studentSession, err := service.RegisterStudent(context.Background(), "student@college.edu", "long-secure-password")
	if err != nil {
		t.Fatal(err)
	}
	handler := NewServer(nil, service)

	noToken := httptest.NewRequest(http.MethodPost, "/api/v1/admin/moderators", strings.NewReader(
		`{"email":"moderator@college.edu","password":"long-secure-password"}`,
	))
	noTokenResponse := httptest.NewRecorder()
	handler.ServeHTTP(noTokenResponse, noToken)
	if noTokenResponse.Code != http.StatusUnauthorized {
		t.Fatalf("missing-token status = %d, want %d", noTokenResponse.Code, http.StatusUnauthorized)
	}

	roleEscalation := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(
		`{"email":"attacker@college.edu","password":"long-secure-password","role":"administrator"}`,
	))
	roleEscalationResponse := httptest.NewRecorder()
	handler.ServeHTTP(roleEscalationResponse, roleEscalation)
	if roleEscalationResponse.Code != http.StatusBadRequest {
		t.Fatalf("role escalation status = %d, want %d", roleEscalationResponse.Code, http.StatusBadRequest)
	}

	studentRequest := httptest.NewRequest(http.MethodPost, "/api/v1/admin/moderators", strings.NewReader(
		`{"email":"moderator@college.edu","password":"long-secure-password"}`,
	))
	studentRequest.Header.Set("Authorization", "Bearer "+studentSession.AccessToken)
	studentResponse := httptest.NewRecorder()
	handler.ServeHTTP(studentResponse, studentRequest)
	if studentResponse.Code != http.StatusForbidden {
		t.Fatalf("student moderator-provision status = %d, want %d", studentResponse.Code, http.StatusForbidden)
	}
}

func TestAdministratorCanProvisionModerator(t *testing.T) {
	service := auth.NewService(newMemoryAuthRepository())
	if _, err := service.CreateAdministrator(context.Background(), "admin@college.edu", "long-secure-password"); err != nil {
		t.Fatal(err)
	}
	adminSession, err := service.Login(context.Background(), "admin@college.edu", "long-secure-password")
	if err != nil {
		t.Fatal(err)
	}
	handler := NewServer(nil, service)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/moderators", strings.NewReader(
		`{"email":"moderator@college.edu","password":"moderator-secure-password"}`,
	))
	request.Header.Set("Authorization", "Bearer "+adminSession.AccessToken)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("provision status = %d, want %d: %s", response.Code, http.StatusCreated, response.Body.String())
	}

	var moderator auth.User
	if err := json.NewDecoder(response.Body).Decode(&moderator); err != nil {
		t.Fatalf("decode moderator response: %v", err)
	}
	if moderator.Role != auth.RoleModerator {
		t.Fatalf("provisioned role = %q, want %q", moderator.Role, auth.RoleModerator)
	}

	moderatorSession, err := service.Login(context.Background(), moderator.Email, "moderator-secure-password")
	if err != nil {
		t.Fatalf("moderator login: %v", err)
	}
	feedbackRequest := httptest.NewRequest(http.MethodPost, "/api/v1/feedback", strings.NewReader(
		`{"courseId":"CS455","content":"Feedback"}`,
	))
	feedbackRequest.Header.Set("Authorization", "Bearer "+moderatorSession.AccessToken)
	feedbackResponse := httptest.NewRecorder()
	handler.ServeHTTP(feedbackResponse, feedbackRequest)
	if feedbackResponse.Code != http.StatusForbidden {
		t.Fatalf("moderator feedback status = %d, want %d", feedbackResponse.Code, http.StatusForbidden)
	}
}
