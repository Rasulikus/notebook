package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Rasulikus/notebook/internal/app"
	"github.com/Rasulikus/notebook/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler: app.App(cfg),
	}

	log.Printf("Server start at addres: " + server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
