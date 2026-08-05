package main

import (
	"log"
	"net/http"

	"sharev2/internal/config"
	"sharev2/internal/handler"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listen:", cfg.ListenAddress)
	log.Println("Data:", cfg.DataDirectory)

	mux := handler.NewMux()

	if err := http.ListenAndServe(cfg.ListenAddress, mux); err != nil {
		log.Fatal(err)
	}
}
