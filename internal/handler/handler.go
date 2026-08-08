package handler

import (
	"net/http"
	"sharev2/internal/upload"
)

type Handler struct {
	uploadManager *upload.Manager
}

func NewHandler(um *upload.Manager) *Handler {
	return &Handler{
		uploadManager: um,
	}
}

func (h *Handler) NewMux() *http.ServeMux {
	/*
		a helper function register mux
	*/
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", root)
	mux.HandleFunc("GET /download/{id}", download)
	mux.HandleFunc("GET /upload", h.createSession)
	mux.HandleFunc("PATCH /upload/{id}", h.upload)
	mux.HandleFunc("POST /upload/{id}/commit", h.commit)

	return mux
}
