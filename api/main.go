package api

import (
	"context"
	"fmt"
	"github.com/Tus1688/library-management-api/authutil"
	"github.com/Tus1688/library-management-api/cache"
	"github.com/Tus1688/library-management-api/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

// Server represents the API server with its dependencies.
type Server struct {
	store   storage.Storage
	cache   cache.Cache
	session authutil.Session

	server *http.Server
}

// NewServer creates a new Server instance.
// Parameters:
// - listenAddr: the address the server will listen on.
// - store: the storage backend.
// - cache: the cache backend.
// - session: the session manager.
// Returns a pointer to the created Server.
func NewServer(listenAddr string, store storage.Storage, cache cache.Cache, session authutil.Session) *Server {
	s := &Server{
		store:   store,
		cache:   cache,
		session: session,
	}

	s.server = &http.Server{
		Addr:    listenAddr,
		Handler: handler(s),
	}

	return s
}

// Shutdown gracefully shuts down the server and its dependencies.
// Parameters:
// - ctx: the context for shutdown.
// Returns an error if any of the shutdown operations fail.
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.store.Shutdown(); err != nil {
		return err
	}

	if err := s.cache.Shutdown(); err != nil {
		return err
	}

	return s.server.Shutdown(ctx)
}

// Run starts the server and listens for incoming requests.
// Returns an error if the server fails to start.
func (s *Server) Run() error {
	log.Print("server is running on ", s.server.Addr)
	return s.server.ListenAndServe()
}

// handler sets up the HTTP routes and middleware for the server.
// Parameters:
// - s: the server instance.
// Returns an http.Handler with the configured routes and middleware.
func handler(s *Server) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", s.Login)
			r.Post("/logout", s.Logout)
			r.Post("/refresh", s.RefreshToken)

			r.Route("/dashboard", func(r chi.Router) {
				r.Use(s.EnforceAuthentication(300, true))

				r.Get("/user", s.GetEmployee)
				r.Post("/user", s.CreateEmployee)
				r.Delete("/user", s.DeleteEmployee)
			})
		})

		r.Route("/collections", func(r chi.Router) {
			// public route
			r.Get("/book", s.GetBook)

			r.Route("/dashboard", func(r chi.Router) {
				r.Use(s.EnforceAuthentication(600, true))

				r.Post("/book", s.CreateBook)
				r.Put("/book", s.UpdateBook)
				r.Delete("/book", s.DeleteBook)

				r.Get("/booking", s.GetBooking)
				r.Post("/booking", s.CreateBooking)
				r.Post("/return", s.ReturnBook)
			})
		})
	})

	_ = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})

	return r
}
