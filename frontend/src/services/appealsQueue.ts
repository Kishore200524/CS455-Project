// Person A owns this file: API calls for the queue and active-appeals pages
// (SCRUM-10, SCRUM-11). Endpoint shapes are in docs/moderator-api.md.
//
// Implement with apiRequest from ./appealsShared and validate the response with
// isRecord, as services/feedback.ts does. The claim call must let ApiError
// (status 409) propagate so QueuePage can show
// "This ticket has already been claimed".

import type { SortOrder, TicketSummary } from "./appealsShared";

// GET /api/v1/appeals?sort=asc|desc -> { appeals: TicketSummary[] }
export async function fetchQueue(
  _accessToken: string,
  _sort: SortOrder,
): Promise<TicketSummary[]> {
  throw new Error("fetchQueue is not implemented yet.");
}

// POST /api/v1/appeals/{id}/claim -> TicketSummary
export async function claimTicket(
  _accessToken: string,
  _ticketId: string,
): Promise<TicketSummary> {
  throw new Error("claimTicket is not implemented yet.");
}

// GET /api/v1/appeals/mine -> { appeals: TicketSummary[] }
export async function fetchMyTickets(
  _accessToken: string,
): Promise<TicketSummary[]> {
  throw new Error("fetchMyTickets is not implemented yet.");
}
