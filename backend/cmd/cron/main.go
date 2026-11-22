package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"backend/config"
	"backend/internal/contexthelper"
	cronService "backend/internal/cron"
	"backend/internal/helper"
	"backend/pkg/logger"
)

const (
	CronServiceName = "cron"
	ShutdownTimeout = 15 * time.Second
)

func main() {
	appConfiguration, err := config.GetConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration: %v", err)
	}

	logger.Init(appConfiguration.LogLevel, contexthelper.GetRequestID)
	logger.Info("Starting application service: %s (%s)", appConfiguration.AppName, CronServiceName)

	connectionContext, cancelConnection := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelConnection()

	servicesHolder, err := helper.ConnectServicesWithRetry(
		connectionContext,
		appConfiguration.DB.DSN,
		appConfiguration.RabbitMQ.URL,
		appConfiguration.Redis.URL,
	)
	if err != nil {
		logger.Fatal("Failed to connect to required services: %v", err)
	}
	defer servicesHolder.RabbitConn.Close()
	defer servicesHolder.DB.Close()
	defer servicesHolder.RedisClient.Close()

	logger.Info("All services connected, starting cron scheduler...")

	applicationContext, stopApplication := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopApplication()

	applicationContext = contexthelper.SetServices(applicationContext, servicesHolder.DB, servicesHolder.RabbitConn, servicesHolder.RedisClient)
	applicationContext = contexthelper.SetConfig(applicationContext, appConfiguration)
	applicationContext = contexthelper.SetRequestID(applicationContext, CronServiceName)

	schedulerInstance := cronService.NewScheduler(applicationContext)

	if registrationError := schedulerInstance.RegisterScheduledJobs(); registrationError != nil {
		logger.Fatal("Failed to register scheduled cron jobs: %v", registrationError)
	}

	schedulerInstance.Start()
	logger.Info("🚀 Cron scheduler started successfully")

	<-applicationContext.Done()
	logger.Info("⏹ Shutting down cron scheduler...")

	shutdownContext, cancelShutdown := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancelShutdown()

	stopContext := schedulerInstance.Stop()

	select {
	case <-stopContext.Done():
		logger.Info("✅ All running cron jobs finished gracefully")
	case <-shutdownContext.Done():
		logger.Warn("⚠️ Shutdown timeout reached, forcing exit of cron scheduler")
	}

	logger.Info("Cron service stopped.")
}

