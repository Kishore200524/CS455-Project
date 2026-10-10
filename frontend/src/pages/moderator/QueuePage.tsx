import type { AuthSession } from "../../services/auth";

// Person A owns this file (page 1, SCRUM-10 and SCRUM-11). Keep the props.
//
// Show unclaimed appeals with a "latest to oldest" / "oldest to latest" toggle
// (default oldest first). Each row: course id, year and professor on one line,
// the reference code on the next. Clicking a row calls claimTicket from
// services/appealsQueue.ts and then onTicketClaimed(id). On a 409 ApiError show
// "This ticket has already been claimed" and refresh the list.

export interface QueuePageProps {
  session: AuthSession;
  onTicketClaimed: (ticketId: string) => void;
}

function QueuePage(_props: QueuePageProps) {
  return (
    <div>
      <h1 className="text-2xl font-bold">Appeal queue</h1>
      <p className="mt-3 text-ink/65">The appeal queue has not been built yet.</p>
    </div>
  );
}

export default QueuePage;
