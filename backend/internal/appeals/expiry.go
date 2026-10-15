package appeals

import (
	"context"
	"log"
	"time"
)

// DefaultSweepInterval is how often RunExpiryWorker looks for expired locks.
const DefaultSweepInterval = time.Minute

// DeletedReviewRetention is how long a deleted review and its appeal details
// remain available to the student by reference code.
const DeletedReviewRetention = 72 * time.Hour

// DefaultDeletedReviewPurgeInterval is how often the purge worker removes
// reviews whose retention period has elapsed.
const DefaultDeletedReviewPurgeInterval = time.Minute

type DeletedReviewPurger interface {
	PurgeDeletedReviews(context.Context, time.Time) (int64, error)
}

// RunExpiryWorker releases tickets whose lock is older than LockDuration until
// ctx is cancelled. It is started once from cmd/api. Queries must also treat
// expired locks as claimable, so a slow sweep never leaves a ticket stuck.
func RunExpiryWorker(ctx context.Context, store ResolveStore, interval time.Duration) {
	if interval <= 0 {
		interval = DefaultSweepInterval
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			released, err := store.ReleaseExpired(ctx, now.UTC())
			switch {
			case err == ErrNotImplemented:
				// Scaffold stub: nothing to sweep until ReleaseExpired is written.
			case err != nil:
				log.Printf("release expired appeal locks: %v", err)
			case released > 0:
				log.Printf("released %d expired appeal lock(s)", released)
			}
		}
	}
}

// RunDeletedReviewPurgeWorker removes deleted reviews and their appeal tickets
// after the student has had DeletedReviewRetention to check the outcome.
func RunDeletedReviewPurgeWorker(ctx context.Context, store DeletedReviewPurger, interval time.Duration) {
	if interval <= 0 {
		interval = DefaultDeletedReviewPurgeInterval
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			deleted, err := store.PurgeDeletedReviews(ctx, now.UTC().Add(-DeletedReviewRetention))
			if err != nil {
				log.Printf("purge expired deleted reviews: %v", err)
			} else if deleted > 0 {
				log.Printf("permanently deleted %d expired review(s) and appeal ticket(s)", deleted)
			}
		}
	}
}
