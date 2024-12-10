package gateway

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/GoJobs/config"
	"github.com/saleh-ghazimoradi/GoJobs/logger"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

func Server() error {
	router := registerRoutes()
	srv := &http.Server{
		Addr:         config.AppConfig.ServerConfig.Port,
		Handler:      router,
		IdleTimeout:  config.AppConfig.ServerConfig.IdleTimeout,
		ReadTimeout:  config.AppConfig.ServerConfig.ReadTimeout,
		WriteTimeout: config.AppConfig.ServerConfig.WriteTimeout,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		logger.Logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		logger.Logger.Info("completing background tasks", "addr", srv.Addr)

		wg.Wait()
		shutdownError <- nil
	}()

	logger.Logger.Info("starting server", "addr", config.AppConfig.ServerConfig.Port, "env", config.AppConfig.ServerConfig.Version)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	logger.Logger.Info("stopped server", "addr", srv.Addr)

	return nil
}