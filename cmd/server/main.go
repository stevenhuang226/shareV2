package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"sharev2/internal/config"
	"sharev2/internal/download"
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
	log.Println("DataDirectory:", cfg.DataDirectory)
	log.Println("MetaDataDirectory:", cfg.MetaDataDirectory)

	storage := &storage.Storage{
		DataRoot:     cfg.DataDirectory,
		MetaDataRoot: cfg.MetaDataDirectory,
	}

	uploadManager, err := upload.CreateManager(storage)
	if err != nil {
		log.Fatal(err)
	}
	downloadManager := download.NewManager(storage)

	cleanup := scheduler.NewCleanupScheduler(
		uploadManager,
		time.Minute,
	)

	go cleanup.Run(ctx)

	httpHandler := handler.NewHandler(uploadManager, downloadManager, cfg.WebDirectory)

	mux := httpHandler.NewMux()

	if err := http.ListenAndServe(cfg.ListenAddress, mux); err != nil {
		log.Fatal(err)
	}
}
