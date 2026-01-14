package handler_test

import (
	"backend/config"
	"backend/internal/contexthelper"
	"backend/internal/cookie"
	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	amqp "github.com/rabbitmq/amqp091-go"
)

type TestDeps struct {
	DB        *sql.DB
	Config    *config.Config
	UserID    uint
	AccessTokenData *contexthelper.AccessTokenData
	RequestID string
	RabbitConn *amqp.Connection
}

func NewTestRequest(
	method string,
	url string,
	body io.Reader,
	deps TestDeps,
) (*http.Request, *httptest.ResponseRecorder) {

	req := httptest.NewRequest(method, url, body)
	ctx := req.Context()

	if deps.DB != nil {
		ctx = contexthelper.SetDb(ctx, deps.DB)
	}

	var cfg *config.Config
	if deps.Config == nil {
		cfg = testConfig()
	} else {
		cfg = deps.Config
	}

	ctx = contexthelper.SetConfig(ctx, cfg)

	if deps.UserID != 0 {
		ctx = contexthelper.SetUserId(ctx, deps.UserID)
	}
	var reqId string
	if deps.RequestID == "" {
		n := rand.Intn(9000) + 1000 // 1000–9999
		reqId =  fmt.Sprintf("test-request-%d", n)
	} else {
		reqId = deps.RequestID
	}
	ctx = contexthelper.SetRequestID(ctx, reqId)

	if deps.AccessTokenData != nil {
		ctx = contexthelper.SetAccessTokenData(ctx, deps.AccessTokenData)
	}
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	return req, rr
}

func IsAccessCookieSet(resp *http.Response) bool {
	cookies := resp.Cookies()

	for _, c := range cookies {
		if c.Name == cookie.AccessTokenKey {
			return true
		}
	}
	return false
}

func testConfig() *config.Config {
	return &config.Config{
		AppEnv: "test",
	}
}
