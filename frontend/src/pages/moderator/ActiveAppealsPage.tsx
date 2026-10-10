import type { AuthSession } from "../../services/auth";

// Person A owns this file (page 3, SCRUM-10 and SCRUM-11). Keep the props.
//
// Show the appeals this moderator has claimed (fetchMyTickets from
// services/appealsQueue.ts), using the same two-line row as the queue.
// Clicking a row calls onOpenTicket(id); it does NOT claim again.

export interface ActiveAppealsPageProps {
  session: AuthSession;
  onOpenTicket: (ticketId: string) => void;
}

function ActiveAppealsPage(_props: ActiveAppealsPageProps) {
  return (
    <div>
      <h1 className="text-2xl font-bold">Active appeals</h1>
      <p className="mt-3 text-ink/65">
        The active appeals list has not been built yet.
      </p>
    </div>
  );
}

export default ActiveAppealsPage;
