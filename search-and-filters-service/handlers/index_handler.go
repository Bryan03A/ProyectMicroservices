package handlers

import (
	"encoding/json"
	"net/http"

	"search-and-filters-service/models"
	"search-and-filters-service/services"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var doc models.Document
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := services.Index(doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
