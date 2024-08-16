package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

type PrometheusHandler struct {
	log            *log.Logger
	totalRequests  *prometheus.CounterVec
	responseStatus *prometheus.CounterVec
	httpDuration   *prometheus.HistogramVec
}

func (h *PrometheusHandler) init() {
	prometheus.MustRegister(h.totalRequests)
	prometheus.MustRegister(h.responseStatus)
	prometheus.MustRegister(h.httpDuration)
}

func NewPrometheusHandler(log *log.Logger) *PrometheusHandler {
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
		totalRequests:  totalRequests,
		responseStatus: responseStatus,
		httpDuration:   httpDuration,
	}
}

func (h *PrometheusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.Header.Get("X-Request-Path")
	status := r.Header.Get("X-Response-Status")
	time := r.Header.Get("X-Response-Time")

	if path != "" {
		log.WithField("path", path).Info("PrometheusMiddleware Executing")
		h.totalRequests.WithLabelValues(path).Inc()
	}

	if status != "" {
		log.WithField("status", status).Info("PrometheusMiddleware Executing")
		h.responseStatus.WithLabelValues(status).Inc()
	}

	if time != "" {
		time = strings.ReplaceAll(time, "ms", "")
		timeF, err := strconv.ParseFloat(time, 64)
		if err != nil {
			log.WithError(err).Error("Failed to parse X-Response-Time to float64")

		} else {
			log.WithField("time", time).Info("PrometheusMiddleware Executing")
			h.httpDuration.WithLabelValues(path).Observe(timeF)
		}
	}
}
