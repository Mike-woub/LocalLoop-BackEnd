package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	db "github.com/mike-woub/User_Auth/db/sqlc"
)

type apiConfig struct {
	DB *db.Queries
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
		DB: queries,
	}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // your frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running on port" + port))
	})
	r.Post("/signup", apiCfg.handlerSignup)
	r.Post("/login", apiCfg.handlerGetUser)
	r.Get("/categories", apiCfg.handlerGetCategories)
	r.With(jwtMiddleware).Post("/posts", apiCfg.handlerCreatePost)

	fmt.Println("starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
