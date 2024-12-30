package main

import (
	"context"
	"errors"
	"github.com/Tus1688/library-management-api/api"
	"github.com/Tus1688/library-management-api/authutil"
	"github.com/Tus1688/library-management-api/cache"
	"github.com/Tus1688/library-management-api/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// main is the entry point of the application.
// It initializes the Postgres, Redis, and session stores, and starts the HTTP server.
// It also handles graceful shutdown on receiving termination signals.
func main() {
	// Initialize Postgres store
	postgres, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal("Unable to connect to database")
	}

	// Initialize Redis store
	redis, err := cache.NewRedisStore(1)
	if err != nil {
		log.Fatal("Unable to connect to redis")
	}

	// Initialize session store
	session, err := authutil.NewSessionStore()
	if err != nil {
		log.Fatal("Unable to create session store")
	}

	// Create a new server
	server := api.NewServer(":8080", postgres, redis, session)
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Set up signal handling for graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds to allow server to clean up
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("gracefully shutdown timed out.. forcing shutdown")
			}
		}()

		// Shutdown server
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal("error shutting down server", err)
		}
		serverStopCtx()
	}()

	// Start server
	err = server.Run()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("error running server", err)
	}

	// Wait for server to stop
	<-serverCtx.Done()
}
