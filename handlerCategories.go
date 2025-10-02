package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (apiCfg *apiConfig) handlerGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := apiCfg.DB.GetCategories(r.Context())
	if err != nil {
		log.Printf("GetCategories error: %v", err)
		http.Error(w, "failed to fetch categores", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}
