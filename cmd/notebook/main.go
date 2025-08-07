package main

import (
	"log"
	"net/http"

	"github.com/Rasulikus/notebook/internal/app"
)

func main() {
	server := http.Server{
		Addr:    ":8081",
		Handler: app.App(),
	}

	log.Printf("Server start at port: " + server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
