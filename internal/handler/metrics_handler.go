package handler

import (
	"encoding/json"
	"net/http"

	"github.com/hans/config-service/internal/service"
)

type MetricsHandler struct {
	metricsSvc *service.MetricsService
}

func NewMetricsHandler(metricsSvc *service.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		metricsSvc: metricsSvc,
	}
}

// GetRequestsPerMinute GET /api/v1/metrics/requests-per-minute
// @Summary Get requests per minute metrics
// @Description Returns the number of requests per minute for the last 24 hours.
// @Tags Metrics
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Router /metrics/requests-per-minute [get]
func (h *MetricsHandler) GetRequestsPerMinute(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.metricsSvc.GetRequestsPerMinute(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
