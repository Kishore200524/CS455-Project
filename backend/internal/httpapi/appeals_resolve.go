package httpapi

// Person B owns this file: the detail, release and resolve handlers
// (SCRUM-12, SCRUM-13). Replace each stub body; keep the method names, because
// moderator_routes.go refers to them. Request and response shapes are in
// docs/moderator-api.md.
//
// Use the signed-in moderator with userFromContext(r.Context()).ID and the
// store with s.appealStore. Reuse writeJSON and decodeRequestJSON.

import "net/http"

// getAppeal handles GET /api/v1/appeals/{id} (read the id with
// r.PathValue("id")). Responds 200 with the full appeals.Ticket only for the
// moderator holding an unexpired lock; 403 for appeals.ErrNotLockHolder (this
// includes an expired lock, so the page can show "no longer assigned to you")
// and 404 for appeals.ErrNotFound.
func (s *Server) getAppeal(w http.ResponseWriter, _ *http.Request) {
	writeNotImplemented(w)
}

// releaseAppeal handles POST /api/v1/appeals/{id}/release (no body).
// Responds 204 on success, 403 for appeals.ErrNotLockHolder.
func (s *Server) releaseAppeal(w http.ResponseWriter, _ *http.Request) {
	writeNotImplemented(w)
}

// resolveAppeal handles POST /api/v1/appeals/{id}/resolve with
// {"action": "restore" | "delete", "comment": "..."}. Decode it with
// decodeRequestJSON. Responds 204 on success; 400 for appeals.ErrInvalidAction
// or appeals.ErrCommentRequired (delete needs a non-empty comment); 403 for
// appeals.ErrNotLockHolder.
func (s *Server) resolveAppeal(w http.ResponseWriter, _ *http.Request) {
	writeNotImplemented(w)
}
