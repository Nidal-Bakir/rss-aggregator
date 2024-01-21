package main

import (
	"database/sql"
	"fmt"

	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Nidal-Bakir/rss-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
)

func startServer() {
	readAndSetEnv(".env")

	port, _ := os.LookupEnv("PORT")
	if port == "" {
		log.Fatal("con not find the port in hte env")
	}

	apiConfig := setupDatabase()

	router := setUpRouter(apiConfig)
	server := &http.Server{Handler: router, Addr: ":" + port}

	fmt.Println("Starting the server on Port:", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("Error while starting the server", err)
	}
}

func setUpRouter(apiConfig apiConfig) (router *chi.Mux) {
	router = chi.NewRouter()

	logger := httplog.NewLogger("rss-agg", httplog.Options{
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		Tags: map[string]string{
			"version": "v1.0",
			"env":     "dev",
		},
		QuietDownRoutes: []string{
			"/",
			"/ping",
		},
		ResponseHeaders: true,
		QuietDownPeriod: 10 * time.Second,
	})

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(httplog.RequestLogger(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/ping"))

	// surely you should not use this!!!
	router.Use(cors.AllowAll().Handler)

	router.Mount("/v1", initV1Router(apiConfig))

	return router
}

type apiConfig struct {
	DB *database.Queries
}

func setupDatabase() apiConfig {
	dbUrl, _ := os.LookupEnv("DB_URL")
	if dbUrl == "" {
		log.Fatal("con not find the DB_URL in hte env")
	}

	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal("can not open the database connection", err)
	}

	return apiConfig{DB: database.New(db)}

}
