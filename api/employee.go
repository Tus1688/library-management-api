package api

import (
	"github.com/Tus1688/library-management-api/jsonutil"
	"github.com/Tus1688/library-management-api/types"
	"net/http"
)

func (s *Server) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var req types.CreateEmployee
	if err := jsonutil.ShouldBind(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, statusCode, err := s.store.CreateEmployee(&req)
	if err.Error != "" {
		err := jsonutil.Render(w, statusCode, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	errResp := jsonutil.Render(w, http.StatusCreated, userId)
	if errResp != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) GetEmployee(w http.ResponseWriter, _ *http.Request) {
	employee, statusCode, err := s.store.GetEmployee()
	if err.Error != "" {
		err := jsonutil.Render(w, statusCode, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	errResp := jsonutil.Render(w, http.StatusOK, employee)
	if errResp != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currentUserId := r.Context().Value("uid").(string)

	statusCode, err := s.store.DeleteEmployee(&currentUserId, &id)
	if err.Error != "" {
		err := jsonutil.Render(w, statusCode, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
