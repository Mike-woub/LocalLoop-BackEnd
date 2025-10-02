package main

import (
	"encoding/json"
	"net/http"
)

func (apiCfg *apiConfig) handlerGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := apiCfg.DB.GetCategories(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch categores", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
