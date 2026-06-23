package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"omada_exporter_go/internal"
	"omada_exporter_go/internal/log"
	"omada_exporter_go/internal/prometheus"
)

func main() {
	Log.Init()
	conf := internal.GetConfig()
	Log.Info("Starting %s", internal.AppName)

	mux := http.NewServeMux()
	mux.Handle(conf.Prometheus.MetricsPath, Prometheus.OmadaMetricsHandler())

	srv := &http.Server{
		Addr:         ":" + conf.Prometheus.MetricsPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		Log.Info("Listening on %s%s", srv.Addr, conf.Prometheus.MetricsPath)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			Log.Error(err, "HTTP server failed")
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	Log.Info("Shutdown signal received, draining connections")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		Log.Error(err, "Graceful shutdown failed")
		os.Exit(1)
	}
	Log.Info("Stopped cleanly")
}
