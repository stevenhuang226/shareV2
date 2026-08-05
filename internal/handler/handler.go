package handler

import "net/http"

type Handler struct {
}

func NewMux() *http.ServeMux {
	/*
		a helper function register mux
	*/
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", root)
	mux.HandleFunc("GET /download/{id}", download)
	mux.HandleFunc("POST /upload", upload)

	return mux
}
