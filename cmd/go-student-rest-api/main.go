package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rajwanraju/go-stundent-rest-api/internal/config"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup
	//setup router

	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	//setup server

	server := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
	}

	fmt.Printf("server started %s", cfg.HTTPServer.Address)
	server.ListenAndServe()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("cannot start server: %s", err.Error())
		}
	}()

	<-done
	slog.Info("Sutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("cannot shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server gracefully stopped")
	os.Exit(0)
}
