import type { AuthSession } from "../../services/auth";

// Person B owns this file (page 2, SCRUM-12 and SCRUM-13). Keep the props.
//
// Load the ticket with fetchTicket from services/appealsResolve.ts and show:
// course id, year and professor; the reference code; the feedback; the reason
// the LLM flagged it; the student's appeal comment. Bottom right, three
// buttons: send back to queue (releaseTicket), publish (restoreTicket) and
// delete (opens a DeleteModal with Cancel and Submit; the comment is required).
// After any of the three succeeds, call onDone with a short success message,
// for example onDone("Review published"). The shell shows it on page 1.
//
// If fetchTicket fails with a 403 or 404 ApiError the lock has expired or the
// ticket was taken: call onDone("This appeal is no longer assigned to you").
// Put DeleteModal in its own file, src/pages/moderator/DeleteModal.tsx.

export interface TicketDetailPageProps {
  session: AuthSession;
  ticketId: string;
  onDone: (message: string) => void;
}

function TicketDetailPage({ ticketId }: TicketDetailPageProps) {
  return (
    <div>
      <h1 className="text-2xl font-bold">Appeal details</h1>
      <p className="mt-3 text-ink/65">
        The appeal details page has not been built yet (ticket {ticketId}).
      </p>
    </div>
  );
}

export default TicketDetailPage;
