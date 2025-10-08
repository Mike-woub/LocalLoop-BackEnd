package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	db "github.com/mike-woub/User_Auth/db/sqlc"
)

func (apiCfg *apiConfig) handlerLikePost(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	postIDStr := chi.URLParam(r, "post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	_, err = apiCfg.DB.LikePost(r.Context(), db.LikePostParams{
		UserID: sql.NullInt32{Int32: int32(userID), Valid: true},
		PostID: sql.NullInt32{Int32: int32(postID), Valid: true},
	})
	if err != nil {
		log.Printf("LikePost error: %v", err)
		http.Error(w, "failed to like post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (apiCfg *apiConfig) handlerUnlikePost(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	postIDStr := chi.URLParam(r, "post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	_, err = apiCfg.DB.UnlikePost(r.Context(), db.UnlikePostParams{
		UserID: sql.NullInt32{Int32: int32(userID), Valid: true},
		PostID: sql.NullInt32{Int32: int32(postID), Valid: true},
	})
	if err != nil {
		log.Printf("UnlikePost error: %v", err)
		http.Error(w, "failed to unlike post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (apiCfg *apiConfig) handlerGetLikeStatus(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	postIDStr := chi.URLParam(r, "post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	liked, err := apiCfg.DB.CheckLiked(r.Context(), db.CheckLikedParams{
		UserID: sql.NullInt32{Int32: int32(userID), Valid: true},
		PostID: sql.NullInt32{Int32: int32(postID), Valid: true},
	})
	if err != nil {
		log.Printf("CheckLiked error: %v", err)
		http.Error(w, "failed to check like status", http.StatusInternalServerError)
		return
	}

	count, err := apiCfg.DB.GetLikeCount(r.Context(), sql.NullInt32{Int32: int32(postID), Valid: true})
	if err != nil {
		log.Printf("GetLikeCount error: %v", err)
		http.Error(w, "failed to get like count", http.StatusInternalServerError)
		return
	}

	resp := struct {
		Liked bool `json:"liked"`
		Count int  `json:"count"`
	}{
		Liked: liked,
		Count: int(count),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
