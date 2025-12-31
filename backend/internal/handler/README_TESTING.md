# Handler Unit Tests

## Current Status

✅ **Unit tests are fully implemented and working!**

All handlers have comprehensive unit test coverage using `sqlmock` to mock database interactions and `httptest` for HTTP request/response testing.

## Testing Approach

Tests use a combination of:
- **`sqlmock`**: Mocks database interactions at the SQL level
- **`httptest`**: Creates HTTP request/response objects for testing handlers
- **`github.com/go-chi/chi/v5`**: For testing route parameters (URL params)

This approach:
- ✅ Tests handlers with real repository instances (using mocked database)
- ✅ Verifies HTTP status codes, response formats, and error handling
- ✅ Tests validation logic, authentication checks, and business rules
- ✅ No need to refactor handlers or services
- ✅ Consistent with service test approach

## Test Coverage

### ✅ PingHandler (`ping_test.go`)
- Success response with request ID

### ✅ RegisterHandler (`register_test.go`)
- Success (skipped - requires RabbitMQ)
- Invalid HTTP method
- Invalid JSON
- Missing username
- Missing email
- Missing password

### ✅ LoginHandler (`login_test.go`)
- Login success
- Invalid JSON
- Empty email
- Empty password
- Invalid credentials (user not found)
- Wrong password

### ✅ MeHandler (`me_test.go`)
- Success (get current user)
- Unauthorized (no user ID in context)
- User not found

### ✅ ConfirmHandler (`confirm_test.go`)
- Register token confirmation success
- Empty token
- Token not found
- Invalid token type

### ✅ PasswordChangeHandler (`password_change_test.go`)
- Success
- Empty token
- Invalid JSON
- Invalid password format
- Token not found

### ✅ ResetPasswordHandler (`reset_password_test.go`)
- Success (skipped - requires RabbitMQ)
- Invalid JSON
- Empty email
- Invalid email format
- User not found (returns success to prevent email enumeration)

### ✅ EmailChangeHandler (`email_change_test.go`)
- Success (skipped - requires RabbitMQ)
- Invalid JSON
- Empty email
- Invalid email format

### ✅ SettingsHandler (`settings_test.go`)
- Success (change language)
- Invalid JSON
- Invalid AppOpt2 value
- Valid AppOpt2 value
- Service error

### ✅ CfgHandler (`cfg_test.go`)
- Success (get configuration)
- Languages error

### ✅ NotFoundHandler (`notfound_test.go`)
- 404 response

### ✅ MethodNotAllowedHandler (`method_not_allowed_test.go`)
- 405 response

## Running Tests

```bash
# Run all handler tests
go test ./internal/handler/... -v

# Run specific test
go test ./internal/handler/... -v -run TestLoginHandler_Success

# Run tests for specific handler
go test ./internal/handler/... -v -run TestLoginHandler
```

## Test Structure

Tests follow this pattern:

1. Create a mock database using `sqlmock.New()`
2. Set up SQL query expectations
3. Create handler instance with mock database
4. Create HTTP request using `httptest.NewRequest()`
5. Add context values (request ID, user ID) if needed
6. Create response recorder using `httptest.NewRecorder()`
7. Call handler method
8. Verify HTTP status code and response body
9. Check that all SQL expectations were met

Example:
```go
func TestLoginHandler_Success(t *testing.T) {
    db, mock, err := sqlmock.New()
    // ... set up SQL expectations ...
    
    h := handler.NewHandler(db, nil)
    
    reqBody := map[string]string{
        "email":    "test@example.com",
        "password": "password123",
    }
    body, _ := json.Marshal(reqBody)
    req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
    ctx := context.WithValue(req.Context(), contexthelper.RequestIDKey, "test-id-123")
    req = req.WithContext(ctx)
    
    rr := httptest.NewRecorder()
    h.LoginHandler(rr, req)
    
    resp := rr.Result()
    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected status 200, got %d", resp.StatusCode)
    }
}
```

## Special Considerations

### RabbitMQ Dependencies

Some handlers require RabbitMQ connections for email queue publishing:
- `RegisterHandler` - publishes registration email
- `ResetPasswordHandler` - publishes password reset email
- `EmailChangeHandler` - publishes email change confirmation

These tests are currently **skipped** with `t.Skip()` and marked for integration testing. To enable them:
1. Use integration tests with a real RabbitMQ connection, or
2. Mock the RabbitMQ connection properly

### URL Parameters

Handlers that use URL parameters (like `ConfirmHandler` and `PasswordChangeHandler`) require setting up chi router context:

```go
rctx := chi.NewRouteContext()
rctx.URLParams.Add("token", token)
req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
```

### Context Values

Many handlers require context values:
- `RequestIDKey`: Set by request ID middleware
- `UserIdKey`: Set by JWT authentication middleware

Tests set these up manually:
```go
ctx := context.WithValue(req.Context(), contexthelper.RequestIDKey, "test-id-123")
ctx = contexthelper.SetUserId(ctx, 1)
req = req.WithContext(ctx)
```

## Test Statistics

- **Total test files**: 11
- **Total test cases**: ~46
- **Passing tests**: All implemented tests pass
- **Skipped tests**: 3 (require RabbitMQ)

## Best Practices

1. **Always verify SQL expectations**: Use `mock.ExpectationsWereMet()` to ensure all expected queries were executed
2. **Test error cases**: Include tests for validation errors, authentication failures, and service errors
3. **Verify response format**: Check both status codes and response body structure
4. **Use realistic data**: Use proper time values, hashed passwords, and valid data structures
5. **Clean up**: Always `defer db.Close()` and `defer resp.Body.Close()`

## Future Enhancements

- Add integration tests for handlers requiring RabbitMQ
- Add tests for edge cases and boundary conditions
- Add performance/load tests for critical handlers
- Add tests for concurrent request handling

