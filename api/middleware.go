package api

import (
	"context"
	"net/http"
)

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
