package contexthelper

type contextKey string

const (
	RequestIDKey contextKey = "requestID"
	UserIdKey    contextKey = "userID"
	// w przyszłości np. UserIDKey, SessionIDKey itd.
)
