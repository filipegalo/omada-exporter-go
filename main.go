package main

import (
	"net/http"

	"omada_exporter_go/internal"
	"omada_exporter_go/internal/log"
	"omada_exporter_go/internal/prometheus"
)

func main() {
	Log.Init()
	conf := internal.GetConfig()
	Log.Info("Starting %s", internal.AppName)
	http.Handle(conf.Prometheus.MetricsPath, Prometheus.OmadaMetricsHandler())
	if err := http.ListenAndServe(":"+conf.Prometheus.MetricsPort, nil); err != nil {
		panic(err)
	}
}
