package contexthelper

import (
	"backend/config"
	"backend/pkg/logger"
	"context"
	"database/sql"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type AccessTokenData struct {
	UserId       uint
	SetCookies   bool
	RefreshToken string
}

type CurrencyOrderLoaderFunc func(ctx context.Context, userId uint) ([]string, error)
type CurrencyOrderHolder struct {
	loader CurrencyOrderLoaderFunc
	order  []string
	loaded bool
	mu     sync.Mutex
}

func InitCurrencyOrder(ctx context.Context, loader CurrencyOrderLoaderFunc) context.Context {
	return context.WithValue(ctx, currencyOrderCtxKey, &CurrencyOrderHolder{
		loader: loader,
	})
}

func GetAccessTokenData(ctx context.Context) (*AccessTokenData, context.Context) {
	data, ok := ctx.Value(accessTokenDataCtxKey).(*AccessTokenData)
	if !ok {
		data = &AccessTokenData{}
		ctx = SetAccessTokenData(ctx, data)
	}
	return data, ctx
}

func SetAccessTokenData(ctx context.Context, data *AccessTokenData) context.Context {
	return context.WithValue(ctx, accessTokenDataCtxKey, data)
}

func GetUserId(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(userIdCtxKey).(uint)
	return userID, ok && userID != 0
}

func SetUserId(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIdCtxKey, userID)
}

func SetClientIp(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, clientIpCtxKey, ip)
}

// Funkcja pomocnicza do pobrania request ID z kontekstu
func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(requestIdCtxKey).(string); ok {
		return reqID
	}
	return ""
}

func GetClientIp(ctx context.Context) string {
	if ip, ok := ctx.Value(clientIpCtxKey).(string); ok {
		return ip
	}
	return ""
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIdCtxKey, requestID)
}

func SetConfig(ctx context.Context, cfg *config.Config) context.Context {
	return context.WithValue(ctx, configCtxKey, cfg)
}

func GetConfig(ctx context.Context) *config.Config {
	if cfg, ok := ctx.Value(configCtxKey).(*config.Config); ok {
		return cfg
	}
	logger.ErrorCtx(ctx, "There is no config in context")
	return nil
}

func SetServices(ctx context.Context, db *sql.DB, rabbitConn *amqp.Connection, redisClient *redis.Client) context.Context {
	ctx = SetDb(ctx, db)
	ctx = SetRabbitConn(ctx, rabbitConn)
	return SetRedisClient(ctx, redisClient)
}

func SetDb(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, dbCtxKey, db)
}

func SetRabbitConn(ctx context.Context, conn *amqp.Connection) context.Context {
	return context.WithValue(ctx, rabbitCtxKey, conn)
}

func SetRedisClient(ctx context.Context, redisClient *redis.Client) context.Context {
	return context.WithValue(ctx, redisClientCtxKey, redisClient)
}

func SetRateLimiter(ctx context.Context, rateLimiter any) context.Context {
	return context.WithValue(ctx, rateLimiterCtxKey, rateLimiter)
}

func GetRateLimiter(ctx context.Context) any {
	if rateLimiter := ctx.Value(rateLimiterCtxKey); rateLimiter != nil {
		return rateLimiter
	}
	logger.ErrorCtx(ctx, "There is no rate limiter in context")
	return nil
}
func GetDb(ctx context.Context) *sql.DB {
	if db, ok := ctx.Value(dbCtxKey).(*sql.DB); ok {
		return db
	}
	logger.ErrorCtx(ctx, "There is no db in context")
	return nil
}
func GetRabbitConn(ctx context.Context) *amqp.Connection {
	if rabbitConn, ok := ctx.Value(rabbitCtxKey).(*amqp.Connection); ok {
		logger.DebugCtx(ctx, "get Rabbit Connection")
		return rabbitConn
	}
	logger.ErrorCtx(ctx, "There is no rabbit connection in context")
	return nil
}

func GetRedisClient(ctx context.Context) *redis.Client {
	if redisClient, ok := ctx.Value(redisClientCtxKey).(*redis.Client); ok {
		logger.DebugCtx(ctx, "get Redis Client")
		return redisClient
	}
	logger.ErrorCtx(ctx, "There is no redis client in context")
	return nil
}

func GetCurrencyOrder(ctx context.Context) []string {
	holder, ok := ctx.Value(currencyOrderCtxKey).(*CurrencyOrderHolder)
	if !ok || holder == nil {
		return nil
	}
	holder.mu.Lock()
	defer holder.mu.Unlock()
	if holder.loaded {
		return holder.order
	}
	userId, ok := GetUserId(ctx)
	if !ok || holder.loader == nil {
		return nil
	}
	// Wywołanie przekazanego loadera przy pierwszym żądaniu
	orderedCurrencyCodes, err := holder.loader(ctx, userId)
	if err == nil {
		holder.order = orderedCurrencyCodes
		holder.loaded = true
	}
	return holder.order
}
