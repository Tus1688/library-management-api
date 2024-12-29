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

func main() {
	postgres, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal("Unable to connect to database")
	}

	redis, err := cache.NewRedisStore(1)
	if err != nil {
		log.Fatal("Unable to connect to redis")
	}

	session, err := authutil.NewSessionStore()
	if err != nil {
		log.Fatal("Unable to create session store")
	}

	server := api.NewServer(":8080", postgres, redis, session)
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// shutdown signal with grace period of 30 seconds to allow server to clean up
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("gracefully shutdown timed out.. forcing shutdown")
			}
		}()

		// shutdown server
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal("error shutting down server", err)
		}
		serverStopCtx()
	}()

	// start server
	err = server.Run()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("error running server", err)
	}

	// wait for server to stop
	<-serverCtx.Done()
}
