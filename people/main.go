package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pavr1/people_project/people/config"
	"github.com/pavr1/people_project/people/handlers/auth"
	_http "github.com/pavr1/people_project/people/handlers/http"
	"github.com/pavr1/people_project/people/handlers/repo"
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

	repoHandler, err := repo.NewRepoHandler(log, config)
	if err != nil {
		log.WithError(err).Error("Failed to create repo handler")

		return
	}

	authHandler := auth.NewAuth(log, config)
	httpHandler := _http.NewHttpHandler(authHandler, repoHandler, log)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("/person/list", httpHandler.Middleware(httpHandler.GetPersonList, httpHandler.PrometheusLog))
	router.HandleFunc("/person/create", httpHandler.Middleware(httpHandler.CreatePerson, httpHandler.PrometheusLog))
	router.HandleFunc("/person/update", httpHandler.Middleware(httpHandler.UpdatePerson, httpHandler.PrometheusLog))
	router.HandleFunc("/person/delete/{id}", httpHandler.Middleware(httpHandler.DeletePerson, httpHandler.PrometheusLog))
	router.HandleFunc("/person/{id}", httpHandler.Middleware(httpHandler.GetPerson, httpHandler.PrometheusLog))

	log.WithField("port", config.Server.Port).Info("Listening to Server...")
	// Start the HTTP server
	log.Error(http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), router))
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
