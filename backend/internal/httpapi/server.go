package httpapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kishore200524/CS455-Project/backend/internal/auth"
	"github.com/Kishore200524/CS455-Project/backend/internal/feedback"
)

type Server struct {
	feedbackStore feedback.Store
	authService   *auth.Service
	cookieSecure  bool
}

func NewServer(feedbackStore feedback.Store, authService *auth.Service, cookieSecure ...bool) http.Handler {
	secure := false
	if len(cookieSecure) > 0 {
		secure = cookieSecure[0]
	}
	server := &Server{feedbackStore: feedbackStore, authService: authService, cookieSecure: secure}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", server.health)
	mux.HandleFunc("POST /api/auth/request-otp", server.requestOTP)
	mux.HandleFunc("POST /api/auth/verify-otp", server.verifyOTP)
	mux.HandleFunc("POST /api/auth/logout", server.logout)
	mux.Handle("GET /api/auth/me", server.requireUser(http.HandlerFunc(server.currentUser)))
	mux.HandleFunc("POST /api/v1/auth/login", server.login)
	mux.HandleFunc("POST /api/auth/student/login", server.loginStudent)
	mux.HandleFunc("POST /api/auth/admin/login", server.loginAdministrator)
	mux.HandleFunc("POST /api/auth/moderator/login", server.loginModerator)
	mux.HandleFunc("POST /api/v1/auth/logout", server.logout)
	mux.Handle("GET /api/v1/auth/me", server.requireUser(http.HandlerFunc(server.currentUser)))
	mux.Handle("POST /api/v1/admin/moderators", server.requireRole(auth.RoleAdministrator, http.HandlerFunc(server.provisionModerator)))
	mux.Handle("DELETE /api/v1/admin/moderators/{email}", server.requireRole(auth.RoleAdministrator, http.HandlerFunc(server.removeModerator)))
	mux.Handle("POST /api/v1/feedback", server.requireRole(auth.RoleStudent, http.HandlerFunc(server.createFeedback)))
	mux.Handle("GET /api/v1/feedback/explore", server.requireRole(auth.RoleStudent, http.HandlerFunc(server.exploreFeedback)))
	mux.Handle("GET /api/v1/feedback/ratings", server.requireRole(auth.RoleStudent, http.HandlerFunc(server.courseRatings)))
	mux.Handle("GET /api/v1/feedback/status/{reference}", server.requireRole(auth.RoleStudent, http.HandlerFunc(server.reviewStatus)))
	mux.Handle("POST /api/v1/feedback/status/{reference}/appeal", server.requireRole(auth.RoleStudent, http.HandlerFunc(server.createAppeal)))
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}
