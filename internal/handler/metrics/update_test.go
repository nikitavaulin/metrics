package metricshandler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	models "github.com/nikitavaulin/metrics/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMetricsService struct {
	mock.Mock
}

func (ms *MockMetricsService) Add(name string, metric models.Metrics) error {
	args := ms.Called(name, metric)
	return args.Error(0)
}

func (ms *MockMetricsService) Get(name string) (any, error) {
	args := ms.Called(name)
	return args.Get(0), args.Error(1)
}

func (ms *MockMetricsService) GetList() (map[string]any, error) {
	args := ms.Called()
	if data := args.Get(0); data != nil {
		return data.(map[string]any), args.Error(1)
	}
	return nil, args.Error(1)
}

func (ms *MockMetricsService) SetMValueFromStorage(m *models.Metrics) error {
	args := ms.Called()
	return args.Error(0)
}

func TestUpdateMetrics(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		setupMock      func(*MockMetricsService)
		expectedStatus int
	}{
		{
			name:   "success counter",
			method: http.MethodPost,
			path:   "/update/counter/test_counter/10",
			setupMock: func(m *MockMetricsService) {
				delta := int64(10)
				m.On("Add", "test_counter", models.Metrics{
					MType: "counter",
					Delta: &delta,
				}).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "success gauge",
			method: http.MethodPost,
			path:   "/update/gauge/test_gauge/3.14",
			setupMock: func(m *MockMetricsService) {
				value := 3.14
				m.On("Add", "test_gauge", models.Metrics{
					MType: "gauge",
					Value: &value,
				}).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty metric type",
			method:         http.MethodPost,
			path:           "/update//test_metric/10",
			setupMock:      func(m *MockMetricsService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid metric type",
			method:         http.MethodPost,
			path:           "/update/wrong_type/test_metric/10",
			setupMock:      func(m *MockMetricsService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty metric name",
			method:         http.MethodPost,
			path:           "/update/counter//10",
			setupMock:      func(m *MockMetricsService) {},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid value",
			method:         http.MethodPost,
			path:           "/update/counter/test_counter/abc",
			setupMock:      func(m *MockMetricsService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid counter value",
			method:         http.MethodPost,
			path:           "/update/counter/test_counter/abc",
			setupMock:      func(m *MockMetricsService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid gauge value",
			method:         http.MethodPost,
			path:           "/update/gauge/test_gauge/abc",
			setupMock:      func(m *MockMetricsService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "service error",
			method: http.MethodPost,
			path:   "/update/counter/test_counter/10",
			setupMock: func(m *MockMetricsService) {
				m.On("Add", "test_counter", mock.Anything).Return(errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockMetricsService)
			tt.setupMock(mockService)

			handler := &MetricsHandler{
				metricsService: mockService,
			}

			r := chi.NewRouter()
			r.Post("/update/{type}/{name}/{value}", handler.UpdateMetrics)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			mockService.AssertExpectations(t)
		})
	}
}
