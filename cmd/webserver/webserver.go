package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"payroll/internal/app"
	"payroll/internal/infrastructure/config"
	"payroll/internal/infrastructure/messagebus"

	"github.com/twmb/franz-go/pkg/kgo"
)

type WebServer struct {
	cfg         *config.Config
	srv         *http.Server
	services    *app.Services
	kafkaClient *kgo.Client
}

func (w *WebServer) Start() {
	// Create a context that is canceled when a shutdown signal is received.
	// This is the idiomatic way to handle graceful shutdowns.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop() // Important to release resources used by NotifyContext.

	// Start Kafka consumer with the cancellable context.
	go func() {
		log.Println("Starting Kafka consumer...")
		messagebus.StartConsumer(ctx, w.kafkaClient, w.cfg.Kafka, w.services)
	}()

	serverErrors := make(chan error, 1)
	// Start the server in a separate goroutine so that it doesn't block.
	go func() {
		log.Printf("Starting server on %s", w.srv.Addr)
		serverErrors <- w.srv.ListenAndServe()
	}()

	// Block until the context is canceled (due to a signal) or the server fails.
	select {
	case err := <-serverErrors:
		// ListenAndServe will return ErrServerClosed on a graceful shutdown.
		// We only want to fatal on other unexpected errors.
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	case <-ctx.Done():
		log.Printf("shutdown signal received: %v", ctx.Err())
	}

	w.Stop()
	log.Println("Server has shut down.")
}

func (w *WebServer) Stop() {
	// Give the server a deadline to finish processing existing requests.
	// It's good practice to have a default timeout in case it's not in the config.
	timeout := 15 * time.Second
	if w.cfg.Server.ShutdownTimeout > 0 {
		timeout = w.cfg.Server.ShutdownTimeout
	}
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()

	// Close the Kafka client.
	log.Println("Closing Kafka client...")
	w.kafkaClient.Close()

	// srv.Shutdown() gracefully shuts down the server.
	if err := w.srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown gracefully: %v", err)
	} else {
		log.Println("Server gracefully stopped.")
	}
}
