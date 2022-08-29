package metrics

import (
	"fmt"

	"{{cookiecutter.module_path}}/pkg/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

//go:generate mockgen -source=metrics.go -destination=mocks/mock.go
type Metrics interface {
	SuccessHTTPRequestsInc()
	ErrorHTTPRequestsInc()
	HTTPSuccessPingRequestsInc()
}

type ServiceMetrics struct {
	successHTTPRequests     prometheus.Counter
	errorHTTPRequests       prometheus.Counter
	httpSuccessPingRequests prometheus.Counter
}

func (l ServiceMetrics) SuccessHTTPRequestsInc() {
	l.successHTTPRequests.Inc()
}

func (l ServiceMetrics) ErrorHTTPRequestsInc() {
	l.errorHTTPRequests.Inc()
}

func (l ServiceMetrics) HTTPSuccessPingRequestsInc() {
	l.httpSuccessPingRequests.Inc()
}

func NewServiceMetrics(cfg *config.Config) Metrics {
	return &ServiceMetrics{
		successHTTPRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_http_requests_total", cfg.Application.Name),
			Help: "The total number of success http requests",
		}),
		errorHTTPRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_http_requests_total", cfg.Application.Name),
			Help: "The total number of error http requests",
		}),
		httpSuccessPingRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_http_ping_requests_total", cfg.Application.Name),
			Help: "The total number of s_success_http_ping_requests_total requests",
		}),
	}
}
