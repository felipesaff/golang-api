package main

import (
	"crud-api/config"
	"crud-api/handlers"
	"log"
	"net/http"
	"time"
)

func main() {
	db := config.LoadDB()
	defer db.Close()

	mux := http.NewServeMux()
	handlers.RegisterHandlers(mux, db)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Server is running on port 8080")
	log.Fatal(srv.ListenAndServe())
}
