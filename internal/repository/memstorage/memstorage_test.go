package memstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	storage := New()
	assert.NotNil(t, storage)
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name        string
		metricName  string
		metricValue any
		setup       func(*MemStorage)
		expectError bool
	}{
		{
			name:        "successful add",
			metricName:  "test_metric",
			metricValue: 42,
			setup:       func(s *MemStorage) {},
			expectError: false,
		},
		{
			name:        "duplicate metric",
			metricName:  "existing_metric",
			metricValue: 100,
			setup: func(s *MemStorage) {
				s.storage.Store("existing_metric", 50)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			tt.setup(s)

			err := s.Add(tt.metricName, tt.metricValue)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				val, ok := s.storage.Load(tt.metricName)
				assert.True(t, ok)
				assert.Equal(t, tt.metricValue, val)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name        string
		metricName  string
		setup       func(*MemStorage)
		expectedVal any
		expectError bool
	}{
		{
			name:        "existing metric",
			metricName:  "test_metric",
			setup:       func(s *MemStorage) { s.storage.Store("test_metric", 42) },
			expectedVal: 42,
			expectError: false,
		},
		{
			name:        "non-existing metric",
			metricName:  "non_existing",
			setup:       func(s *MemStorage) {},
			expectedVal: nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			tt.setup(s)

			val, err := s.Get(tt.metricName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVal, val)
			}
		})
	}
}

func TestGetOrNil(t *testing.T) {
	tests := []struct {
		name        string
		metricName  string
		setup       func(*MemStorage)
		expectedVal any
	}{
		{
			name:        "existing metric",
			metricName:  "test_metric",
			setup:       func(s *MemStorage) { s.storage.Store("test_metric", 42) },
			expectedVal: 42,
		},
		{
			name:        "non-existing metric",
			metricName:  "non_existing",
			setup:       func(s *MemStorage) {},
			expectedVal: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			tt.setup(s)

			val, err := s.GetOrNil(tt.metricName)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedVal, val)
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name        string
		metricName  string
		newValue    any
		setup       func(*MemStorage)
		expectError bool
	}{
		{
			name:        "successful update",
			metricName:  "test_metric",
			newValue:    100,
			setup:       func(s *MemStorage) { s.storage.Store("test_metric", 42) },
			expectError: false,
		},
		{
			name:        "update non-existing",
			metricName:  "non_existing",
			newValue:    100,
			setup:       func(s *MemStorage) {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			tt.setup(s)

			err := s.Update(tt.metricName, tt.newValue)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				val, ok := s.storage.Load(tt.metricName)
				assert.True(t, ok)
				assert.Equal(t, tt.newValue, val)
			}
		})
	}
}

func TestGetList(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*MemStorage)
		expected  map[string]any
		expectErr bool
	}{
		{
			name:      "empty storage",
			setup:     func(s *MemStorage) {},
			expected:  map[string]any{},
			expectErr: false,
		},
		{
			name: "multiple metrics",
			setup: func(s *MemStorage) {
				s.storage.Store("metric1", 42)
				s.storage.Store("metric2", "value")
				s.storage.Store("metric3", 3.14)
			},
			expected: map[string]any{
				"metric1": 42,
				"metric2": "value",
				"metric3": 3.14,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			tt.setup(s)

			metrics, err := s.GetList()

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, metrics)
			}
		})
	}
}
