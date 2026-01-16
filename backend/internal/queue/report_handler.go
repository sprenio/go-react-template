package queue

import (
	"backend/pkg/logger"
	"context"
	"encoding/json"
)

func (c *Consumer) HandleReportTask(ctx context.Context, task string, rawMessage json.RawMessage) error {
	logger.InfoCtx(ctx, "📄 Generating report: %s", task)

	// Tu np. generowanie PDF, zapis do S3 itp.
	return nil
}
