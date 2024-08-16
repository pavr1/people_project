package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/pavr1/people_project/people/handlers/auth"
	repohandler "github.com/pavr1/people_project/people/handlers/repo"
	"github.com/pavr1/people_project/people/models"
)

type HttpHandler struct {
	log  *log.Logger
	repo *repohandler.RepoHandler
	auth *auth.Auth
}

func NewHttpHandler(auth *auth.Auth, repo *repohandler.RepoHandler, log *log.Logger) *HttpHandler {
	return &HttpHandler{
		auth: auth,
		repo: repo,
		log:  log,
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, 200}
}

func (h *HttpHandler) Middleware(peopleHandler http.HandlerFunc, prometheusHandler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		peopleHandler.ServeHTTP(w, r)
		timeEnd := time.Now()

		route := mux.CurrentRoute(r)
		path, err := route.GetPathTemplate()
		if err != nil {
			h.log.WithError(err).Error("Failed to get path template")
		}

		rw := newResponseWriter(w)

		r.Header.Add("X-Request-Path", path)
		r.Header.Add("X-Response-Status", string(rw.statusCode))
		r.Header.Add("X-Response-Time", timeEnd.Sub(timeStart).String())

		h.log.WithFields(log.Fields{"path": r.Header.Get("X-Request-Path"), "status": r.Header.Get("X-Response-Status"), "time": r.Header.Get("X-Response-Time")}).Info("Middleware Completed")

		prometheusHandler.ServeHTTP(w, r)
	})
}

func (h *HttpHandler) GetPersonList(w http.ResponseWriter, r *http.Request) {
	h.log.Info("GetPersonList")

	isValid := h.validate(r, w, http.MethodGet)
	if !isValid {
		return
	}

	people, err := h.repo.GetPersonList()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	bytes, err := json.Marshal(people)
	if err != nil {
		h.log.WithError(err).Error("Failed to marshal person list")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (h *HttpHandler) GetPerson(w http.ResponseWriter, r *http.Request) {
	h.log.Info("GetPerson")

	isValid := h.validate(r, w, http.MethodGet)
	if !isValid {
		return
	}

	id := r.PathValue("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID is required"))

		return
	}

	person, err := h.repo.GetPerson(id)
	if err != nil {
		//will need to check for not found
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	if person == nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Person not found"))

		return
	}

	bytes, err := json.Marshal(person)
	if err != nil {
		h.log.WithError(err).Error("Failed to marshal person")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (h *HttpHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	h.log.Info("CreatePerson")

	isValid := h.validate(r, w, http.MethodPost)
	if !isValid {
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.WithError(err).Error("Failed to read request body")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	person := models.Person{}
	err = json.Unmarshal(body, &person)
	if err != nil {
		h.log.WithError(err).Error("Failed to unmarshal request body")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = h.repo.CreatePerson(&person)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	//might need to change this to return the created person
	w.Write([]byte("Person successfully created"))
}

func (h *HttpHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	h.log.Info("UpdatePerson")

	isValid := h.validate(r, w, http.MethodPut)
	if !isValid {
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.WithError(err).Error("Failed to read request body")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	person := models.Person{}

	err = json.Unmarshal(body, &person)
	if err != nil {
		h.log.WithError(err).Error("Failed to unmarshal request body")

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if person.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID is required"))

		return
	}
	if person.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Name is required"))

		return
	}
	if person.LastName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("LastName is required"))

		return
	}

	if person.Age == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Age is required"))

		return
	}

	err = h.repo.UpdatePerson(&person)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Person successfully updated"))
}

func (h *HttpHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	h.log.Info("DeletePerson")

	isValid := h.validate(r, w, http.MethodDelete)
	if !isValid {
		return
	}

	id := r.PathValue("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID is required"))

		return
	}

	err := h.repo.DeletePerson(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Person with id " + id + " successfully deleted"))
}

func (h *HttpHandler) validate(r *http.Request, w http.ResponseWriter, method string) bool {
	isValid := h.isValidRequest(r, w, method)
	if !isValid {
		return false
	}

	isValid = h.isValidToken(r, w)

	return isValid
}

func (h *HttpHandler) isValidToken(r *http.Request, w http.ResponseWriter) bool {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	resCode, body, err := h.auth.IsValidToken(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.log.Warn("Failed to validate token")

		return false
	}

	if resCode != http.StatusOK {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(body))
		h.log.Warn("Invalid token")

		return false
	}

	h.log.Info("Valid token")
	return true
}

func (h *HttpHandler) isValidRequest(r *http.Request, w http.ResponseWriter, method string) bool {
	if r.Method != method {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		h.log.Warn("Invalid request method")
		return false
	}

	if r.Header.Get("Authorization") == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Authorization header is required"))
		h.log.Warn("Authorization header is required")

		return false
	}

	if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Authorization header is invalid"))
		h.log.Warn("Authorization header is invalid")

		return false
	}

	h.log.Info("Valid request")
	return true
}

func (h *HttpHandler) PrometheusLog(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(log.Fields{"path": r.Header.Get("X-Request-Path"), "status": r.Header.Get("X-Response-Status"), "time": r.Header.Get("X-Response-Time")}).Info("Middleware Prometheus")

	req, err := http.NewRequest("GET", "http://prometheus:9000/prometheus/log", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.WithError(err).Error("Failed to create prometheus request")

		return
	}

	req.Header = r.Header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.WithError(err).Error("Failed to send prometheus request")

		return
	}

	defer resp.Body.Close()
}
