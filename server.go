package main

import (
	"fmt"

	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
)

func startServer() {
	readAndSetEnv(".env")

	port, _ := os.LookupEnv("PORT")
	router := setUpRouter()
	server := &http.Server{Handler: router, Addr: ":" + port}

	fmt.Println("Starting the server on Port:", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("Error while starting the server", err)
	}
}

func setUpRouter() (router *chi.Mux) {
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

	router.Mount("/v1", initV1Router())

	return router
}
