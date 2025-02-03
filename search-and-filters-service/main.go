package main

import (
	"log"
	"net/http"

	"search-and-filters-service/config"
	"search-and-filters-service/handlers"

	"github.com/gorilla/mux"
)

func main() {
	config.LoadConfig()

	router := mux.NewRouter()
	router.HandleFunc("/search", handlers.SearchHandler).Methods("GET")
	router.HandleFunc("/index", handlers.IndexHandler).Methods("POST")

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
