package httpapi

// Person A owns this file: the queue, claim and active-appeals handlers
// (SCRUM-10, SCRUM-11). Replace each stub body; keep the method names, because
// moderator_routes.go refers to them. Request and response shapes are in
// docs/moderator-api.md.
//
// Use the signed-in moderator with userFromContext(r.Context()).ID and the
// store with s.appealStore. Reuse writeJSON and decodeRequestJSON.

import "net/http"

// listAppeals handles GET /api/v1/appeals?sort=asc|desc.
// Default sort is asc (oldest first). Reject any other sort value with 400.
// Responds 200 with {"appeals": []appeals.Summary}; an empty queue is an empty
// array, never null.
func (s *Server) listAppeals(w http.ResponseWriter, _ *http.Request) {
	writeNotImplemented(w)
}

// listMyAppeals handles GET /api/v1/appeals/mine.
// Responds 200 with {"appeals": []appeals.Summary} for the signed-in moderator.
func (s *Server) listMyAppeals(w http.ResponseWriter, _ *http.Request) {
	writeNotImplemented(w)
}

// claimAppeal handles POST /api/v1/appeals/{id}/claim (read the id with
// r.PathValue("id")). Responds 200 with the claimed appeals.Summary,
// 409 {"error": "This ticket has already been claimed"} when
// appeals.ErrAlreadyClaimed, and 404 when appeals.ErrNotFound.
func (s *Server) claimAppeal(w http.ResponseWriter, _ *http.Request) {
	writeNotImplemented(w)
}
