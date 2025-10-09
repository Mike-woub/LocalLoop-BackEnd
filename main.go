package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	db "github.com/mike-woub/User_Auth/db/sqlc"
	"github.com/patrickmn/go-cache"
)

type apiConfig struct {
	DB       *db.Queries
	OTPStore *cache.Cache
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set up in environment")
	}
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can not connect to database: ", err)
	}
	queries := db.New(conn)
	apiCfg := apiConfig{
		DB:       queries,
		OTPStore: cache.New(5*time.Minute, 10*time.Minute), // default expiration and cleanup interval
	}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",        // local dev
			"http://localhost:5175",        // alternate dev port
			"https://localloop.vercel.app", // your Vercel frontend
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running on port" + port))
	})
	r.Post("/request-signup-otp", apiCfg.handlerSignup)
	r.Post("/verify-signup", apiCfg.handlerVerifySignup)
	r.Post("/verify-password-reset", apiCfg.handlerVerifyPasswordReset)
	r.Post("/request-password-reset", apiCfg.handlerRequestPasswordReset)
	r.Post("/login", apiCfg.handlerGetUser)
	r.Get("/categories", apiCfg.handlerGetCategories)
	r.With(jwtMiddleware).Post("/posts", apiCfg.handlerCreatePost)
	r.With(jwtMiddleware).Post("/comments", apiCfg.handlerCreateComments)
	r.Get("/comments", apiCfg.handlerGetComments)
	r.Get("/posts", apiCfg.handlerGetPosts)
	r.Get("/posts/{post_id}", apiCfg.handlerGetCertainPost)
	r.With(jwtMiddleware).Post("/posts/{post_id}/like", apiCfg.handlerLikePost)
	r.With(jwtMiddleware).Delete("/posts/{post_id}/like", apiCfg.handlerUnlikePost)
	r.With(jwtMiddleware).Get("/posts/{post_id}/likes", apiCfg.handlerGetLikeStatus)
	r.With(jwtMiddleware).Delete("/posts/{post_id}", apiCfg.handlerDeletePosts)
	r.With(jwtMiddleware).Put("/users/username", apiCfg.handlerUpdateUsername)
	r.With(jwtMiddleware).Put("/users/email", apiCfg.handlerUpdateEmail)
	r.With(jwtMiddleware).Put("/users/password", apiCfg.handlerUpdatePassword)
	r.With(jwtMiddleware).Put("/users/avatar", apiCfg.handlerUpdateAvatar)

	fmt.Println("starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
