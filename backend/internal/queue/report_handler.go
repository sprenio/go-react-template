package queue

import (
	"encoding/json"
	"log"
	"context"
)

func (c *Consumer) HandleReportTask(ctx context.Context, task string, rawMessage json.RawMessage) error {
	log.Printf("📄 Generating report: %s", task)

	// Tu np. generowanie PDF, zapis do S3 itp.
	return nil
}
