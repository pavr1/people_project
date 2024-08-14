package prometheus

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pavr1/people/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type PrometheusHandler struct {
	log            *log.Logger
	config         *config.Config
	totalRequests  *prometheus.CounterVec
	responseStatus *prometheus.CounterVec
	httpDuration   *prometheus.HistogramVec
}

func (h *PrometheusHandler) init() {
	prometheus.MustRegister(h.totalRequests)
	prometheus.MustRegister(h.responseStatus)
	prometheus.MustRegister(h.httpDuration)
}

func NewPrometheusHandler(log *log.Logger, config *config.Config) *PrometheusHandler {
	totalRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"path"},
	)

	responseStatus := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Status of HTTP response",
		},
		[]string{"status"},
	)

	httpDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path1"})

	return &PrometheusHandler{
		log:            log,
		config:         config,
		totalRequests:  totalRequests,
		responseStatus: responseStatus,
		httpDuration:   httpDuration,
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (h *PrometheusHandler) PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(h.httpDuration.WithLabelValues(path))
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		h.responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		h.totalRequests.WithLabelValues(path).Inc()

		timer.ObserveDuration()

		log.WithFields(log.Fields{"status": statusCode, "path": path, "timer": timer.ObserveDuration()}).Info("PrometheusMiddleware Executed")
	})
}

func (h *PrometheusHandler) Listen() {
	router := mux.NewRouter()
	router.Use(h.PrometheusMiddleware)

	// Prometheus endpoint
	router.Path("/prometheus").Handler(promhttp.Handler())

	// Serving static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	h.log.WithField("Port", h.config.Prometheus.Port).Info("Serving prometheus requests")
	log.Error(http.ListenAndServe(fmt.Sprintf(":%d", h.config.Prometheus.Port), router))
}
