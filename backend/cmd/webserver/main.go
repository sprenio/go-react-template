package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/config"
	"backend/internal/contexthelper"
	"backend/internal/helper"
	"backend/internal/router"
	"backend/pkg/logger"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration: %v", err)
	}

	logger.Init(cfg.LogLevel, contexthelper.GetRequestID)
	logger.Info("Starting application", "app_name", cfg.AppName)

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
		logger.Fatal("Failed to connect to services", "error", err)
	}
	defer svcs.RabbitConn.Close()
	defer svcs.DB.Close()
	defer svcs.RedisClient.Close()

	logger.Info("All services ready, starting backend...")

	// Initialize handlers and router
	webHostPort := fmt.Sprintf("%s:%v", cfg.WebServer.Host, cfg.WebServer.HTTPPort)
	srv := &http.Server{
		Addr:    webHostPort,
		Handler: router.SetupRouter(cfg, svcs.DB, svcs.RabbitConn, svcs.RedisClient),
	}

	// Signal handling (graceful shutdown)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Starting server", "port", webHostPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error: %v", err)
		}
	}()

	<-stop // wait for SIGINT/SIGTERM
	logger.Info("Stopping application...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error shutting down server", "error", err)
	}

	logger.Info("Aplikacja zatrzymana.")
}
