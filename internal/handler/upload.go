package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sharev2/internal/storage"
	"sharev2/internal/upload"
	"strconv"
)

type uploadResponse struct {
	ID     string `json:"id"`
	Offset int64  `json:"offset"`
}

func (h *Handler) upload(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	offsetStr := r.Header.Get("Upload-Offset")

	if id == "" || offsetStr == "" {
		http.Error(w, "wrong id or offset", http.StatusBadRequest)
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		http.Error(w, "wrong offset value", http.StatusBadRequest)
		return
	}

	n, err := h.uploadManager.Append(id, offset, r.Body)

	switch {
	case errors.Is(err, upload.ErrNotFound):
		http.Error(w, "session not found", http.StatusNotFound)
		return

	case errors.Is(err, storage.ErrOffsetMismatch),
		errors.Is(err, storage.ErrOffsetWrong):
		nextOffset, err := h.uploadManager.NextOffset(id)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusConflict)

		err = json.NewEncoder(w).Encode(uploadResponse{
			ID:     id,
			Offset: nextOffset,
		})

		if err != nil {
			log.Println(err)
		}

		return

	case err != nil:
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(uploadResponse{
		ID:     id,
		Offset: offset + n,
	})

	if err != nil {
		log.Println(err)
	}
}

func (h *Handler) createSession(w http.ResponseWriter, r *http.Request) {
	id, err := h.uploadManager.CreateSession()

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(uploadResponse{
		ID:     id,
		Offset: 0,
	})

	if err != nil {
		log.Println(err)
	}
}

func (h *Handler) commit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.uploadManager.Commit(id)

	switch {
	case errors.Is(err, upload.ErrNotFound):
		http.Error(w, "session not found", http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, "internal server error, commit failed", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
