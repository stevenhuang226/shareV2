package handler

import (
	"net/http"
	"sharev2/internal/upload"
)

type Handler struct {
	uploadManager *upload.Manager
	webRoot       string
}

func NewHandler(um *upload.Manager, webRoot string) *Handler {
	return &Handler{
		uploadManager: um,
		webRoot:       webRoot,
	}
}

func (h *Handler) NewMux() *http.ServeMux {
	/*
		a helper function register mux
	*/
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(h.webRoot)))
	mux.HandleFunc("GET /download/{id}", download)
	mux.HandleFunc("GET /upload", h.createSession)
	mux.HandleFunc("PATCH /upload/{id}", h.upload)
	mux.HandleFunc("POST /upload/{id}/commit", h.commit)

	return mux
}
