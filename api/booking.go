package api

import (
	"github.com/Tus1688/library-management-api/jsonutil"
	"github.com/Tus1688/library-management-api/types"
	"net/http"
	"strconv"
)

func (s *Server) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req types.CreateBooking
	if err := jsonutil.ShouldBind(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uid := r.Context().Value("uid").(string)

	bookingId, statusCode, err := s.store.CreateBooking(&uid, &req)
	if err.Error != "" {
		err := jsonutil.Render(w, statusCode, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	errResp := jsonutil.Render(w, http.StatusCreated, bookingId)
	if errResp != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) ReturnBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	statusCode, err := s.store.ReturnBook(&id)
	if err.Error != "" {
		err := jsonutil.Render(w, statusCode, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetBooking(w http.ResponseWriter, r *http.Request) {
	lastId, _ := strconv.Atoi(r.URL.Query().Get("last_id"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	bookings, statusCode, err := s.store.GetBooking(&lastId, &limit)
	if err.Error != "" {
		err := jsonutil.Render(w, statusCode, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	errResp := jsonutil.Render(w, http.StatusOK, bookings)
	if errResp != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
