package database

// Person A owns this file: queue, claim and active-appeals queries
// (SCRUM-10, SCRUM-11). Replace each stub body; keep the signatures.
//
// Rules from docs/moderator-api.md that these queries must follow:
//   - A ticket is claimable when status is under_appeal, OR status is locked
//     and lockedAt <= appeals.LockExpiryCutoff(now). Use the same filter for
//     ListUnclaimed and for Claim, so an expired lock can be claimed even if
//     the sweep has not run yet.
//   - Claim must be ONE atomic update (FindOneAndUpdate with the claimable
//     filter in the query), never a read followed by a write. No match means
//     ErrAlreadyClaimed if the ticket exists, ErrNotFound if it does not.
//   - Never read or write any student identity here.

import (
	"context"

	"github.com/Kishore200524/CS455-Project/backend/internal/appeals"
)

var _ appeals.QueueStore = (*Mongo)(nil)

func (m *Mongo) ListUnclaimed(_ context.Context, _ string) ([]appeals.Summary, error) {
	return nil, appeals.ErrNotImplemented
}

func (m *Mongo) Claim(_ context.Context, _, _ string) (appeals.Summary, error) {
	return appeals.Summary{}, appeals.ErrNotImplemented
}

func (m *Mongo) ListMine(_ context.Context, _ string) ([]appeals.Summary, error) {
	return nil, appeals.ErrNotImplemented
}
