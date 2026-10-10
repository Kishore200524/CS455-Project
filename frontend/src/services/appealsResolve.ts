// Person B owns this file: API calls for the ticket detail page
// (SCRUM-12, SCRUM-13). Endpoint shapes are in docs/moderator-api.md.
//
// Implement with apiRequest from ./appealsShared and validate the response with
// isRecord, as services/feedback.ts does. Let ApiError propagate: a 403 or 404
// from fetchTicket means the lock expired or the ticket is no longer this
// moderator's, which TicketDetailPage should report via onDone instead of
// crashing.

import type { TicketDetail } from "./appealsShared";

// GET /api/v1/appeals/{id} -> TicketDetail
export async function fetchTicket(
  _accessToken: string,
  _ticketId: string,
): Promise<TicketDetail> {
  throw new Error("fetchTicket is not implemented yet.");
}

// POST /api/v1/appeals/{id}/release -> 204
export async function releaseTicket(
  _accessToken: string,
  _ticketId: string,
): Promise<void> {
  throw new Error("releaseTicket is not implemented yet.");
}

// POST /api/v1/appeals/{id}/resolve with { action: "restore" } -> 204
export async function restoreTicket(
  _accessToken: string,
  _ticketId: string,
): Promise<void> {
  throw new Error("restoreTicket is not implemented yet.");
}

// POST /api/v1/appeals/{id}/resolve with { action: "delete", comment } -> 204
// The comment is required and is what the student later sees.
export async function deleteTicket(
  _accessToken: string,
  _ticketId: string,
  _comment: string,
): Promise<void> {
  throw new Error("deleteTicket is not implemented yet.");
}
