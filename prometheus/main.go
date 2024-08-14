package main

import (
	"fmt"
	"net/http"

	"os"

	"github.com/gorilla/mux"
	"github.com/pavr1/prometheus/config"
	"github.com/pavr1/prometheus/handler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := mux.NewRouter()

	log := setupLogger()

	config, err := config.NewConfig()
	if err != nil {
		log.WithError(err).Error("Failed to create config")
		return
	}

	prometheus := handler.NewPrometheusHandler(log, config)
	router.Use(prometheus.PrometheusMiddleware)

	// Prometheus endpoint
	router.Path("/prometheus").Handler(promhttp.Handler())

	// Serving static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.WithField("Port", config.Prometheus.Port).Info("Serving prometheus requests")
	log.Error(http.ListenAndServe(fmt.Sprintf(":%d", config.Prometheus.Port), router))
}

func setupLogger() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	logger.SetReportCaller(true)
	logger.SetLevel(log.DebugLevel)

	// Set the output to stdout
	logger.SetOutput(os.Stdout)

	return logger
}
