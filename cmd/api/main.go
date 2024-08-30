package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func main() {
	cfg := config.New()

	router := chi.NewMux()

	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	}))

	killSignal := make(chan os.Signal, 1)

	signal.Notify(killSignal, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		err := server.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			cfg.Logger.Info("Server shutdown complete")
			// log.Printf("Server shutdown complete")
		} else if err != nil {
			cfg.Logger.Fatal("server failed to start", zap.Error(err))
			// log.Printf("Failed to start server")
			// os.Exit(1)
		}
	}()

	cfg.Logger.Info("Server started...")
	// log.Printf("Server started")

	// Wait for killsignal
	<-killSignal

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		cfg.Logger.Fatal("server shutdown failed", zap.Error(err))
		// log.Printf("Server shutdown failed")
		// os.Exit(1)
	}
}
