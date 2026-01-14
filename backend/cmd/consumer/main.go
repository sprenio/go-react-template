package main

import (
	"backend/config"
	"backend/internal/contexthelper"
	"backend/internal/helper"
	"backend/internal/queue"
	"backend/pkg/logger"
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Fatal("Nie można załadować konfiguracji: %v", err)
	}

	logger.Init(cfg.LogLevel, contexthelper.GetRequestID)

	// Kontekst z timeoutem na połączenie z DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	db, rabbitConn, err := helper.ConnectServicesWithRetry(ctx, cfg.DB.DSN, cfg.RabbitMQ.URL)
	if err != nil {
		logger.Fatal("Nie udało się połączyć z usługami: %v", err)
	}
	defer rabbitConn.Close()
	defer db.Close()

	logger.Info("Wszystkie usługi gotowe, start backendu...")

	// Główny kontekst aplikacji – będzie anulowany na SIGINT/SIGTERM
	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Info("🚀 Starting consumers...")
	c := queue.NewConsumer(db, rabbitConn)

	// Warto dodać, żeby consumer dostał kontekst (zatrzyma się na cancel)
	go func() {
		if err := c.StartEmailConsumer(appCtx); err != nil {
			log.Printf("Email consumer stopped with error: %v", err)
		}
	}()
	go func() {
		if err := c.StartReportConsumer(appCtx); err != nil {
			log.Printf("Report consumer stopped with error: %v", err)
		}
	}()

	// Czekaj na sygnał zakończenia
	<-appCtx.Done()
	logger.Info("⏹ Zatrzymywanie consumerów...")

	// jeśli masz c.Stop() → można tu wywołać
	// c.Stop()

	logger.Info("✅ Consumers zatrzymani")
}
