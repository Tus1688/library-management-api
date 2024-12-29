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

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	refresh, err := r.Cookie("refresh")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// verify the refresh signature
	if err := s.session.VerifyRefreshToken(&refresh.Value); err.Error != "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// delete the refresh token from the cache
	if err := s.cache.DeleteRefreshToken(&refresh.Value); err.Error != "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// delete access & refresh token
	access := http.Cookie{
		Name:     "access",
		Value:    "",
		Path:     "/api",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	}

	newRefresh := http.Cookie{
		Name:     "refresh",
		Value:    "",
		Path:     "/api",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		MaxAge:   -1,
	}

	http.SetCookie(w, &access)
	http.SetCookie(w, &newRefresh)

	w.WriteHeader(http.StatusOK)
}

func (s *Server) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("refresh")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// verify the authenticity of the token
	if err := s.session.VerifyRefreshToken(&token.Value); err.Error != "" {
		if err := jsonutil.Render(w, http.StatusUnauthorized, err); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	uid, errResp := s.cache.GetRefreshToken(&token.Value)
	if errResp.Error != "" {
		if err := jsonutil.Render(w, http.StatusInternalServerError, err); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	sessionToken, errResp := s.session.CreateSessionToken(&uid)
	if errResp.Error != "" {
		if err := jsonutil.Render(w, http.StatusInternalServerError, err); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	access := http.Cookie{
		Name:     "access",
		Value:    sessionToken,
		Path:     "/api",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	}

	http.SetCookie(w, &access)
	w.WriteHeader(http.StatusOK)
}
