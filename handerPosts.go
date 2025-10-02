package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	db "github.com/mike-woub/User_Auth/db/sqlc"
)

func (apiCfg *apiConfig) handlerCreatePost(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	type postParams struct {
		Category  string     `json:"category"`
		Title     string     `json:"title"`
		Content   string     `json:"content"`
		ImageUrl  []string   `json:"image_url"`
		ExpiresAt *time.Time `json:"expires_at,omitempty"`
	}

	var params postParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "invalid Json", http.StatusBadRequest)
		return
	}
	var expires sql.NullTime
	if params.ExpiresAt != nil {
		expires = sql.NullTime{Time: *params.ExpiresAt, Valid: true}
	} else {
		expires = sql.NullTime{Valid: false}
	}

	post, err := apiCfg.DB.CreatePost(r.Context(), db.CreatePostParams{
		UserID:    userID,
		Category:  params.Category,
		Title:     params.Title,
		Content:   params.Content,
		ImageUrl:  params.ImageUrl,
		ExpiresAt: expires,
	})
	if err != nil {
		log.Printf("CreatePost error: %v", err)
		http.Error(w, "failed to create post", http.StatusInternalServerError)
		return
	}
	w.WriteHeader((http.StatusCreated))
	json.NewEncoder(w).Encode(post)

}
