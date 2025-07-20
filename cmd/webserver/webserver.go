package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"payroll/internal/infrastructure/config"
	"payroll/internal/infrastructure/messagebus"

	"github.com/twmb/franz-go/pkg/kgo"
)

type WebServer struct {
	cfg         *config.Config
	srv         *http.Server
	kafkaClient *kgo.Client
}

func (w *WebServer) Start() {
	// Start Kafka consumer
	go func() {
		log.Println("Starting Kafka consumer...")
		messagebus.StartConsumer(context.Background(), w.kafkaClient)
	}()

	serverErrors := make(chan error, 1)
	// Start the server in a separate goroutine so that it doesn't block.
	go func() {
		log.Printf("Starting server on %s", w.srv.Addr)
		serverErrors <- w.srv.ListenAndServe()
	}()

	// Set up a channel to listen for OS signals for graceful shutdown.
	quit := make(chan os.Signal, 1)
	// Listen for interrupt (Ctrl+C) or termination signals. SIGKILL cannot be caught.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a shutdown signal or the server fails.
	select {
	case err := <-serverErrors:
		// ListenAndServe will return ErrServerClosed on a graceful shutdown.
		// We only want to fatal on other unexpected errors.
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	case sig := <-quit:
		log.Printf("Received signal: %v. Initiating graceful shutdown...", sig)
		w.Stop()
	}
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
