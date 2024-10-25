package controller

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsHandler struct {
	Log *log.Logger
	Db  *sqlx.DB
}

func NewMetricsHandler(logger *log.Logger, db *sqlx.DB) *MetricsHandler {
	return &MetricsHandler{
		Log: logger.WithPrefix("Metrics"),
		Db:  db,
	}
}

func (c *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Log.Debug("got metrics request")

	handler := promhttp.Handler()
	handler.ServeHTTP(w, r)
}

func (*MetricsHandler) Pattern() string {
	return "/metrics"
}
