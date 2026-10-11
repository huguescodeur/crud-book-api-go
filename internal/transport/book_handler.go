package transport

import (
	"crud-book/internal/model"
	"crud-book/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) HandlerBook(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		books, err := h.service.GetAllBook()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"total": len(books),
			"data":  books,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		var book model.Book

		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		createdBook, err := h.service.CreateBook(&book)
		if err != nil {
			if errors.Is(err, service.ErrAllFieldsRequired) {
				http.Error(w, "all fields required", http.StatusBadRequest)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(createdBook); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) HandlerBookByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		book, err := h.service.GetBookByID(id)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "data not found", http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodPut:
		var b model.Book

		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		updatedBook, err := h.service.UpdateBook(id, &b)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "data not found", http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(updatedBook); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodDelete:
		if err := h.service.DeleteBook(id); err != nil {
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "data not exist", http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

}
