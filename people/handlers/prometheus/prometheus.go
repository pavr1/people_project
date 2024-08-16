package prometheus

// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/gorilla/mux"
// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/promauto"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// 	log "github.com/sirupsen/logrus"
// )

// type PrometheusHandler struct {
// 	log            *log.Logger
// 	client         *http.Client
// }

// func NewPrometheusHandler(log *log.Logger) PrometheusHandler {
// 	return PrometheusHandler{
// 		log:    log,
// 		client: &http.Client{},
// 	}
// }

// type responseWriter struct {
// 	http.ResponseWriter
// 	statusCode int
// }

// func newResponseWriter(w http.ResponseWriter) *responseWriter {
// 	return &responseWriter{w, http.StatusOK}
// }

// func (h *PrometheusHandler) PrometheusMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		route := mux.CurrentRoute(r)
// 		path, _ := route.GetPathTemplate()

// 		log.WithField("path", path).Info("PrometheusMiddleware Executing")

// 		timer := prometheus.NewTimer(h.httpDuration.WithLabelValues(path))
// 		rw := newResponseWriter(w)
// 		next.ServeHTTP(rw, r)

// 		statusCode := rw.statusCode

// 		log.WithField("status", statusCode).Info("PrometheusMiddleware Executing")

// 		h.responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
// 		h.totalRequests.WithLabelValues(path).Inc()

// 		timer.ObserveDuration()

// 		log.WithFields(log.Fields{"status": statusCode, "path": path, "timer": timer.ObserveDuration()}).Info("PrometheusMiddleware Executed")
// 	})
// }

// func (h *PrometheusHandler) Listen(log *log.Logger, port int) {
// 	go func() {
// 		router := mux.NewRouter()

// 		Prometheus endpoint
// 		router.Path("/prometheus").Handler(promhttp.Handler())

// 		Serving static files
// 		router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

// 		log.WithField("Port", port).Info("Serving prometheus requests")
// 		log.Error(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
// 	}()
// }
