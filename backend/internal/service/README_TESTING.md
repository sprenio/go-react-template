# Service Unit Tests

## Current Status

✅ **Unit tests are implemented and working!**

Tests use `sqlmock` (similar to handler tests) to mock database interactions, allowing us to test services with concrete repository types without needing to refactor to interfaces.

## Testing Approach

Instead of using mock interfaces, we use `sqlmock` to create a mock database connection. The repositories are instantiated with this mock database, and we set up SQL query expectations. This approach:

- ✅ Works with existing concrete repository types
- ✅ Tests actual SQL queries and interactions
- ✅ No need to refactor services to use interfaces
- ✅ Consistent with handler test approach

## Test Coverage

### ✅ AuthService (`auth_test.go`)
- Login success
- Invalid credentials
- User not found
- GetDataById error

### ✅ UserService (`user_test.go`)
- GetUserResponseData success
- GetUserResponseData user not found
- GetUserResponseData empty ID

### ✅ LanguageService (`language_test.go`)
- GetLanguages success
- GetLanguages empty list
- GetLanguages repository error

### ✅ UserSettingsService (`user_settings_test.go`)
- Update change language
- Update no update needed
- Update user not authenticated
- Update get by user ID error

## Running Tests

```bash
# Run all service tests
go test ./internal/service/... -v

# Run specific test
go test ./internal/service/... -v -run TestAuthService_Login_Success
```

## Test Structure

Tests follow this pattern:

1. Create a mock database using `sqlmock.New()`
2. Set up SQL query expectations using `mock.ExpectQuery()` or `mock.ExpectExec()`
3. Create repository instances with the mock database
4. Create service instances with the repositories
5. Call service methods
6. Verify results and check that all expectations were met

Example:
```go
func TestAuthService_Login_Success(t *testing.T) {
    ctx := context.Background()
    db, mock, err := sqlmock.New()
    // ... set up expectations ...
    
    userRepo := repository.NewUserRepository(db)
    authService := service.NewAuthService(userRepo)
    user, err := authService.Login(ctx, "test@example.com", "password123")
    
    // ... assertions ...
}
```

## Remaining Test Coverage

Some services don't have full test coverage yet (they require more complex setups or RabbitMQ mocking):

- **RegisterService**: Requires RabbitMQ connection for email queue
- **RegisterConfirmationService**: Requires proper token payload construction
- **PasswordService**: ResetPassword requires RabbitMQ connection
- **EmailChangeService**: Requires RabbitMQ connection and user context setup

These can be added as integration tests or with proper RabbitMQ mocking.
