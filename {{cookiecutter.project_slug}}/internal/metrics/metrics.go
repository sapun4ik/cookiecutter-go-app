package metrics

import (
	"fmt"

	"{{cookiecutter.module_path}}/pkg/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type SearchMicroserviceMetrics struct {
	SuccessHTTPRequests prometheus.Counter
	ErrorHTTPRequests   prometheus.Counter

	HTTPSuccessPingRequests prometheus.Counter
}

func NewSearchMicroserviceMetrics(cfg *config.Config) *SearchMicroserviceMetrics {
	return &SearchMicroserviceMetrics{

		SuccessHTTPRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_http_requests_total", cfg.Application.Name),
			Help: "The total number of success http requests",
		}),
		ErrorHTTPRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_http_requests_total", cfg.Application.Name),
			Help: "The total number of error http requests",
		}),
		HTTPSuccessPingRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_http_ping_requests_total", cfg.Application.Name),
			Help: "The total number of s_success_http_ping_requests_total requests",
		}),
	}
}
