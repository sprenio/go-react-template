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
	"backend/internal/handler"
	"backend/internal/helper"
	"backend/internal/router"
	"backend/pkg/logger"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Fatal("Nie można załadować konfiguracji: %v", err)
	}

	logger.Init(cfg.LogLevel)
	logger.Info("Uruchamianie aplikacji: %s", cfg.AppName)

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

	// Inicjalizacja handlerów i routera
	h := handler.NewHandler(db, rabbitConn)
	webHostPort := fmt.Sprintf("%s:%v", cfg.WebServer.Host, cfg.WebServer.HTTPPort)
	srv := &http.Server{
		Addr:    webHostPort,
		Handler: router.SetupRouter(h),
	}
	
	// Obsługa sygnałów (graceful shutdown)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Start serwera: %s", webHostPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Błąd serwera: %v", err)
		}
	}()

	<-stop // czekaj na SIGINT/SIGTERM
	logger.Info("Zatrzymywanie aplikacji...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Błąd przy zamykaniu serwera: %v", err)
	}

	logger.Info("Aplikacja zatrzymana.")
}
