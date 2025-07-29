package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lakshsetia/learn-RESTAPI/internal/config"
	"github.com/lakshsetia/learn-RESTAPI/internal/handlers"
	"github.com/lakshsetia/learn-RESTAPI/internal/middlewares"
	"github.com/lakshsetia/learn-RESTAPI/internal/storage/postgresql"
)

func main() {
	// load config
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config file\n", err)
	}

	// load database
	postgresql, err := postgresql.New(config)
	if err != nil {
		log.Fatal("failed to connect to database\n", err)
	}

	// setup router
	router := http.NewServeMux()

	// setup handler
	router.Handle("/user", middlewares.Middleware(handlers.UserHandler(postgresql)))
	router.Handle("/user/", middlewares.Middleware(handlers.UserByIdHandler(postgresql)))

	// setup server
	server := http.Server{
		Addr: config.HTTPServer.Address,
		Handler: router,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	slog.Info("server starting at", slog.String("address", config.HTTPServer.Address))

	// start server
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("failed to start server\n", err)
		}
	}()
	<-done

	// shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	server.Shutdown(ctx)
}