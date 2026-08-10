package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) listFiles(w http.ResponseWriter, r *http.Request) {
	metaData, err := h.downloadManager.ListMetaData()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(metaData)
}
