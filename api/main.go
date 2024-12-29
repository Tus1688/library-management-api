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

type Server struct {
	store   storage.Storage
	cache   cache.Cache
	session authutil.Session

	server *http.Server
}

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

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.store.Shutdown(); err != nil {
		return err
	}

	if err := s.cache.Shutdown(); err != nil {
		return err
	}

	return s.server.Shutdown(ctx)
}

func (s *Server) Run() error {
	log.Print("server is running on ", s.server.Addr)
	return s.server.ListenAndServe()
}

func handler(s *Server) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", s.Login)
		})
	})

	_ = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})

	return r
}
