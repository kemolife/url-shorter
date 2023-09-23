// package build short url base on provide from user url
// Example:
// https://test.domain.com/home?request_id=1&user_id=1 -> https://short.domain.com/xr4tff

// Store all related url data into MongoDB.
// Base CRUD for managing url data
// Provide yaml base configurations
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"shorter-url/internal/auth"
	"shorter-url/internal/config"
	"shorter-url/internal/db"
	"shorter-url/internal/github"
	"shorter-url/internal/server"
	"shorter-url/internal/shorter"
	"shorter-url/internal/storage/shortening"
	"shorter-url/internal/storage/user"
	"syscall"
	"time"
)

func main() {
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer dbCancel()

	dbClient, err := db.Connect(dbCtx, config.Get().Storage.Config["mongo_url"])

	if err != nil {
		log.Fatalf("storage init false: %s", err)
	}

	mongoDb := dbClient.GetClient().Database("url")

	var (
		shorteningStorage = shortening.NewMongoDB(mongoDb)
		userStorage       = user.NewMongoDB(mongoDb)
		shortener         = shorter.NewService(shorteningStorage)
		githubClient      = github.NewClient()
		authenticator     = auth.NewService(
			githubClient, userStorage, config.Get().GitHub.ClientId, config.Get().GitHub.ClientSecret,
		)
		srv = server.New(shortener, authenticator)
	)

	srv.AddCloser(dbClient.CloseConnection)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	u, _ := url.Parse(config.Get().Address)

	go func() {
		if err := http.ListenAndServe(u.Host, srv); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error running server: %v", err)
		}
	}()

	log.Println("server started")
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("error closing server: %v", err)
	}

	log.Println("server stopped")
}
