package cron

import (
	"context"
	"fmt"
	"testing"

	robfigCron "github.com/robfig/cron/v3"
)

func TestCronLoggerAdapterEnrichJobName(t *testing.T) {
	adapter := newCronLoggerAdapter(context.Background())
	adapter.RegisterJob(1, "heartbeat")
	adapter.RegisterJob(2, "retry_pending_queues")

	keysAndValues := []any{
		"now", "some-timestamp",
		"entry", robfigCron.EntryID(1),
		"next", "next-timestamp",
	}

	enriched := adapter.enrichKeysAndValuesWithJobName(keysAndValues)

	expectedEntry := fmt.Sprintf("%d (%s)", 1, "heartbeat")
	if enriched[3] != expectedEntry {
		t.Errorf("expected %q at index 3, got %q", expectedEntry, enriched[3])
	}

	unregisteredKeysAndValues := []any{
		"now", "some-timestamp",
		"entry", robfigCron.EntryID(99),
		"next", "next-timestamp",
	}

	notEnriched := adapter.enrichKeysAndValuesWithJobName(unregisteredKeysAndValues)
	if notEnriched[3] != robfigCron.EntryID(99) {
		t.Errorf("expected unchanged value for unregistered entry, got %v", notEnriched[3])
	}
}
