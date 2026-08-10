package handler

import (
	"net/http"
	"sharev2/internal/download"
	"sharev2/internal/upload"
)

type Handler struct {
	uploadManager   *upload.Manager
	downloadManager *download.Manager
	webRoot         string
}

func NewHandler(um *upload.Manager, dm *download.Manager, webRoot string) *Handler {
	return &Handler{
		uploadManager:   um,
		downloadManager: dm,
		webRoot:         webRoot,
	}
}

func (h *Handler) NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(h.webRoot)))
	mux.HandleFunc("GET /download/{id}", h.download)
	mux.HandleFunc("GET /3usSroyU/files", h.listFiles)
	mux.HandleFunc("POST /upload", h.createSession)
	mux.HandleFunc("PATCH /upload/{id}", h.upload)
	mux.HandleFunc("POST /upload/{id}/commit", h.commit)

	return mux
}
