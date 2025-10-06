package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	db "github.com/mike-woub/User_Auth/db/sqlc"
)

func (apiCfg *apiConfig) handlerCreatePost(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	type postParams struct {
		CategoryID int        `json:"category_id"`
		Title      string     `json:"title"`
		Content    string     `json:"content"`
		ImageUrl   []string   `json:"image_url"`
		ExpiresAt  *time.Time `json:"expires_at,omitempty"`
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
		expires = sql.NullTime{Time: time.Now().AddDate(0, 0, 14), Valid: true}
	}
	if len(params.ImageUrl) == 0 {
		params.ImageUrl = []string{"https://cdn3.iconfinder.com/data/icons/news-65/64/paper_plane-send-message-mail-communication-publish-origami-512.png"}
	}
	if params.Title == "" || params.CategoryID == 0 || params.Content == "" {
		log.Printf("Invalid post params: %+v", params)
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	post, err := apiCfg.DB.CreatePost(r.Context(), db.CreatePostParams{
		UserID:     userID,
		CategoryID: int32(params.CategoryID),
		Title:      params.Title,
		Content:    params.Content,
		ImageUrl:   params.ImageUrl,
		ExpiresAt:  expires,
	})

	if err != nil {
		log.Printf("CreatePost error: %v", err)
		http.Error(w, "failed to create post", http.StatusInternalServerError)
		return
	}
	w.WriteHeader((http.StatusCreated))
	json.NewEncoder(w).Encode(post)

}

func (apiCfg *apiConfig) handlerGetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := apiCfg.DB.GetPosts(r.Context())
	if err != nil {
		log.Printf("error fetching posts %v", err)
		http.Error(w, "cant fetch posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
func (apiCfg *apiConfig) handlerGetCertainPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	post, err := apiCfg.DB.GetCertainPost(r.Context(), int32(postID))
	if err != nil {
		log.Printf("error fetching post %v", err)
		http.Error(w, "couldnt fetch post", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}
