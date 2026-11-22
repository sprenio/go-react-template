package main

import (
	"backend/config"
	"backend/internal/contexthelper"
	"backend/internal/helper"
	queueConsumer "backend/internal/queue/consumer"
	"backend/pkg/logger"
	"context"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration: %v", err)
	}

	logger.Init(cfg.LogLevel, contexthelper.GetRequestID)

	// Context with timeout for DB connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	svcs, err := helper.ConnectServicesWithRetry(
		ctx,
		cfg.DB.DSN,
		cfg.RabbitMQ.URL,
		cfg.Redis.URL,
	)
	if err != nil {
		logger.Fatal("Failed to connect to services: %v", err)
	}
	defer svcs.RabbitConn.Close()
	defer svcs.DB.Close()
	defer svcs.RedisClient.Close()
	logger.Info("All services ready, starting backend...")

	// Main application context – will be cancelled on SIGINT/SIGTERM
	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	appCtx = contexthelper.SetServices(appCtx, svcs.DB, svcs.RabbitConn, svcs.RedisClient)
	appCtx = contexthelper.SetConfig(appCtx, cfg)
	appCtx = contexthelper.SetRequestID(appCtx, "consumer")

	logger.InfoCtx(appCtx, "🚀 Starting consumers...")

	c := queueConsumer.NewConsumer(appCtx)

	var wg sync.WaitGroup

	startConsumer := func(name string, run func(context.Context) error) {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := run(appCtx); err != nil {
				logger.ErrorCtx(appCtx, "Consumer stopped with error", "consumer", name, "error", err)

				// Error in one consumer causes the others to stop.
				stop()
			}

			logger.InfoCtx(appCtx, "Consumer stopped", "consumer", name)
		}()
	}

	startConsumer("email", c.StartEmailConsumer)
	startConsumer("register", c.StartRegisterConsumer)
	startConsumer("reset-password", c.StartResetPasswordConsumer)

	// Wait for SIGTERM, SIGINT, or an error from one of the consumers.
	<-appCtx.Done()

	logger.InfoCtx(appCtx, "⏹ Stopping consumers...")

	// Wait for all goroutines to finish.
	wg.Wait()

	logger.InfoCtx(appCtx, "✅ All consumers stopped")
}
