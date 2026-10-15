package appeals

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestIsLockExpiredBoundary(t *testing.T) {
	now := time.Date(2026, 10, 12, 10, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		lockedAt time.Time
		want     bool
	}{
		{"just taken", now, false},
		{"29m59s old", now.Add(-LockDuration + time.Second), false},
		{"exactly 30 minutes old", now.Add(-LockDuration), true},
		{"31 minutes old", now.Add(-31 * time.Minute), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsLockExpired(tc.lockedAt, now); got != tc.want {
				t.Fatalf("IsLockExpired = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestLockExpiryCutoffMatchesIsLockExpired(t *testing.T) {
	now := time.Date(2026, 10, 12, 10, 0, 0, 0, time.UTC)
	cutoff := LockExpiryCutoff(now)
	if !IsLockExpired(cutoff, now) {
		t.Fatal("a lock taken exactly at the cutoff must be expired")
	}
	if IsLockExpired(cutoff.Add(time.Nanosecond), now) {
		t.Fatal("a lock taken after the cutoff must not be expired")
	}
}

type sweepCounter struct {
	calls atomic.Int32
}

func (s *sweepCounter) Get(context.Context, string, string) (Ticket, error) {
	return Ticket{}, ErrNotImplemented
}
func (s *sweepCounter) Release(context.Context, string, string) error { return ErrNotImplemented }
func (s *sweepCounter) ExtendLock(context.Context, string, string) (time.Time, error) {
	return time.Time{}, ErrNotImplemented
}
func (s *sweepCounter) Resolve(context.Context, string, string, string, string) error {
	return ErrNotImplemented
}
func (s *sweepCounter) ReleaseExpired(context.Context, time.Time) (int64, error) {
	s.calls.Add(1)
	return 0, ErrNotImplemented
}

func TestRunExpiryWorkerSweepsUntilCancelled(t *testing.T) {
	store := &sweepCounter{}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		RunExpiryWorker(ctx, store, 5*time.Millisecond)
		close(done)
	}()

	deadline := time.After(2 * time.Second)
	for store.calls.Load() < 2 {
		select {
		case <-deadline:
			t.Fatal("worker did not sweep repeatedly")
		case <-time.After(5 * time.Millisecond):
		}
	}

	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not stop after cancellation")
	}
}

type purgeCounter struct {
	cutoffs chan time.Time
}

func (p *purgeCounter) PurgeDeletedReviews(_ context.Context, before time.Time) (int64, error) {
	p.cutoffs <- before
	return 0, nil
}

func TestRunDeletedReviewPurgeWorkerUsesRetentionWindow(t *testing.T) {
	store := &purgeCounter{cutoffs: make(chan time.Time, 1)}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		RunDeletedReviewPurgeWorker(ctx, store, 5*time.Millisecond)
		close(done)
	}()

	var cutoff time.Time
	select {
	case cutoff = <-store.cutoffs:
	case <-time.After(2 * time.Second):
		t.Fatal("purge worker did not run")
	}
	now := time.Now().UTC()
	if cutoff.After(now.Add(-DeletedReviewRetention)) || cutoff.Before(now.Add(-DeletedReviewRetention-time.Second)) {
		t.Fatalf("purge cutoff = %s, want approximately %s", cutoff, now.Add(-DeletedReviewRetention))
	}

	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("purge worker did not stop after cancellation")
	}
}
