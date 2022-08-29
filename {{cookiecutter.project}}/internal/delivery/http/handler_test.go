package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"{{cookiecutter.module_path}}/internal/metrics"
	mock_metrics "{{cookiecutter.module_path}}/internal/metrics/mocks"
	"{{cookiecutter.module_path}}/pkg/config"
	"{{cookiecutter.module_path}}/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {

	h := NewHandler(logger.NewAPILogger(&config.Config{}), &metrics.ServiceMetrics{})

	require.IsType(t, &Handler{}, h)
}

func TestNewHandler_Init(t *testing.T) {

	// Init Dependencies
	c := gomock.NewController(t)
	defer c.Finish()

	mockBehavior := func(m *mock_metrics.MockMetrics) {
		m.EXPECT().HTTPSuccessPingRequestsInc()
	}

	mockMetrics := mock_metrics.NewMockMetrics(c)
	mockBehavior(mockMetrics)

	h := NewHandler(logger.NewAPILogger(&config.Config{}), mockMetrics)

	router := h.Init(&config.Config{})

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)
}
