package handler

import (
	"net/http"
)

/*
func download(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	res := "try to download:" + r.PathValue("id")

	_, _ = w.Write([]byte(res))
}
*/

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
