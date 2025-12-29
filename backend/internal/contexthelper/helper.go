package contexthelper

import (
	"context"
)

// GetUserId pozwala wyciągnąć user_id z kontekstu w handlerach
func GetUserId(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIdKey).(uint)
	return userID, ok
}

func SetUserId(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, UserIdKey, userID)
}


// Funkcja pomocnicza do pobrania request ID z kontekstu
func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}


