package api

import (
	"github.com/Tus1688/library-management-api/jsonutil"
	"github.com/Tus1688/library-management-api/types"
	"net/http"
)

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
	if err := jsonutil.ShouldBind(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uid, statusCode, err := s.store.Login(&req)
	if err.Error != "" {
		err := jsonutil.Render(w, statusCode, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	sessionToken, err := s.session.CreateSessionToken(&uid)
	if err.Error != "" {
		err := jsonutil.Render(w, http.StatusInternalServerError, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	refreshToken, err := s.session.CreateRefreshToken()
	if err.Error != "" {
		err := jsonutil.Render(w, http.StatusInternalServerError, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if err := s.cache.SaveRefreshToken(&refreshToken, &uid); err.Error != "" {
		err := jsonutil.Render(w, http.StatusInternalServerError, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = s.session.SignRefreshToken(&refreshToken)
	if err.Error != "" {
		err := jsonutil.Render(w, http.StatusInternalServerError, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	access := http.Cookie{
		Name:     "access",
		Value:    sessionToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/api",
		Secure:   true,
	}

	refresh := http.Cookie{
		Name:     "refresh",
		Value:    refreshToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/api",
		Secure:   true,
		// 24 hours
		MaxAge: 24 * 60 * 60,
	}

	http.SetCookie(w, &access)
	http.SetCookie(w, &refresh)

	w.WriteHeader(http.StatusOK)
}
