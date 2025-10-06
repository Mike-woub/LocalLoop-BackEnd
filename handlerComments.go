package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "github.com/mike-woub/User_Auth/db/sqlc"
)

func (apiCfg apiConfig) handlerCreateComments(w http.ResponseWriter, r *http.Request) {
	type commentParams struct {
		PostID  string `json:"post_id"`
		UserID  string `json:"user_id"`
		Content string `json:"content"`
	}

	var params commentParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		fmt.Printf("json:%v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Convert PostID and UserID to integers
	postID, err := strconv.Atoi(params.PostID)
	if err != nil {
		http.Error(w, "Invalid post_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(params.UserID)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Create comment
	comment, err := apiCfg.DB.CreateComment(r.Context(), db.CreateCommentParams{
		PostID:    int32(postID),
		UserID:    int32(userID),
		Content:   params.Content,
		CreatedAt: time.Now(),
	})
	if err != nil {
		http.Error(w, "Error creating comment", http.StatusInternalServerError)
		return
	}

	// Respond with the created comment
	response := struct {
		ID        int32     `json:"id"`
		PostID    int32     `json:"post_id"`
		UserID    int32     `json:"user_id"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}{
		ID:        comment.ID,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (apiCfg *apiConfig) handlerGetComments(w http.ResponseWriter, r *http.Request) {
	// Extract post_id from query string
	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		http.Error(w, "Missing post_id", http.StatusBadRequest)
		return
	}

	// Convert post_id to int32
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post_id", http.StatusBadRequest)
		return
	}

	// Fetch comments from DB
	comments, err := apiCfg.DB.GetComments(r.Context(), int32(postID))
	if err != nil {
		http.Error(w, "Can't fetch comments", http.StatusInternalServerError)
		return
	}

	// Respond with comments
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
