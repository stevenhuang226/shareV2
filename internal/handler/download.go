package handler

import (
	"net/http"
)

func (h *Handler) download(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	file, metaData, err := h.downloadManager.OpenDownload(id)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	http.ServeContent(w, r, metaData.Name, metaData.UploadTime, file)
}
