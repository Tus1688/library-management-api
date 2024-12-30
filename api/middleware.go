package api

import (
	"context"
	"net/http"
)

// EnforceAuthentication is a middleware that enforces authentication on incoming HTTP requests.
// It checks for the presence of an access token in the request cookies and validates it.
// If the token is valid, it optionally passes the user ID to the request context.
// Parameters:
// - expiredIn: The time-to-live (TTL) for the token in seconds.
// - passUserId: A boolean indicating whether to pass the user ID to the request context.
// Returns:
// - A middleware function that wraps the next HTTP handler.
func (s *Server) EnforceAuthentication(expiredIn uint32, passUserId bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			access, err := r.Cookie("access")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			uid, errResp := s.session.ValidateToken(access.Value, expiredIn)
			if errResp.Error != "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if passUserId {
				ctx := context.WithValue(r.Context(), "uid", uid)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
