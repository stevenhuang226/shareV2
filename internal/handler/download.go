package handler

import "net/http"

func download(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	res := "try to download:" + r.PathValue("id")

	_, _ = w.Write([]byte(res))
}
