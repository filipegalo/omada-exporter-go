package main

import (
	"net/http"

	"omada_exporter_go/internal"
	"omada_exporter_go/internal/log"
	"omada_exporter_go/internal/prometheus"
)

func main() {
	Log.Init()
	Log.Info("Starting %s", internal.AppName)
	http.Handle("/metrics", Prometheus.OmadaMetricsHandler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
