package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ByChanderZap/rss-web-server/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port must be set on the enviroment")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL must be set on the enviroment")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	timeBetweenRequests, err := time.ParseDuration(os.Getenv("TIME_B_REQ"))
	if err != nil {
		log.Fatalf("error while parsing duration %v", err)
	}
	go startScrapping(db, 10, timeBetweenRequests)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"POST", "GET", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handleRediness)
	v1Router.Get("/err", handleErr)
	// User routes
	v1Router.Post("/users", apiCfg.handleCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUserByApiKey))
	// Feed routes
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	v1Router.Get("/feeds", apiCfg.handleGetFeeds)
	// Feed follow routes
	v1Router.Post("/feeds-follow", apiCfg.middlewareAuth(apiCfg.handleCreateFeedFollow))
	v1Router.Get("/feeds-follow", apiCfg.middlewareAuth(apiCfg.handleGetFeedsFollows))
	v1Router.Delete("/feeds-follow/{feedFollowId}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeedsFollows))
	// Posts
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handleGetPostsForUser))

	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server started on port %v", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
