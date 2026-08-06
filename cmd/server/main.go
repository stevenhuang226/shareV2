package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"sharev2/internal/config"
	"sharev2/internal/handler"
	"sharev2/internal/scheduler"
	"sharev2/internal/storage"
	"sharev2/internal/upload"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listen:", cfg.ListenAddress)
	log.Println("Data:", cfg.DataDirectory)

	storage := &storage.Storage{
		RootPath: "./test",
	}

	uploadManager, err := upload.CreateManager(storage)
	if err != nil {
		log.Fatal(err)
	}

	cleanup := scheduler.NewCleanupScheduler(
		uploadManager,
		time.Minute,
	)

	go cleanup.Run(ctx)

	httpHandler := handler.NewHandler(uploadManager)

	mux := httpHandler.NewMux()

	if err := http.ListenAndServe(cfg.ListenAddress, mux); err != nil {
		log.Fatal(err)
	}
}
