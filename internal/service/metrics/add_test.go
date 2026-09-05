package metricsservice

import (
	"errors"
	"testing"

	models "github.com/nikitavaulin/metrics/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMetricsStorage struct {
	mock.Mock
}

func (m *MockMetricsStorage) Add(name string, metric any) error {
	args := m.Called(name, metric)
	return args.Error(0)
}

func (m *MockMetricsStorage) GetOrNil(name string) (any, error) {
	args := m.Called(name)
	return args.Get(0), args.Error(1)
}

func (m *MockMetricsStorage) Update(name string, updated any) error {
	args := m.Called(name, updated)
	return args.Error(0)
}

func (m *MockMetricsStorage) Get(name string) (any, error) {
	args := m.Called(name)
	return args.Get(0), args.Error(1)
}

func (m *MockMetricsStorage) GetList() (map[string]any, error) {
	args := m.Called()
	if data := args.Get(0); data != nil {
		return data.(map[string]any), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestMetricsService_Add(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*MockMetricsStorage)
		mName     string
		metric    models.Metrics
		wantErr   bool
	}{
		{
			name: "successful add new gauge metric",
			setupMock: func(m *MockMetricsStorage) {
				m.On("GetOrNil", "testGauge").Return(nil, nil)
				m.On("Add", "testGauge", float64(123.45)).Return(nil)
			},
			mName: "testGauge",
			metric: models.Metrics{
				ID:    "testGauge",
				MType: models.Gauge,
				Value: func() *float64 { v := 123.45; return &v }(),
			},
			wantErr: false,
		},
		{
			name: "successful add new counter metric",
			setupMock: func(m *MockMetricsStorage) {
				m.On("GetOrNil", "testCounter").Return(nil, nil)
				m.On("Add", "testCounter", int64(10)).Return(nil)
			},
			mName: "testCounter",
			metric: models.Metrics{
				ID:    "testCounter",
				MType: models.Counter,
				Delta: func() *int64 { v := int64(10); return &v }(),
			},
			wantErr: false,
		},
		{
			name: "successful update existing gauge",
			setupMock: func(m *MockMetricsStorage) {
				m.On("GetOrNil", "testGauge").Return(float64(100.0), nil)
				m.On("Update", "testGauge", float64(123.45)).Return(nil)
			},
			mName: "testGauge",
			metric: models.Metrics{
				ID:    "testGauge",
				MType: models.Gauge,
				Value: func() *float64 { v := 123.45; return &v }(),
			},
			wantErr: false,
		},
		{
			name: "successful update existing counter",
			setupMock: func(m *MockMetricsStorage) {
				m.On("GetOrNil", "testCounter").Return(int64(5), nil)
				m.On("Update", "testCounter", int64(15)).Return(nil)
			},
			mName: "testCounter",
			metric: models.Metrics{
				ID:    "testCounter",
				MType: models.Counter,
				Delta: func() *int64 { v := int64(10); return &v }(),
			},
			wantErr: false,
		},
		{
			name: "empty metric name",
			setupMock: func(m *MockMetricsStorage) {
			},
			mName: "",
			metric: models.Metrics{
				ID:    "",
				MType: models.Gauge,
				Value: func() *float64 { v := 123.45; return &v }(),
			},
			wantErr: true,
		},
		{
			name: "empty metric - both value and delta nil",
			setupMock: func(m *MockMetricsStorage) {
			},
			mName: "test",
			metric: models.Metrics{
				ID:    "test",
				MType: models.Gauge,
				Value: nil,
				Delta: nil,
			},
			wantErr: true,
		},
		{
			name: "unknown metric type",
			setupMock: func(m *MockMetricsStorage) {
			},
			mName: "test",
			metric: models.Metrics{
				ID:    "test",
				MType: "unknown",
				Value: func() *float64 { v := 123.45; return &v }(),
			},
			wantErr: true,
		},
		{
			name: "counter with nil delta",
			setupMock: func(m *MockMetricsStorage) {
			},
			mName: "test",
			metric: models.Metrics{
				ID:    "test",
				MType: models.Counter,
				Delta: nil,
			},
			wantErr: true,
		},
		{
			name: "gauge with nil value",
			setupMock: func(m *MockMetricsStorage) {
			},
			mName: "test",
			metric: models.Metrics{
				ID:    "test",
				MType: models.Gauge,
				Value: nil,
			},
			wantErr: true,
		},
		{
			name: "error when getting from storage",
			setupMock: func(m *MockMetricsStorage) {
				m.On("GetOrNil", "test").Return(nil, errors.New("storage error"))
			},
			mName: "test",
			metric: models.Metrics{
				ID:    "test",
				MType: models.Gauge,
				Value: func() *float64 { v := 123.45; return &v }(),
			},
			wantErr: true,
		},
		{
			name: "error when adding to storage",
			setupMock: func(m *MockMetricsStorage) {
				m.On("GetOrNil", "test").Return(nil, nil)
				m.On("Add", "test", float64(123.45)).Return(errors.New("add error"))
			},
			mName: "test",
			metric: models.Metrics{
				ID:    "test",
				MType: models.Gauge,
				Value: func() *float64 { v := 123.45; return &v }(),
			},
			wantErr: true,
		},
		{
			name: "error when updating storage",
			setupMock: func(m *MockMetricsStorage) {
				m.On("GetOrNil", "test").Return(float64(100.0), nil)
				m.On("Update", "test", float64(123.45)).Return(errors.New("update error"))
			},
			mName: "test",
			metric: models.Metrics{
				ID:    "test",
				MType: models.Gauge,
				Value: func() *float64 { v := 123.45; return &v }(),
			},
			wantErr: true,
		},
		{
			name: "counter type conversion error",
			setupMock: func(m *MockMetricsStorage) {
				m.On("GetOrNil", "test").Return("invalid_type", nil)
			},
			mName: "test",
			metric: models.Metrics{
				ID:    "test",
				MType: models.Counter,
				Delta: func() *int64 { v := int64(10); return &v }(),
			},
			wantErr: true,
		},
		{
			name: "counter with large delta value",
			setupMock: func(m *MockMetricsStorage) {
				m.On("GetOrNil", "testCounter").Return(int64(9223372036854775800), nil)
				m.On("Update", "testCounter", int64(9223372036854775807)).Return(nil)
			},
			mName: "testCounter",
			metric: models.Metrics{
				ID:    "testCounter",
				MType: models.Counter,
				Delta: func() *int64 { v := int64(7); return &v }(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockMetricsStorage)
			tt.setupMock(mockStorage)
			s := New(mockStorage)

			gotErr := s.Add(tt.mName, tt.metric)
			if tt.wantErr {
				assert.Error(t, gotErr, "should be error")
			} else {
				assert.NoError(t, gotErr, "should be nil")
			}
			mockStorage.AssertExpectations(t)
		})
	}
}

func Test_validateMetric(t *testing.T) {
	tests := []struct {
		name    string
		mName   string
		metric  models.Metrics
		wantErr bool
	}{
		{
			name:  "valid counter metric",
			mName: "test_counter",
			metric: models.Metrics{
				MType: models.Counter,
				Delta: func() *int64 { v := int64(10); return &v }(),
			},
			wantErr: false,
		},
		{
			name:  "valid gauge metric",
			mName: "test_gauge",
			metric: models.Metrics{
				MType: models.Gauge,
				Value: func() *float64 { v := 10.5; return &v }(),
			},
			wantErr: false,
		},
		{
			name:  "empty metric name",
			mName: "",
			metric: models.Metrics{
				MType: models.Counter,
				Delta: func() *int64 { v := int64(10); return &v }(),
			},
			wantErr: true,
		},
		{
			name:  "counter without delta",
			mName: "test_counter",
			metric: models.Metrics{
				MType: models.Counter,
				Delta: nil,
			},
			wantErr: true,
		},
		{
			name:  "gauge without value",
			mName: "test_gauge",
			metric: models.Metrics{
				MType: models.Gauge,
				Value: nil,
			},
			wantErr: true,
		},
		{
			name:  "unknown metric type",
			mName: "test_unknown",
			metric: models.Metrics{
				MType: "unknown_type",
			},
			wantErr: true,
		},
		{
			name:  "counter with zero delta",
			mName: "test_counter_zero",
			metric: models.Metrics{
				MType: models.Counter,
				Delta: func() *int64 { v := int64(0); return &v }(),
			},
			wantErr: false,
		},
		{
			name:  "gauge with zero value",
			mName: "test_gauge_zero",
			metric: models.Metrics{
				MType: models.Gauge,
				Value: func() *float64 { v := 0.0; return &v }(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := validateMetric(tt.mName, tt.metric)
			if tt.wantErr {
				assert.Error(t, gotErr, "should be error")
			} else {
				assert.NoError(t, gotErr, "should not be error")
			}
		})
	}
}
