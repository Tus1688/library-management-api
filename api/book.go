package api

import (
	"github.com/Tus1688/library-management-api/jsonutil"
	"github.com/Tus1688/library-management-api/types"
	"net/http"
	"strconv"
)

func (s *Server) GetBook(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")
	lastId, _ := strconv.Atoi(r.URL.Query().Get("last_id"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	books, statusCode, err := s.store.GetBook(&searchQuery, &lastId, &limit)
	if err.Error != "" {
		err := jsonutil.Render(w, statusCode, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	errResp := jsonutil.Render(w, http.StatusOK, books)
	if errResp != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req types.CreateBook
	if err := jsonutil.ShouldBind(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bookId, statusCode, err := s.store.CreateBook(&req)
	if err.Error != "" {
		err := jsonutil.Render(w, statusCode, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	errResp := jsonutil.Render(w, http.StatusCreated, bookId)
	if errResp != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	statusCode, err := s.store.DeleteBook(&id)
	if err.Error != "" {
		err := jsonutil.Render(w, statusCode, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var req types.UpdateBook
	if err := jsonutil.ShouldBind(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	statusCode, err := s.store.UpdateBook(&req)
	if err.Error != "" {
		err := jsonutil.Render(w, statusCode, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
