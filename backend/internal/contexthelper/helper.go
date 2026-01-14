package contexthelper

import (
	"backend/config"
	"backend/pkg/logger"
	"context"
	"database/sql"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AccessTokenData struct {
	UserId       uint
	SetCookies   bool
	RefreshToken string
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

func SetServices(ctx context.Context, db *sql.DB, rabbitConn *amqp.Connection) context.Context {
	ctx = SetDb(ctx, db)
	return SetRabbitConn(ctx, rabbitConn)
}

func SetDb(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, dbCtxKey, db)
}

func SetRabbitConn(ctx context.Context, conn *amqp.Connection) context.Context {
	return context.WithValue(ctx, rabbitCtxKey, conn)
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
		log.Printf("get Rabbit Connection")
		return rabbitConn
	}
	logger.ErrorCtx(ctx, "There is no rabbit connection in context")
	return nil
}
