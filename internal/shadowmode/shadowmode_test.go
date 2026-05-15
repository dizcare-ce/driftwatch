package shadowmode_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"driftwatch/internal/drift"
	"driftwatch/internal/shadowmode"
)

func clean() drift.Result  { return drift.Result{Service: "svc", Diffs: nil} }
func drifted() drift.Result { return drift.Result{Service: "svc", Diffs: []drift.Diff{{Field: "x"}}} }

func TestRun_StoresRecord(t *testing.T) {
	r := shadowmode.New(5)
	r.Run(context.Background(), func(_ context.Context) ([]drift.Result, error) {
		return []drift.Result{clean()}, nil
	})
	got := r.Last(10)
	if len(got) != 1 {
		t.Fatalf("expected 1 record, got %d", len(got))
	}
	if got[0].Err != nil {
		t.Fatalf("unexpected error: %v", got[0].Err)
	}
	if len(got[0].Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got[0].Results))
	}
}

func TestRun_RecordsError(t *testing.T) {
	r := shadowmode.New(5)
	want := errors.New("boom")
	r.Run(context.Background(), func(_ context.Context) ([]drift.Result, error) {
		return nil, want
	})
	got := r.Last(1)
	if !errors.Is(got[0].Err, want) {
		t.Fatalf("expected %v, got %v", want, got[0].Err)
	}
}

func TestRun_RecoversPanic(t *testing.T) {
	r := shadowmode.New(5)
	r.Run(context.Background(), func(_ context.Context) ([]drift.Result, error) {
		panic("unexpected panic")
	})
	got := r.Last(1)
	if got[0].Err == nil {
		t.Fatal("expected error from panic, got nil")
	}
}

func TestLast_RespectsCapacity(t *testing.T) {
	r := shadowmode.New(3)
	for i := 0; i < 6; i++ {
		r.Run(context.Background(), func(_ context.Context) ([]drift.Result, error) {
			return nil, nil
		})
	}
	got := r.Last(10)
	if len(got) != 3 {
		t.Fatalf("expected 3 records (cap), got %d", len(got))
	}
}

func TestLast_OrderedOldestFirst(t *testing.T) {
	r := shadowmode.New(5)
	times := make([]time.Time, 3)
	for i := range times {
		r.Run(context.Background(), func(_ context.Context) ([]drift.Result, error) {
			return nil, nil
		})
		times[i] = r.Last(10)[i].RunAt
		time.Sleep(time.Millisecond)
	}
	if !times[0].Before(times[1]) || !times[1].Before(times[2]) {
		t.Fatal("records not in ascending time order")
	}
}

func TestReset_ClearsRecords(t *testing.T) {
	r := shadowmode.New(5)
	r.Run(context.Background(), func(_ context.Context) ([]drift.Result, error) {
		return []drift.Result{drifted()}, nil
	})
	r.Reset()
	if got := r.Last(10); len(got) != 0 {
		t.Fatalf("expected 0 records after reset, got %d", len(got))
	}
}
