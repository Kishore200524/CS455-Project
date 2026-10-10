package database

// Person B owns this file: detail, release, resolve and the expiry sweep
// (SCRUM-12, SCRUM-13). Replace each stub body; keep the signatures.
//
// Rules from docs/moderator-api.md that these operations must follow:
//   - Get, Release and Resolve only work for the moderator in lockedBy, and
//     only while the lock is unexpired (see appeals.IsLockExpired). Do the
//     lockedBy/status check inside the update filter, not in Go after reading.
//   - Release sets status back to under_appeal and clears lockedBy/lockedAt.
//   - Resolve sets status to restored or deleted, stores resolutionComment and
//     resolvedAt, and updates the review in the feedback collection
//     (restored -> published, deleted -> deleted).
//   - ReleaseExpired is one UpdateMany: status locked and
//     lockedAt <= appeals.LockExpiryCutoff(now) -> under_appeal.
//   - Never read or write any student identity here.

import (
	"context"
	"time"

	"github.com/Kishore200524/CS455-Project/backend/internal/appeals"
)

var _ appeals.ResolveStore = (*Mongo)(nil)

func (m *Mongo) Get(_ context.Context, _, _ string) (appeals.Ticket, error) {
	return appeals.Ticket{}, appeals.ErrNotImplemented
}

func (m *Mongo) Release(_ context.Context, _, _ string) error {
	return appeals.ErrNotImplemented
}

func (m *Mongo) Resolve(_ context.Context, _, _, _, _ string) error {
	return appeals.ErrNotImplemented
}

func (m *Mongo) ReleaseExpired(_ context.Context, _ time.Time) (int64, error) {
	return 0, appeals.ErrNotImplemented
}
