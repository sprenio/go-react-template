package cron

import (
	"context"
	"fmt"
	"sync"

	"backend/pkg/apperror"
	"backend/pkg/logger"

	robfigCron "github.com/robfig/cron/v3"
)

type cronLoggerAdapter struct {
	appContext                  context.Context
	registeredJobNamesByEntryID map[int]string
	jobNamesMutex               sync.RWMutex
}

func newCronLoggerAdapter(appContext context.Context) *cronLoggerAdapter {
	return &cronLoggerAdapter{
		appContext:                  appContext,
		registeredJobNamesByEntryID: make(map[int]string),
	}
}

func (adapter *cronLoggerAdapter) RegisterJob(entryID int, jobName string) {
	adapter.jobNamesMutex.Lock()
	defer adapter.jobNamesMutex.Unlock()
	adapter.registeredJobNamesByEntryID[entryID] = jobName
}

func (adapter *cronLoggerAdapter) getJobName(entryID int) (string, bool) {
	adapter.jobNamesMutex.RLock()
	defer adapter.jobNamesMutex.RUnlock()
	jobName, exists := adapter.registeredJobNamesByEntryID[entryID]
	return jobName, exists
}

func (adapter *cronLoggerAdapter) Info(message string, keysAndValues ...any) {
	enrichedKeysAndValues := adapter.enrichKeysAndValuesWithJobName(keysAndValues)
	msg := fmt.Sprintf("[CRON] %s", message)
	logger.InfoCtx(adapter.appContext, msg, enrichedKeysAndValues...)
}

func (adapter *cronLoggerAdapter) Error(err error, message string, keysAndValues ...any) {
	enrichedKeysAndValues := adapter.enrichKeysAndValuesWithJobName(keysAndValues)
	msg := fmt.Sprintf("[CRON ERROR] %s", message)
	logErr := apperror.Wrap(err, msg, enrichedKeysAndValues...)
	logger.ErrorCtx(adapter.appContext, logErr.Error())
}

func (adapter *cronLoggerAdapter) enrichKeysAndValuesWithJobName(keysAndValues []any) []any {
	if len(keysAndValues) == 0 {
		return keysAndValues
	}

	enrichedList := make([]any, len(keysAndValues))
	copy(enrichedList, keysAndValues)

	for index := 0; index < len(enrichedList)-1; index += 2 {
		keyString, isString := enrichedList[index].(string)
		if isString && keyString == "entry" {
			entryID, isParsed := extractEntryID(enrichedList[index+1])
			if isParsed {
				if jobName, exists := adapter.getJobName(entryID); exists {
					enrichedList[index+1] = fmt.Sprintf("%d (%s)", entryID, jobName)
				}
			}
		}
	}

	return enrichedList
}

func extractEntryID(value any) (int, bool) {
	switch typedValue := value.(type) {
	case int:
		return typedValue, true
	case int32:
		return int(typedValue), true
	case int64:
		return int(typedValue), true
	case robfigCron.EntryID:
		return int(typedValue), true
	default:
		return 0, false
	}
}
