package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/jmnjung/rssagg/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	apiCfg := &apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", handlerHealthz)
	mux.HandleFunc("GET /v1/err", handlerErr)

	mux.HandleFunc("POST /v1/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	mux.HandleFunc("GET /v1/feeds", apiCfg.handlerGetAllFeeds)

	mux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))
	mux.HandleFunc("GET /v1/feed_follows/", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
